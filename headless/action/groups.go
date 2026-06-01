package action

import (
	"fmt"
	"sort"
	"sync"

	"flclash-headless/app"
	"flclash-headless/coreclient"
	"flclash-headless/model"
)

const delayTestConcurrency = 8

var groupTypes = map[string]bool{
	"Selector":    true,
	"URLTest":     true,
	"Fallback":    true,
	"LoadBalance": true,
	"Relay":       true,
}

func isGroupType(t string) bool {
	return groupTypes[t]
}

func RefreshGroups(a *app.App) error {
	if a.CoreClient.Status() != coreclient.StatusRunning {
		return nil
	}

	proxies, err := a.CoreClient.GetProxies()
	if err != nil {
		return err
	}

	order := make(map[string]int, len(proxies.All))
	for i, name := range proxies.All {
		order[name] = i
	}

	groups := make([]model.GroupSummary, 0)
	for groupName := range proxies.Proxies {
		groupRaw := proxies.Proxies[groupName]
		groupMap, ok := groupRaw.(map[string]interface{})
		if !ok {
			continue
		}

		groupType := getStringField(groupMap, "type")
		if !isGroupType(groupType) {
			continue
		}

		summary := model.GroupSummary{
			Name: groupName,
			Type: groupType,
			Now:  getStringField(groupMap, "now"),
		}
		if all, ok := groupMap["all"].([]interface{}); ok {
			summary.Total = len(all)
		}
		groups = append(groups, summary)

		oldDelays := groupDelayMap(a.State.GetGroupDetail(groupName))
		detail := &model.GroupDetail{
			Name: groupName,
			Type: summary.Type,
			Now:  summary.Now,
		}
		if all, ok := groupMap["all"].([]interface{}); ok {
			currentNow := summary.Now
			for _, nodeRaw := range all {
				var nodeName string
				switch v := nodeRaw.(type) {
				case string:
					nodeName = v
				case map[string]interface{}:
					nodeName = getStringField(v, "name")
				default:
					continue
				}
				delay := -1
				if oldDelay, ok := oldDelays[nodeName]; ok {
					delay = oldDelay
				}
				node := model.GroupNode{
					Name:  nodeName,
					Type:  getProxyType(proxies, nodeName),
					Delay: delay,
					Now:   nodeName == currentNow,
				}
				detail.Nodes = append(detail.Nodes, node)
			}
		}
		a.State.SetGroupDetail(detail)
	}

	sort.SliceStable(groups, func(i, j int) bool {
		left, leftOK := order[groups[i].Name]
		right, rightOK := order[groups[j].Name]
		if leftOK && rightOK {
			return left < right
		}
		if leftOK != rightOK {
			return leftOK
		}
		return groups[i].Name < groups[j].Name
	})

	a.State.SetGroups(groups)
	return nil
}

func groupDelayMap(detail *model.GroupDetail) map[string]int {
	if detail == nil {
		return nil
	}
	delays := make(map[string]int, len(detail.Nodes))
	for _, node := range detail.Nodes {
		if node.Delay > 0 {
			delays[node.Name] = node.Delay
		}
	}
	return delays
}

func getProxyType(proxies *coreclient.ProxiesData, name string) string {
	if proxies == nil {
		return ""
	}
	proxyRaw, ok := proxies.Proxies[name]
	if !ok {
		return ""
	}
	proxyMap, ok := proxyRaw.(map[string]interface{})
	if !ok {
		return ""
	}
	return getStringField(proxyMap, "type")
}

func SwitchNode(a *app.App, groupName, proxyName string) error {
	if err := a.CoreClient.ChangeProxy(groupName, proxyName); err != nil {
		return err
	}

	prefs := a.StateStore.Get()
	if prefs.SelectedMap == nil {
		prefs.SelectedMap = make(map[string]string)
	}
	prefs.SelectedMap[groupName] = proxyName
	a.StateStore.SetSelectedMap(prefs.SelectedMap)

	return RefreshGroups(a)
}

func SwitchGlobalExit(a *app.App, target string) error {
	if target == "" {
		return fmt.Errorf("全局出口不能为空")
	}
	if err := a.CoreClient.ChangeProxy("GLOBAL", target); err != nil {
		return err
	}
	prefs := a.StateStore.Get()
	selected := prefs.SelectedMap
	if selected == nil {
		selected = make(map[string]string)
	}
	selected["GLOBAL"] = target
	a.StateStore.SetSelectedMap(selected)
	a.StateStore.Save()
	return RefreshGroups(a)
}

func EnsureGlobalProxySelected(a *app.App) error {
	prefs := a.StateStore.Get()
	if prefs.Mode != model.ModeGlobal {
		return nil
	}
	if err := RefreshGroups(a); err != nil {
		return err
	}
	detail := a.State.GetGroupDetail("GLOBAL")
	if detail == nil {
		return nil
	}
	if detail.Now != "" && detail.Now != "DIRECT" && detail.Now != "REJECT" {
		return nil
	}
	target := ""
	for _, n := range detail.Nodes {
		if n.Name == "\U0001f680 节点选择" {
			target = n.Name
			break
		}
	}
	if target == "" {
		for _, n := range detail.Nodes {
			if n.Name != "DIRECT" && n.Name != "REJECT" {
				target = n.Name
				break
			}
		}
	}
	if target == "" {
		return nil
	}
	if err := a.CoreClient.ChangeProxy("GLOBAL", target); err != nil {
		return err
	}
	selected := prefs.SelectedMap
	if selected == nil {
		selected = make(map[string]string)
	}
	selected["GLOBAL"] = target
	a.StateStore.SetSelectedMap(selected)
	a.StateStore.Save()
	return RefreshGroups(a)
}

func TestGroupDelay(a *app.App, groupName, testUrl string) error {
	detail := a.State.GetGroupDetail(groupName)
	if detail == nil {
		return nil
	}
	if len(detail.Nodes) == 0 {
		a.State.SetGroupDetail(detail)
		return nil
	}

	type delayResult struct {
		index int
		delay int
	}

	workerCount := delayTestConcurrency
	if len(detail.Nodes) < workerCount {
		workerCount = len(detail.Nodes)
	}

	jobs := make(chan int)
	results := make(chan delayResult, len(detail.Nodes))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				delay, err := a.CoreClient.TestDelay(testUrl, detail.Nodes[index].Name)
				if err != nil || delay <= 0 {
					delay = -1
				}
				results <- delayResult{index: index, delay: delay}
			}
		}()
	}

	for i := range detail.Nodes {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	close(results)

	for result := range results {
		detail.Nodes[result.index].Delay = result.delay
	}

	a.State.SetGroupDetail(detail)
	return nil
}

func getStringField(m map[string]interface{}, field string) string {
	if v, ok := m[field].(string); ok {
		return v
	}
	return ""
}

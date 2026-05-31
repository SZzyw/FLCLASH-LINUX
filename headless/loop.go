package main

import (
	"fmt"
	"log"

	"flclash-headless/action"
	"flclash-headless/app"
	"flclash-headless/coreclient"
	"flclash-headless/input"
	"flclash-headless/model"
	"flclash-headless/renderer"
	"flclash-headless/storage"
	"flclash-headless/util"
	"time"
)

type Loop struct {
	app            *app.App
	running        bool
	stopCh         chan struct{}
	refreshCh      chan struct{}
	buffer         string
	refreshBlocked bool
}

func NewLoop(a *app.App) *Loop {
	return &Loop{
		app:       a,
		stopCh:    make(chan struct{}),
		refreshCh: make(chan struct{}, 10),
	}
}

func (l *Loop) Run() {
	l.running = true

	if l.app.CoreClient != nil {
		l.app.CoreClient.AddListener(func(event coreclient.CoreEvent) {
			if event.Type == coreclient.EventLog {
				if data, ok := event.Data.(map[string]interface{}); ok {
					logType, _ := data["type"].(string)
					payload, _ := data["payload"].(string)
					l.app.State.AddLog(model.LogEntry{
						Type:    logType,
						Payload: payload,
					})
				}
			}
		})
	}

	go l.backgroundRefresh()

	l.renderCurrentPage()

	for l.running {
		route := l.app.PageStack.Current()
		canAutoRefresh := (route == app.RouteDashboard || route == app.RouteTraffic) && !l.refreshBlocked

		if canAutoRefresh {
			inp := l.readInputWithLiveRefresh()
			if inp == nil {
				continue
			}
			if inp.Command == "EOF" {
				l.cleanExit()
				break
			}
			l.handleInput(inp)
		} else {
			inp := input.ReadInput()
			if !l.running {
				break
			}
			if inp.Command == "EOF" {
				l.cleanExit()
				break
			}
			l.handleInput(inp)
		}

		l.renderCurrentPage()
	}

	renderer.RenderMessage("退出", "程序已正常退出。")
}

func (l *Loop) readInputWithLiveRefresh() *input.InputResult {
	inputCh := make(chan *input.InputResult, 1)
	go func() {
		inputCh <- input.ReadInput()
	}()

	for {
		select {
		case inp := <-inputCh:
			return inp
		case <-l.refreshCh:
			l.renderCurrentPage()
		}
	}
}

func (l *Loop) blockRefresh() {
	l.refreshBlocked = true
}

func (l *Loop) unblockRefresh() {
	l.refreshBlocked = false
}

func (l *Loop) cleanExit() {
	l.running = false
	l.app.StateStore.Save()
	l.app.ProfileStore.Save()
	if l.app.CoreClient != nil {
		l.app.CoreClient.Stop()
	}
}

func (l *Loop) Stop() {
	l.running = false
	l.stopCh <- struct{}{}
}

func (l *Loop) backgroundRefresh() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.refreshRuntimeData()
			select {
			case l.refreshCh <- struct{}{}:
			default:
			}
		case <-l.stopCh:
			return
		}
	}
}

func (l *Loop) refreshRuntimeData() {
	if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
		traffic, err := l.app.CoreClient.GetTraffic()
		if err != nil {
			log.Printf("GetTraffic error: %v", err)
		} else if traffic != nil {
			l.app.State.SetTraffic(*traffic)
		}
		total, err := l.app.CoreClient.GetTotalTraffic()
		if err != nil {
			log.Printf("GetTotalTraffic error: %v", err)
		} else if total != nil {
			l.app.State.SetTotalTraffic(*total)
		}
		if err := action.RefreshGroups(l.app); err != nil {
			log.Printf("RefreshGroups error: %v", err)
		}
	}
}

func (l *Loop) renderCurrentPage() {
	route := l.app.PageStack.Current()
	prefs := l.app.StateStore.Get()
	manifest := l.app.ProfileStore

	switch route {
	case app.RouteDashboard:
		renderer.RenderDashboard(l.app.State, prefs, manifest)
	case app.RouteProfiles:
		renderer.RenderProfilesList(manifest)
	case app.RouteGroupList:
		groups := l.app.State.GetGroups()
		renderer.RenderGroupList(groups)
	case app.RouteGroupDetail:
		detailName := l.buffer
		detail := l.app.State.GetGroupDetail(detailName)
		if detail != nil {
			renderer.RenderGroupDetail(detail)
		} else {
			l.app.PageStack.Pop()
			return
		}
	case app.RouteLogs:
		logs := l.app.State.GetLogs()
		renderer.RenderLogs(logs)
	case app.RouteTraffic:
		renderer.RenderTraffic(l.app.State)
	case app.RouteImportURL:
		renderer.RenderImportURLForm()
	case app.RouteImportFile:
		renderer.RenderImportFileForm()
	case app.RouteModeSwitch:
		renderer.RenderModeSwitch()
	case app.RouteConfigDetail:
		currentProfile := manifest.GetManifest().GetCurrentProfile()
		profileName := "无"
		if currentProfile != nil {
			profileName = currentProfile.Name
		}
		renderer.RenderConfigDetail(
			profileName,
			storage.ConfigFilePath(),
			string(prefs.Mode),
			prefs.MixedPort,
			prefs.TunEnabled,
			prefs.ExternalController,
		)
	case app.RouteGlobalExit:
		groups := l.app.State.GetGroups()
		detail := l.app.State.GetGroupDetail("GLOBAL")
		renderer.RenderGlobalExit(groups, detail)
	default:
		renderer.RenderDashboard(l.app.State, prefs, manifest)
	}
	renderer.RenderPrompt()
}

func (l *Loop) handleInput(inp *input.InputResult) {
	route := l.app.PageStack.Current()

	switch route {
	case app.RouteDashboard:
		l.handleDashboardInput(inp)
	case app.RouteProfiles:
		l.handleProfilesInput(inp)
	case app.RouteGroupList:
		l.handleGroupListInput(inp)
	case app.RouteGroupDetail:
		l.handleGroupDetailInput(inp)
	case app.RouteLogs:
		if inp.Command == "q" {
			l.app.PageStack.Pop()
		}
	case app.RouteTraffic:
		if inp.Command == "b" {
			l.app.PageStack.Pop()
		}
	case app.RouteImportURL:
		l.handleImportURL(inp)
	case app.RouteImportFile:
		l.handleImportFile(inp)
	case app.RouteModeSwitch:
		l.handleModeSwitch(inp)
	case app.RouteConfigDetail:
		if inp.Command == "n" {
			l.handleToggleTun()
		} else if inp.Command == "b" {
			l.app.PageStack.Pop()
		}
	case app.RouteGlobalExit:
		l.handleGlobalExitInput(inp)
	}
}

func (l *Loop) handleDashboardInput(inp *input.InputResult) {
	switch inp.Command {
	case "q":
		l.cleanExit()
	case "a":
		l.app.PageStack.Push(app.RouteImportURL)
	case "f":
		l.app.PageStack.Push(app.RouteImportFile)
	case "r":
		l.handleStartStop()
	case "g":
		groups := l.app.State.GetGroups()
		if len(groups) == 0 {
			renderer.RenderMessage("提示", "当前没有代理组，无法打开代理组。")
			return
		}
		l.app.PageStack.Push(app.RouteGroupList)
	case "l":
		l.app.PageStack.Push(app.RouteLogs)
	case "t":
		l.app.PageStack.Push(app.RouteTraffic)
	case "m":
		l.app.PageStack.Push(app.RouteModeSwitch)
	case "n":
		l.handleToggleTun()
	case "x":
		groups := l.app.State.GetGroups()
		found := false
		for _, g := range groups {
			if g.Name == "GLOBAL" {
				found = true
				break
			}
		}
		if !found {
			renderer.RenderMessage("提示", "当前没有 GLOBAL 代理组，无法选择全局出口。")
			return
		}
		if err := action.RefreshGroups(l.app); err != nil {
			log.Printf("RefreshGroups error: %v", err)
		}
		l.app.PageStack.Push(app.RouteGlobalExit)
	case "c":
		l.app.PageStack.Push(app.RouteConfigDetail)
	case "p":
		l.app.PageStack.Push(app.RouteProfiles)
	case "u":
		l.handleUpdateSubscription()
	}
}

func (l *Loop) handleToggleTun() {
	l.blockRefresh()
	defer l.unblockRefresh()

	prefs := l.app.StateStore.Get()
	next := !prefs.TunEnabled

	statusText := "开启"
	if !next {
		statusText = "关闭"
	}
	fmt.Printf("  即将%s TUN。\n", statusText)

	running := l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning
	if running {
		fmt.Println("  切换 TUN 需要重启核心，会短暂断开代理连接。")
		fmt.Print("  是否继续? [y/N]: ")
		if !input.ReadConfirm(false) {
			return
		}
	}

	old := prefs.TunEnabled
	l.app.StateStore.SetTunEnabled(next)
	if err := l.app.StateStore.Save(); err != nil {
		renderer.RenderError("保存 TUN 状态失败: " + err.Error())
		renderer.RenderPressEnter()
		return
	}

	if running {
		if err := action.StopCore(l.app); err != nil {
			l.app.StateStore.SetTunEnabled(old)
			l.app.StateStore.Save()
			renderer.RenderError("停止核心失败，已回滚 TUN 状态: " + err.Error())
			renderer.RenderPressEnter()
			return
		}

		if err := action.StartCore(l.app); err != nil {
			l.app.StateStore.SetTunEnabled(old)
			l.app.StateStore.Save()
			renderer.RenderError("启动核心失败，已回滚 TUN 状态，请手动重启: " + err.Error())
			renderer.RenderPressEnter()
			return
		}
	}

	if next {
		renderer.RenderSuccess("TUN 已开启")
	} else {
		renderer.RenderSuccess("TUN 已关闭")
	}
	renderer.RenderPressEnter()
}

func (l *Loop) handleGlobalExitInput(inp *input.InputResult) {
	switch inp.Command {
	case "b":
		l.app.PageStack.Pop()
		l.buffer = ""
		return
	}

	if !inp.IsNum || inp.Number < 1 {
		return
	}

	detail := l.app.State.GetGroupDetail("GLOBAL")
	if detail == nil {
		renderer.RenderError("无法获取 GLOBAL 组信息")
		renderer.RenderPressEnter()
		return
	}

	groups := l.app.State.GetGroups()
	options := action.BuildGlobalExitOptions(groups, detail)
	idx := inp.Number - 1
	if idx < 0 || idx >= len(options) {
		renderer.RenderError("无效编号")
		renderer.RenderPressEnter()
		return
	}

	target := options[idx].Name
	if err := action.SwitchGlobalExit(l.app, target); err != nil {
		renderer.RenderError(fmt.Sprintf("切换全局出口失败: %s", err.Error()))
	} else {
		renderer.RenderSuccess(fmt.Sprintf("全局出口已切换至: %s", target))
	}
	renderer.RenderPressEnter()
	l.app.PageStack.Pop()
	l.buffer = ""
}

func (l *Loop) handleProfilesInput(inp *input.InputResult) {
	switch inp.Command {
	case "b":
		l.app.PageStack.Pop()
	case "a":
		l.app.PageStack.Push(app.RouteImportURL)
	case "f":
		l.app.PageStack.Push(app.RouteImportFile)
	case "u":
		l.handleUpdateSubscription()
	default:
		if inp.IsNum && inp.Number >= 1 {
			profiles := l.app.ProfileStore.GetManifest().Profiles
			idx := inp.Number - 1
			if idx < len(profiles) {
				profile := profiles[idx]
				fmt.Printf("  正在切换到配置: %s\n", profile.Name)
				if err := action.SwitchProfile(l.app, &profile); err != nil {
					renderer.RenderError(err.Error())
					renderer.RenderPressEnter()
					return
				}
				renderer.RenderSuccess("已切换到配置: " + profile.Name)
				l.app.PageStack.PopTo(app.RouteDashboard)
			}
		}
	}
}

func (l *Loop) handleGroupListInput(inp *input.InputResult) {
	switch inp.Command {
	case "d":
		fmt.Println("  正在测试全部代理组延迟...")
		action.TestAllGroupsDelay(l.app, "https://www.gstatic.com/generate_204")
		renderer.RenderPressEnter()
	case "b":
		l.app.PageStack.Pop()
		l.buffer = ""
	default:
		if inp.IsNum && inp.Number >= 1 {
			groups := l.app.State.GetGroups()
			idx := inp.Number - 1
			if idx < len(groups) {
				l.buffer = groups[idx].Name
				l.app.PageStack.Push(app.RouteGroupDetail)
			}
		}
	}
}

func (l *Loop) handleGroupDetailInput(inp *input.InputResult) {
	detailName := l.buffer
	detail := l.app.State.GetGroupDetail(detailName)

	switch inp.Command {
	case "b":
		l.app.PageStack.Pop()
		l.buffer = ""
	case "t":
		if detail != nil {
			testURL := "https://www.gstatic.com/generate_204"
			fmt.Println("  正在测试延迟，请稍候...")
			action.TestGroupDelay(l.app, detailName, testURL)
			fmt.Println("  测试完成。")
			renderer.RenderPressEnter()
		}
	default:
		if inp.IsNum && inp.Number >= 1 && detail != nil {
			idx := inp.Number - 1
			if idx < len(detail.Nodes) {
				node := detail.Nodes[idx]
				if err := action.SwitchNode(l.app, detailName, node.Name); err != nil {
					renderer.RenderError(fmt.Sprintf("切换失败: %s", err.Error()))
				} else {
					fmt.Printf("  已切换: %s -> %s\n", detailName, node.Name)
				}
				renderer.RenderPressEnter()
			}
		}
	}
}

func (l *Loop) handleImportURL(inp *input.InputResult) {
	l.blockRefresh()
	defer l.unblockRefresh()

	url := inp.Text
	if url == "" || url == "b" {
		l.app.PageStack.Pop()
		return
	}

	fmt.Print("  可选备注名 (留空则自动生成): ")
	name, _ := input.ReadLine()

	autoApply := true
	fmt.Print("  是否在导入后立即切换? [Y/n]: ")
	autoApply = input.ReadConfirm(true)

	profile, err := action.ImportFromURL(l.app.ProfileStore, url, name, autoApply)
	if err != nil {
		renderer.RenderError(fmt.Sprintf("导入失败: %s", err.Error()))
		renderer.RenderPressEnter()
		l.app.PageStack.Pop()
		return
	}

	renderer.RenderSuccess("导入成功!")
	fmt.Printf("  - 名称: %s\n", profile.Name)
	fmt.Printf("  - 类型: 订阅\n")
	fmt.Printf("  - 来源: %s\n", profile.Source)

	l.app.ProfileStore.Save()
	l.app.StateStore.Save()

	if autoApply {
		fmt.Println("  正在切换到新配置...")
		if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
			if err := action.SwitchProfile(l.app, profile); err != nil {
				renderer.RenderError(fmt.Sprintf("切换失败: %s", err.Error()))
			} else {
				renderer.RenderSuccess("配置已应用")
			}
		} else {
			l.app.ProfileStore.SetCurrent(profile.ID)
			l.app.StateStore.SetCurrentProfileID(profile.ID)
			fmt.Println("  已设为当前配置，但未运行。")
		}
	}

	l.app.StateStore.Save()
	l.app.PageStack.PopTo(app.RouteDashboard)
	renderer.RenderPressEnter()
}

func (l *Loop) handleImportFile(inp *input.InputResult) {
	l.blockRefresh()
	defer l.unblockRefresh()

	filePath := inp.Text
	if filePath == "" || filePath == "b" {
		l.app.PageStack.Pop()
		return
	}

	fmt.Print("  可选备注名 (留空则使用文件名): ")
	name, _ := input.ReadLine()

	autoApply := true
	fmt.Print("  是否在导入后立即切换? [Y/n]: ")
	autoApply = input.ReadConfirm(true)

	profile, err := action.ImportFromFile(l.app.ProfileStore, filePath, name, autoApply)
	if err != nil {
		renderer.RenderError(fmt.Sprintf("导入失败: %s", err.Error()))
		renderer.RenderPressEnter()
		l.app.PageStack.Pop()
		return
	}

	renderer.RenderSuccess("导入成功!")
	fmt.Printf("  - 名称: %s\n", profile.Name)
	fmt.Printf("  - 类型: 文件\n")
	fmt.Printf("  - 来源: %s\n", profile.Source)

	if autoApply {
		fmt.Println("  正在切换到新配置...")
		if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
			if err := action.SwitchProfile(l.app, profile); err != nil {
				renderer.RenderError(fmt.Sprintf("切换失败: %s", err.Error()))
			} else {
				renderer.RenderSuccess("配置已应用")
			}
		} else {
			l.app.ProfileStore.SetCurrent(profile.ID)
			l.app.StateStore.SetCurrentProfileID(profile.ID)
			fmt.Println("  已设为当前配置，但未运行。")
		}
	}

	l.app.StateStore.Save()
	l.app.ProfileStore.Save()
	l.app.PageStack.PopTo(app.RouteDashboard)
	renderer.RenderPressEnter()
}

func (l *Loop) handleStartStop() {
	l.blockRefresh()
	defer l.unblockRefresh()

	if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
		fmt.Println("  核心当前正在运行。")
		fmt.Print("  是否立即停止? [y/N]: ")
		if input.ReadConfirm(false) {
			if err := action.StopCore(l.app); err != nil {
				renderer.RenderError(err.Error())
			} else {
				renderer.RenderSuccess("核心已停止")
			}
			renderer.RenderPressEnter()
		}
	} else {
		fmt.Println("  核心当前已停止。")
		fmt.Print("  是否立即启动? [y/N]: ")
		if input.ReadConfirm(false) {
			if err := action.StartCore(l.app); err != nil {
				renderer.RenderError(err.Error())
			} else {
				renderer.RenderSuccess("核心已启动")
			}
			renderer.RenderPressEnter()
		}
	}
}

func (l *Loop) handleUpdateSubscription() {
	l.blockRefresh()
	defer l.unblockRefresh()

	manifest := l.app.ProfileStore.GetManifest()
	profile := manifest.GetCurrentProfile()
	if profile == nil {
		renderer.RenderMessage("提示", "当前没有可用配置。")
		return
	}
	if profile.Type == model.ProfileTypeFile {
		renderer.RenderMessage("提示", "当前配置不是订阅类型，无法执行更新。")
		return
	}

	fmt.Println("  正在更新当前订阅...")
	newProfile, err := action.ImportFromURL(l.app.ProfileStore, profile.Source, profile.Name, false)
	if err != nil {
		renderer.RenderError(fmt.Sprintf("更新失败: %s", err.Error()))
		renderer.RenderPressEnter()
		return
	}

	renderer.RenderSuccess("更新成功!")
	fmt.Print("  是否立即重新应用? [Y/n]: ")
	if input.ReadConfirm(true) {
		if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
			action.SwitchProfile(l.app, newProfile)
		}
		renderer.RenderSuccess("配置已重新应用")
	} else {
		fmt.Println("  更新成功，但未立即应用。")
	}

	l.app.ProfileStore.Save()
	l.app.StateStore.Save()
	renderer.RenderPressEnter()
}

func (l *Loop) handleModeSwitch(inp *input.InputResult) {
	switch inp.Command {
	case "b":
		l.app.PageStack.Pop()
	case "1":
		l.app.StateStore.SetMode(model.ModeRule)
	case "2":
		l.app.StateStore.SetMode(model.ModeGlobal)
	case "3":
		l.app.StateStore.SetMode(model.ModeDirect)
	}

	if inp.Command == "1" || inp.Command == "2" || inp.Command == "3" {
		modes := map[string]model.Mode{"1": model.ModeRule, "2": model.ModeGlobal, "3": model.ModeDirect}
		mode := modes[inp.Command]
		modeName := util.FormatMode(string(mode))
		fmt.Println()
		renderer.RenderSuccess("运行模式已切换为: " + modeName)

		if l.app.CoreClient != nil && l.app.CoreClient.Status() == coreclient.StatusRunning {
			profile := l.app.ProfileStore.GetManifest().GetCurrentProfile()
			if profile != nil {
				action.ApplyProfile(l.app, profile)
			}
		}

		l.app.StateStore.Save()
		renderer.RenderPressEnter()
		l.app.PageStack.Pop()
	}
}

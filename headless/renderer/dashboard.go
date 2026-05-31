package renderer

import (
	"fmt"
	"time"

	"flclash-headless/app"
	"flclash-headless/coreclient"
	"flclash-headless/model"
	"flclash-headless/storage"
	"flclash-headless/util"
)

func RenderDashboard(appState *app.RuntimeState, prefs *model.RuntimePrefs, manifest *storage.ProfileStore) {
	util.ClearScreen()

	hostname := util.GetHostname()
	title := "FlClash 无桌面终端版 v" + model.Version
	if hostname != "" {
		title += "  |  主机: " + hostname
	}
	RenderHeader(title)

	coreStatus := appState.GetCoreStatus()
	currentProfile := manifest.GetManifest().GetCurrentProfile()

	statusText := ""
	switch coreStatus {
	case coreclient.StatusRunning:
		statusText = util.ColorGreen("运行中")
		uptime := appState.CoreStartTime
		if !uptime.IsZero() {
			statusText += "  |  运行时长: " + util.FormatDuration(time.Since(uptime))
		}
	case coreclient.StatusStarting:
		statusText = util.ColorYellow("启动中")
	case coreclient.StatusError:
		statusText = util.ColorRed("错误")
	default:
		statusText = util.ColorYellow("未启动")
	}

	RenderInfoLine("核心状态", statusText)

	if currentProfile != nil {
		typeText := "订阅"
		if currentProfile.Type == model.ProfileTypeFile {
			typeText = "文件"
		}
		RenderInfoLine("当前配置", currentProfile.Name+"  ("+typeText+")")
	} else {
		RenderInfoLine("当前配置", util.ColorYellow("无"))
	}

	RenderInfoLine("当前模式", util.FormatMode(string(prefs.Mode)))
	globalExitText := getGlobalExitText(appState, prefs)
	if globalExitText != "" {
		RenderInfoLine("全局出口", globalExitText)
	}
	RenderInfoLine("TUN", boolText(prefs.TunEnabled))
	RenderInfoLine("混合端口", fmt.Sprintf("%d", prefs.MixedPort))

	traffic := appState.GetTraffic()
	total := appState.GetTotalTraffic()
	RenderInfoLine("上传速度", util.FormatSpeed(traffic.Up))
	RenderInfoLine("下载速度", util.FormatSpeed(traffic.Down))
	RenderInfoLine("总上传", util.FormatTraffic(total.Up))
	RenderInfoLine("总下载", util.FormatTraffic(total.Down))

	groups := appState.GetGroups()
	if prefs.Mode == model.ModeGlobal {
		for _, g := range groups {
			if g.Name == "GLOBAL" && (g.Now == "DIRECT" || g.Now == "REJECT") {
				fmt.Println("  " + util.ColorYellow("注意: 全局模式下 GLOBAL 为 DIRECT，外网流量不会走代理。"))
				fmt.Println("  " + util.ColorYellow("请按 [x] 选择全局出口。"))
			}
		}
	}

	if len(groups) > 0 {
		RenderSubHeader("代理组摘要")
		for i, g := range groups {
			now := g.Now
			if now == "" {
				now = "--"
			}
			fmt.Printf("  [%d] %-20s -> %s\n", i+1, util.Truncate(g.Name, 20), util.Truncate(now, 20))
		}
	} else {
		RenderSubHeader("代理组摘要")
		if coreStatus == coreclient.StatusRunning {
			fmt.Println("  " + util.ColorDim("暂无代理组数据"))
		} else {
			fmt.Println("  " + util.ColorDim("核心未运行，无法获取代理组"))
		}
	}

	RenderSection("配置相关", "")
	RenderMenu([]MenuEntry{
		{Key: "p", Label: "查看当前配置摘要"},
		{Key: "a", Label: "从 URL 添加配置"},
		{Key: "f", Label: "从文件路径添加配置"},
		{Key: "u", Label: "更新当前订阅"},
	})

	RenderSection("运行控制", "")
	runLabel := "启动核心"
	if coreStatus == coreclient.StatusRunning {
		runLabel = "停止核心"
	}
	RenderMenu([]MenuEntry{
		{Key: "r", Label: runLabel},
		{Key: "m", Label: "切换运行模式"},
		{Key: "g", Label: "打开代理组"},
		{Key: "n", Label: "切换 TUN"},
		{Key: "x", Label: "选择全局出口"},
		{Key: "c", Label: "打开运行配置详情"},
		{Key: "l", Label: "打开日志页"},
		{Key: "t", Label: "打开流量页"},
		{Key: "q", Label: "退出"},
	})

	msg := appState.GetMessage()
	if msg != "" {
		fmt.Println()
		fmt.Println("  " + util.ColorYellow("消息: "+msg))
		appState.SetMessage("")
	}
}

func boolText(v bool) string {
	if v {
		return util.ColorGreen("开启")
	}
	return util.ColorRed("关闭")
}

func getGlobalExitText(appState *app.RuntimeState, prefs *model.RuntimePrefs) string {
	if prefs.Mode != model.ModeGlobal {
		return ""
	}
	groups := appState.GetGroups()
	for _, g := range groups {
		if g.Name == "GLOBAL" {
			if g.Now == "" || g.Now == "DIRECT" || g.Now == "REJECT" {
				return util.ColorYellow(g.Now)
			}
			detail := appState.GetGroupDetail("GLOBAL")
			if detail != nil && len(detail.Nodes) > 0 {
				for _, node := range detail.Nodes {
					if node.Now {
						return fmt.Sprintf("GLOBAL -> %s", node.Name)
					}
				}
			}
			return g.Now
		}
	}
	return ""
}

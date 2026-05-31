package renderer

import (
	"fmt"

	"flclash-headless/util"
)

func RenderImportURLForm() {
	util.ClearScreen()
	RenderHeader("从 URL 添加配置")
	fmt.Println()
	fmt.Println("  " + util.ColorDim("请输入订阅地址"))
}

func RenderImportFileForm() {
	util.ClearScreen()
	RenderHeader("从文件路径添加配置")
	fmt.Println()
	fmt.Println("  " + util.ColorDim("请输入配置文件路径"))
}

func RenderModeSwitch() {
	util.ClearScreen()
	RenderHeader("运行模式")
	fmt.Println()
	fmt.Println("  [1] 规则模式")
	fmt.Println("      按配置规则分流")
	fmt.Println()
	fmt.Println("  [2] 全局模式")
	fmt.Println("      绝大多数流量走全局代理")
	fmt.Println()
	fmt.Println("  [3] 直连模式")
	fmt.Println("      不走代理")
	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "1-3", Label: "选择模式"},
		{Key: "b", Label: "返回主面板"},
	})
}

func RenderConfigDetail(profileName, configPath, mode string, mixedPort int, tunEnabled bool, externalController string) {
	util.ClearScreen()
	RenderHeader("运行配置摘要")
	fmt.Println()
	RenderInfoLine("配置名", profileName)
	RenderInfoLine("运行文件", configPath)
	RenderInfoLine("混合端口", fmt.Sprintf("%d", mixedPort))
	RenderInfoLine("模式", util.FormatMode(mode))
	tunText := "关闭"
	if tunEnabled {
		tunText = "开启"
	}
	RenderInfoLine("TUN", tunText)
	RenderInfoLine("外部控制", externalController)
	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "n", Label: "切换 TUN"},
		{Key: "b", Label: "返回主面板"},
	})
}

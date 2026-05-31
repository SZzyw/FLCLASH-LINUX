package renderer

import (
	"fmt"

	"flclash-headless/action"
	"flclash-headless/model"
	"flclash-headless/util"
)

func RenderGlobalExit(groups []model.GroupSummary, globalDetail *model.GroupDetail) {
	util.ClearScreen()
	RenderHeader("全局出口选择")

	if globalDetail == nil {
		fmt.Println()
		fmt.Println("  " + util.ColorDim("无法获取 GLOBAL 组信息。"))
		fmt.Println()
		RenderMenu([]MenuEntry{
			{Key: "b", Label: "返回主面板"},
		})
		return
	}

	fmt.Println()
	fmt.Println("  " + util.ColorBold("当前全局出口: ") + globalDetail.Now)
	fmt.Println()
	fmt.Println("  " + util.ColorDim("选择 GLOBAL 组使用的出口："))
	fmt.Println()

	options := action.BuildGlobalExitOptions(groups, globalDetail)
	if options == nil {
		fmt.Println("  " + util.ColorDim("无可用的出口选项"))
		fmt.Println()
		RenderMenu([]MenuEntry{
			{Key: "b", Label: "返回主面板"},
		})
		return
	}

	for i, node := range options {
		prefix := "  "
		if node.Now {
			prefix = util.ColorGreen("  >")
		}

		display := node.Name
		for _, g := range groups {
			if g.Name == node.Name && g.Now != "" {
				display = fmt.Sprintf("%s -> %s", node.Name, g.Now)
				break
			}
		}

		suffix := ""
		if node.Now {
			suffix = util.ColorGreen("  <- 当前")
		}
		fmt.Printf("%s[%d] %-35s %s\n", prefix, i+1, util.Truncate(display, 35), suffix)
	}

	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "数字", Label: "选择出口"},
		{Key: "b", Label: "返回主面板"},
	})
}

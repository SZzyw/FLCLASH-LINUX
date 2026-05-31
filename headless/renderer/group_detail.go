package renderer

import (
	"fmt"

	"flclash-headless/model"
	"flclash-headless/util"
)

func RenderGroupDetail(detail *model.GroupDetail) {
	util.ClearScreen()
	RenderHeader("代理组: " + detail.Name)

	currentNow := detail.Now
	if currentNow == "" {
		currentNow = "--"
	}
	RenderInfoLine("当前节点", currentNow)
	fmt.Println()

	for i, node := range detail.Nodes {
		prefix := "  "
		suffix := ""
		if node.Now {
			prefix = util.ColorGreen("  >")
			suffix = util.ColorGreen("  <- 当前")
		}
		delayText := util.FormatDelay(node.Delay)
		if node.Delay > 0 {
			delayText = util.ColorByDelay(node.Delay) + delayText + util.Reset
		}
		fmt.Printf("%s[%d] %-25s %s%s\n", prefix, i+1, util.Truncate(node.Name, 25), delayText, suffix)
	}

	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "数字", Label: "切换到该节点"},
		{Key: "t", Label: "测试全部节点延迟"},
		{Key: "b", Label: "返回代理组列表"},
	})
}

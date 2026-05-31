package renderer

import (
	"fmt"

	"flclash-headless/model"
	"flclash-headless/util"
)

func RenderGroupList(groups []model.GroupSummary) {
	util.ClearScreen()
	RenderHeader("代理组列表")

	if len(groups) == 0 {
		fmt.Println()
		fmt.Println("  " + util.ColorDim("当前没有可切换的代理组。"))
		fmt.Println()
		fmt.Println("  可能原因：")
		fmt.Println("  - 当前配置尚未成功应用")
		fmt.Println("  - 当前模式为直连模式")
		fmt.Println("  - 当前配置本身没有可选组")
		fmt.Println()
		RenderMenu([]MenuEntry{
			{Key: "b", Label: "返回主面板"},
		})
		return
	}

	for i, g := range groups {
		now := g.Now
		if now == "" {
			now = "--"
		}
		fmt.Printf("  [%d] %-20s 当前: %s\n", i+1, util.Truncate(g.Name, 20), util.Truncate(now, 20))
	}

	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "数字", Label: "打开代理组"},
		{Key: "d", Label: "测试全部延迟"},
		{Key: "b", Label: "返回主面板"},
	})
}

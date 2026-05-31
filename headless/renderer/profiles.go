package renderer

import (
	"fmt"

	"flclash-headless/model"
	"flclash-headless/storage"
	"flclash-headless/util"
)

func RenderProfilesList(manifest *storage.ProfileStore) {
	util.ClearScreen()
	RenderHeader("配置列表")

	currentID := manifest.GetManifest().CurrentProfileID
	profiles := manifest.GetManifest().Profiles

	if currentID > 0 {
		current := manifest.GetManifest().GetCurrentProfile()
		if current != nil {
			RenderInfoLine("当前正在使用", current.Name)
		}
	} else {
		RenderInfoLine("当前正在使用", util.ColorYellow("无"))
	}

	fmt.Println()

	if len(profiles) == 0 {
		fmt.Println("  " + util.ColorDim("暂无可用配置"))
	} else {
		for i, p := range profiles {
			typeText := "订阅"
			if p.Type == model.ProfileTypeFile {
				typeText = "文件"
			}
			statusText := "可切换"
			if p.ID == currentID {
				statusText = util.ColorGreen("当前使用中")
			}
			fmt.Printf("  [%d] %s\n", i+1, util.ColorBold(p.Name))
			fmt.Printf("      类型: %s    来源: %s\n", typeText, util.Truncate(p.Source, 30))
			fmt.Printf("      最近更新: %s    状态: %s\n", util.FormatTimeAgo(p.UpdatedAt), statusText)
			fmt.Println()
		}
	}

	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "a", Label: "从 URL 添加"},
		{Key: "f", Label: "从文件路径添加"},
		{Key: "u", Label: "更新当前订阅"},
		{Key: "数字", Label: "切换到指定配置"},
		{Key: "b", Label: "返回主面板"},
	})
}

func RenderProfilesCompact(manifest *storage.ProfileStore) {
	currentID := manifest.GetManifest().CurrentProfileID
	profiles := manifest.GetManifest().Profiles
	for i, p := range profiles {
		typeText := "订阅"
		if p.Type == model.ProfileTypeFile {
			typeText = "文件"
		}
		statusText := ""
		if p.ID == currentID {
			statusText = " <<"
		}
		fmt.Printf("  [%d] %-20s %-8s %s%s\n", i+1, util.Truncate(p.Name, 20), typeText, util.FormatTimeAgo(p.UpdatedAt), util.ColorGreen(statusText))
	}
}

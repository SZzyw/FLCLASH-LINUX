package renderer

import (
	"fmt"

	"flclash-headless/app"
	"flclash-headless/util"
)

func RenderTraffic(appState *app.RuntimeState) {
	util.ClearScreen()
	RenderHeader("流量信息")

	traffic := appState.GetTraffic()
	total := appState.GetTotalTraffic()

	RenderInfoLine("当前上传速度", util.FormatSpeed(traffic.Up))
	RenderInfoLine("当前下载速度", util.FormatSpeed(traffic.Down))
	fmt.Println()
	RenderInfoLine("累计上传", util.FormatTraffic(total.Up))
	RenderInfoLine("累计下载", util.FormatTraffic(total.Down))

	fmt.Println()
	RenderMenu([]MenuEntry{
		{Key: "b", Label: "返回主面板"},
	})
}

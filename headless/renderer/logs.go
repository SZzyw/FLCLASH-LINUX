package renderer

import (
	"fmt"

	"flclash-headless/model"
	"flclash-headless/util"
)

func RenderLogs(logs []model.LogEntry) {
	util.ClearScreen()
	RenderHeader("日志  (按 q 返回主面板)")

	if len(logs) == 0 {
		fmt.Println("  " + util.ColorDim("暂无日志"))
		return
	}

	start := 0
	if len(logs) > 50 {
		start = len(logs) - 50
	}
	for i := start; i < len(logs); i++ {
		entry := logs[i]
		timeStr := util.FormatTimeAgo(model.TimeNow())
		levelText := util.FormatLogLevel(entry.Type)
		colored := ""
		switch entry.Type {
		case "error":
			colored = util.ColorRed(levelText)
		case "warning":
			colored = util.ColorYellow(levelText)
		default:
			colored = levelText
		}
		msg := util.Truncate(entry.Payload, 80)
		fmt.Printf("  [%s] [%s] %s\n", timeStr, colored, msg)
	}
}

package renderer

import (
	"fmt"

	"flclash-headless/util"
)

func RenderMessage(title, message string) {
	util.ClearScreen()
	RenderHeader(title)
	fmt.Println()
	lines := splitLines(message, BoxWidth-4)
	for _, line := range lines {
		fmt.Println("  " + line)
	}
	fmt.Println()
	RenderPressEnter()
}

func RenderConfirm(title, message string) {
	util.ClearScreen()
	RenderHeader(title)
	fmt.Println()
	lines := splitLines(message, BoxWidth-4)
	for _, line := range lines {
		fmt.Println("  " + line)
	}
	fmt.Println()
	fmt.Print("  " + util.ColorYellow("确认? [y/N]: "))
}

func splitLines(s string, maxLen int) []string {
	var result []string
	runes := []rune(s)
	for len(runes) > 0 {
		if len(runes) <= maxLen {
			result = append(result, string(runes))
			break
		}
		result = append(result, string(runes[:maxLen]))
		runes = runes[maxLen:]
	}
	return result
}

package renderer

import (
	"fmt"
	"strings"

	"flclash-headless/input"
	"flclash-headless/util"
)

const BoxWidth = 66

type MenuEntry struct {
	Key   string
	Label string
}

func RenderHeader(title string) {
	w := BoxWidth
	fmt.Println(util.ColorCyan(util.RepeatChar("=", w)))
	fmt.Println(util.ColorBold(util.PadCenter(title, w)))
	fmt.Println(util.ColorCyan(util.RepeatChar("=", w)))
}

func RenderSubHeader(title string) {
	fmt.Println()
	fmt.Println(util.ColorCyan(util.RepeatChar("-", BoxWidth)))
	fmt.Println(util.ColorBold("  " + title))
	fmt.Println(util.ColorCyan(util.RepeatChar("-", BoxWidth)))
}

func RenderInfoLine(label, value string) {
	fmt.Printf("  %s: %s\n", util.ColorBold(label), value)
}

func RenderSection(title string, content string) {
	fmt.Println()
	fmt.Println(util.ColorCyan("  " + title))
	fmt.Println(util.ColorCyan(util.RepeatChar("-", BoxWidth-2)))
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fmt.Println("  " + line)
	}
}

func RenderMenu(entries []MenuEntry) {
	fmt.Println()
	for _, e := range entries {
		keyStr := util.ColorYellow(fmt.Sprintf("[%s]", e.Key))
		fmt.Printf("  %s %s\n", keyStr, e.Label)
	}
}

func RenderPrompt() {
	fmt.Print("\n  " + util.ColorGreen("输入> "))
}

func RenderError(msg string) {
	fmt.Println(util.ColorRed(fmt.Sprintf("  !! 错误: %s", msg)))
}

func RenderSuccess(msg string) {
	fmt.Println(util.ColorGreen(fmt.Sprintf("  OK %s", msg)))
}

func RenderWait(text string) {
	fmt.Println("  " + text)
}

func RenderPressEnter() {
	fmt.Print("  " + util.ColorDim("按回车继续..."))
	input.ReadLine()
}

package renderer

import (
	"fmt"

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

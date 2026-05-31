package util

import "unicode/utf8"

func DisplayWidth(s string) int {
	return utf8.RuneCountInString(s)
}

func PadCenter(s string, width int) string {
	w := DisplayWidth(s)
	if w >= width {
		return s
	}
	left := (width - w) / 2
	right := width - w - left
	return RepeatChar(" ", left) + s + RepeatChar(" ", right)
}

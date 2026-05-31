package util

import (
	"fmt"
	"time"
	"unicode/utf8"
)

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func FormatMode(mode string) string {
	switch mode {
	case "rule":
		return "规则模式"
	case "global":
		return "全局模式"
	case "direct":
		return "直连模式"
	default:
		return mode
	}
}

func FormatDelay(delay int) string {
	if delay <= 0 {
		return "--"
	}
	return fmt.Sprintf("%d ms", delay)
}

var hostname string

func GetHostname() string {
	return hostname
}

func SetHostname(h string) {
	hostname = h
}

func Truncate(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen-1]) + "..."
}

func AutoName(url string) string {
	if url == "" {
		return "未命名"
	}
	return url
}

func RepeatChar(ch string, count int) string {
	if count <= 0 {
		return ""
	}
	result := make([]byte, count)
	for i := range result {
		result[i] = ch[0]
	}
	return string(result)
}

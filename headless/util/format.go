package util

import (
	"fmt"
	"time"
	"unicode/utf8"
)

func FormatTraffic(bytes int64) string {
	if bytes < 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	if bytes == 0 {
		return "0 B"
	}
	size := float64(bytes)
	unitIdx := 0
	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}
	if unitIdx == 0 {
		return fmt.Sprintf("%.0f %s", size, units[unitIdx])
	}
	return fmt.Sprintf("%.1f %s", size, units[unitIdx])
}

func FormatSpeed(bytesPerSec int64) string {
	return FormatTraffic(bytesPerSec) + "/s"
}

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

func FormatTimeAgo(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return "刚刚"
	case diff < time.Hour:
		return fmt.Sprintf("%d 分钟前", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d 小时前", int(diff.Hours()))
	default:
		return fmt.Sprintf("%d 天前", int(diff.Hours()/24))
	}
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

func FormatLogLevel(level string) string {
	switch level {
	case "debug":
		return "调试"
	case "info":
		return "信息"
	case "warning":
		return "警告"
	case "error":
		return "错误"
	default:
		return level
	}
}

func FormatDelay(delay int) string {
	if delay <= 0 {
		return "--"
	}
	return fmt.Sprintf("%d ms", delay)
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

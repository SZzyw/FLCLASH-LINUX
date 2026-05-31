package util

import "fmt"

const (
	Reset   = "[0m"
	Bold    = "[1m"
	Dim     = "[2m"
	Red     = "[31m"
	Green   = "[32m"
	Yellow  = "[33m"
	Blue    = "[34m"
	Magenta = "[35m"
	Cyan    = "[36m"
	White   = "[37m"
)

func ColorRed(s string) string { return Red + s + Reset }
func ColorGreen(s string) string { return Green + s + Reset }
func ColorYellow(s string) string { return Yellow + s + Reset }
func ColorBlue(s string) string { return Blue + s + Reset }
func ColorCyan(s string) string { return Cyan + s + Reset }
func ColorMagenta(s string) string { return Magenta + s + Reset }
func ColorBold(s string) string { return Bold + s + Reset }
func ColorDim(s string) string { return Dim + s + Reset }

func ColorByDelay(delay int) string {
	switch {
	case delay <= 0:
		return Red
	case delay <= 200:
		return Green
	case delay <= 500:
		return Yellow
	default:
		return Red
	}
}

func ClearScreen() {
	fmt.Print("[2J[H")
}

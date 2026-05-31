package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var globalReader *bufio.Reader

func ensureReader() *bufio.Reader {
	if globalReader == nil {
		globalReader = bufio.NewReader(os.Stdin)
	}
	return globalReader
}

func ResetReader() {
	globalReader = nil
}

type InputResult struct {
	Command string
	Text    string
	Number  int
	IsNum   bool
	Confirm *bool
}

func ReadLine() (string, error) {
	reader := ensureReader()
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func ReadInput() *InputResult {
	text, err := ReadLine()
	if err != nil {
		return &InputResult{Command: "EOF"}
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return &InputResult{Command: ""}
	}

	result := &InputResult{
		Command: strings.ToLower(text),
		Text:    text,
	}

	if num, err := parseInt(text); err == nil {
		result.Number = num
		result.IsNum = true
	}

	return result
}

func ReadText(prompt string) string {
	fmt.Print(prompt)
	text, _ := ReadLine()
	return strings.TrimSpace(text)
}

func ReadConfirm(defaultYes bool) bool {
	text, _ := ReadLine()
	text = strings.TrimSpace(strings.ToLower(text))
	if text == "" {
		return defaultYes
	}
	return text == "y" || text == "yes"
}

func ReadNotEmpty(prompt string) string {
	for {
		text := ReadText(prompt)
		if text != "" {
			return text
		}
		fmt.Println("  输入不能为空，请重新输入。")
	}
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

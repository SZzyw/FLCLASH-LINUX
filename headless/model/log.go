package model

type LogEntry struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type LogLevel string

const (
	LogDebug   LogLevel = "debug"
	LogInfo    LogLevel = "info"
	LogWarning LogLevel = "warning"
	LogError   LogLevel = "error"
)

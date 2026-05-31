package model

type Mode string

const (
	ModeRule   Mode = "rule"
	ModeGlobal Mode = "global"
	ModeDirect Mode = "direct"
)

type RuntimePrefs struct {
	CurrentProfileID  int64             `json:"current_profile_id"`
	Mode              Mode              `json:"mode"`
	TunEnabled        bool              `json:"tun_enabled"`
	SystemProxy       bool              `json:"system_proxy"`
	MixedPort         int               `json:"mixed_port"`
	ExternalController string           `json:"external_controller"`
	LogLevel          string            `json:"log_level"`
	SelectedMap       map[string]string `json:"selected_map"`
	LastRunning       bool              `json:"last_running"`
}

func NewRuntimePrefs() *RuntimePrefs {
	return &RuntimePrefs{
		Mode:              ModeRule,
		MixedPort:         7890,
		ExternalController: "127.0.0.1:9090",
		LogLevel:          "info",
		SelectedMap:       make(map[string]string),
	}
}

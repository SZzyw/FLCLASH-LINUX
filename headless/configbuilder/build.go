package configbuilder

import (
	"fmt"
	"os"

	"flclash-headless/model"
	"flclash-headless/storage"
	"gopkg.in/yaml.v3"
)

type BuiltConfig struct {
	Path string
}

const tunDeviceName = "flclash0"

func Build(profile *model.ProfileRecord, prefs *model.RuntimePrefs) (*BuiltConfig, error) {
	configPath := storage.ConfigFilePath()

	raw, err := os.ReadFile(profile.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read profile: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(raw, &config); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	config["mixed-port"] = prefs.MixedPort
	config["mode"] = string(prefs.Mode)
	config["log-level"] = prefs.LogLevel
	config["external-controller"] = prefs.ExternalController
	config["allow-lan"] = true
	config["ipv6"] = false

	config["tun"] = map[string]interface{}{
		"enable":                prefs.TunEnabled,
		"device":                tunDeviceName,
		"stack":                 "system",
		"auto-route":            true,
		"auto-detect-interface": true,
		"strict-route":          true,
		"dns-hijack":            []string{"any:53"},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("marshal YAML: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write config: %w", err)
	}

	return &BuiltConfig{Path: configPath}, nil
}

func WriteRawYAML(profileID int64, data []byte) (string, error) {
	path := storage.ProfileFilePath(profileID)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func ValidateRawYAML(data []byte) error {
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	if len(config) == 0 {
		return fmt.Errorf("配置不是有效的 Clash YAML")
	}
	if _, ok := config["proxies"]; !ok {
		if _, ok := config["proxy-providers"]; !ok {
			return fmt.Errorf("配置缺少 proxies 或 proxy-providers")
		}
	}
	return nil
}

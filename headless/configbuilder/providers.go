package configbuilder

import (
	"fmt"
	"os"
	"path/filepath"

	"flclash-headless/storage"
)

func fixProviderPaths(config map[string]interface{}, profileID int64) error {
	providerDir := storage.ProviderDirForProfile(profileID)
	if err := os.MkdirAll(providerDir, 0755); err != nil {
		return fmt.Errorf("create provider dir: %w", err)
	}

	fixProxyProviders(config, providerDir)
	fixRuleProviders(config, providerDir)
	return nil
}

func fixProxyProviders(config map[string]interface{}, providerDir string) {
	raw, ok := config["proxy-providers"]
	if !ok {
		return
	}
	providers, ok := raw.(map[string]interface{})
	if !ok {
		return
	}
	for name, provider := range providers {
		p, ok := provider.(map[string]interface{})
		if !ok {
			continue
		}
		if path, ok := p["path"].(string); ok && path != "" {
			p["path"] = filepath.Join(providerDir, "proxies", name, filepath.Base(path))
		}
	}
}

func fixRuleProviders(config map[string]interface{}, providerDir string) {
	raw, ok := config["rule-providers"]
	if !ok {
		return
	}
	providers, ok := raw.(map[string]interface{})
	if !ok {
		return
	}
	for name, provider := range providers {
		p, ok := provider.(map[string]interface{})
		if !ok {
			continue
		}
		if path, ok := p["path"].(string); ok && path != "" {
			p["path"] = filepath.Join(providerDir, "rules", name, filepath.Base(path))
		}
	}
}

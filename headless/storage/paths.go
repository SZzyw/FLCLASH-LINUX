package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DataDirName    = "flclash-headless"
	ProfilesDir    = "profiles"
	ProvidersDir   = "providers"
	CacheDir       = "cache"
	LogsDir        = "logs"
	StateFile      = "state.json"
	ProfilesFile   = "profiles.json"
	ConfigFile     = "config.yaml"
	assetsRel      = "../assets/data"
)

var dataDir string

func GetDataDir() string {
	if dataDir != "" {
		return dataDir
	}
	if dir := findAssetsDataDir(); dir != "" {
		dataDir = dir
		return dataDir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		dataDir = filepath.Join("/tmp", DataDirName)
		return dataDir
	}
	dataDir = filepath.Join(home, ".local", "share", DataDirName)
	return dataDir
}

func findAssetsDataDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Clean(filepath.Join(filepath.Dir(exe), assetsRel))
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir
	}
	return ""
}

func SetDataDir(dir string) {
	dataDir = dir
}

func EnsureDirs() error {
	dirs := []string{
		GetDataDir(),
		ProfileDir(),
		ProvidersRootDir(),
		CacheDirPath(),
		LogDirPath(),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}
	return nil
}

func ProfileDir() string {
	return filepath.Join(GetDataDir(), ProfilesDir)
}

func ProvidersRootDir() string {
	return filepath.Join(GetDataDir(), ProvidersDir)
}

func CacheDirPath() string {
	return filepath.Join(GetDataDir(), CacheDir)
}

func LogDirPath() string {
	return filepath.Join(GetDataDir(), LogsDir)
}

func StateFilePath() string {
	return filepath.Join(GetDataDir(), StateFile)
}

func ProfilesFilePath() string {
	return filepath.Join(GetDataDir(), ProfilesFile)
}

func ConfigFilePath() string {
	return filepath.Join(GetDataDir(), ConfigFile)
}

func ProfileFilePath(id int64) string {
	return filepath.Join(ProfileDir(), fmt.Sprintf("%d.yaml", id))
}

func ProviderDirForProfile(profileID int64) string {
	return filepath.Join(ProvidersRootDir(), fmt.Sprintf("%d", profileID))
}

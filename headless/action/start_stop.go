package action

import (
	"fmt"
	"os"

	"flclash-headless/app"
	"flclash-headless/configbuilder"
	"flclash-headless/coreclient"
	"flclash-headless/model"
	"flclash-headless/storage"
)

func StartCore(a *app.App) error {
	manifest := a.ProfileStore.GetManifest()
	profile := manifest.GetCurrentProfile()
	if profile == nil {
		return fmt.Errorf("当前没有可用配置，无法启动核心")
	}

	fmt.Println("  正在准备运行配置...")
	prefs := a.StateStore.Get()

	built, err := configbuilder.Build(profile, prefs)
	if err != nil {
		return fmt.Errorf("配置构建失败: %w", err)
	}
	_ = built

	corePath := findCorePath()
	if corePath == "" {
		return fmt.Errorf("未找到 FlClashCore 二进制文件")
	}
	a.InitCoreClient(corePath)

	fmt.Println("  正在启动核心...")
	if err := a.CoreClient.Start(); err != nil {
		return fmt.Errorf("核心启动失败: %w", err)
	}

	fmt.Println("  正在初始化核心...")
	if err := a.CoreClient.InitClash(1); err != nil {
		a.CoreClient.Stop()
		return fmt.Errorf("核心初始化失败: %w", err)
	}

	fmt.Println("  正在应用配置...")
	if _, err := a.CoreClient.ValidateConfig(built.Path); err != nil {
		a.CoreClient.Stop()
		return fmt.Errorf("配置校验失败: %w", err)
	}

	if err := a.CoreClient.SetupConfig(prefs.SelectedMap, ""); err != nil {
		a.CoreClient.Stop()
		return fmt.Errorf("配置应用失败: %w", err)
	}

	fmt.Println("  正在启动监听...")
	if err := a.CoreClient.StartListener(); err != nil {
		a.CoreClient.Stop()
		return fmt.Errorf("监听启动失败: %w", err)
	}

	a.CoreClient.AddListener(func(event coreclient.CoreEvent) {
		if event.Type == coreclient.EventLog {
			if data, ok := event.Data.(map[string]interface{}); ok {
				logType, _ := data["type"].(string)
				payload, _ := data["payload"].(string)
				a.State.AddLog(model.LogEntry{
					Type:    logType,
					Payload: payload,
				})
			}
		}
	})
	a.CoreClient.StartLog()
	_ = EnsureGlobalProxySelected(a)
	a.State.SetCoreStatus(coreclient.StatusRunning)
	a.State.CoreStartTime = model.TimeNow()
	a.StateStore.SetLastRunning(true)

	return nil
}

func StopCore(a *app.App) error {
	fmt.Println("  正在停止监听...")
	a.CoreClient.StopLog()
	a.CoreClient.StopListener()

	fmt.Println("  正在停止核心...")
	a.CoreClient.Shutdown()
	a.CoreClient.Stop()

	a.State.SetCoreStatus(coreclient.StatusStopped)
	a.StateStore.SetLastRunning(false)

	return nil
}

func findCorePath() string {
	candidates := []string{
		"FlClashCore",
		"./FlClashCore",
		"../FlClashCore",
		"/usr/local/bin/FlClashCore",
		"/usr/bin/FlClashCore",
		storage.GetDataDir() + "/FlClashCore",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

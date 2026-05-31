package action

import (
	"fmt"

	"flclash-headless/app"
	"flclash-headless/configbuilder"
	"flclash-headless/coreclient"
	"flclash-headless/model"
)

func ApplyProfile(a *app.App, profile *model.ProfileRecord) error {
	fmt.Println("  正在准备运行配置...")

	prefs := a.StateStore.Get()
	built, err := configbuilder.Build(profile, prefs)
	if err != nil {
		return fmt.Errorf("配置构建失败: %w", err)
	}
	_ = built

	fmt.Println("  正在应用配置...")
	if _, err := a.CoreClient.ValidateConfig(built.Path); err != nil {
		return fmt.Errorf("配置校验失败: %w", err)
	}

	if err := a.CoreClient.SetupConfig(prefs.SelectedMap, ""); err != nil {
		return fmt.Errorf("配置应用失败: %w", err)
	}

	_ = EnsureGlobalProxySelected(a)

	a.StateStore.SetCurrentProfileID(profile.ID)
	a.ProfileStore.SetCurrent(profile.ID)

	if a.State.GetCoreStatus() == coreclient.StatusRunning {
		fmt.Println("  正在重载核心...")
	}

	return nil
}

func SwitchProfile(a *app.App, profile *model.ProfileRecord) error {
	if profile == nil {
		return fmt.Errorf("配置不存在")
	}

	fmt.Printf("  正在应用配置: %s\n", profile.Name)
	fmt.Println("  正在检查配置...")

	if err := ApplyProfile(a, profile); err != nil {
		return fmt.Errorf("切换失败: %w", err)
	}

	fmt.Println("  完成。")
	return nil
}

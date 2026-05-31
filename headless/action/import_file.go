package action

import (
	"fmt"
	"os"
	"time"

	"flclash-headless/configbuilder"
	"flclash-headless/model"
	"flclash-headless/storage"
)

func ImportFromFile(profileStore *storage.ProfileStore, filePath, name string, autoApply bool) (*model.ProfileRecord, error) {
	fmt.Println("  正在读取文件...")

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, fmt.Errorf("无法读取文件: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("指定路径是一个目录，不是普通文件")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败 (可能无权限): %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("配置文件为空")
	}

	fmt.Println("  正在校验配置...")
	if err := configbuilder.ValidateRawYAML(data); err != nil {
		return nil, fmt.Errorf("配置校验未通过: %w", err)
	}

	now := time.Now()
	id := now.UnixMilli()
	profilePath, err := configbuilder.WriteRawYAML(id, data)
	if err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}

	if name == "" {
		name = info.Name()
	}

	profile := model.ProfileRecord{
		ID:        id,
		Name:      name,
		Type:      model.ProfileTypeFile,
		Source:    filePath,
		FilePath:  profilePath,
		CreatedAt: now,
		UpdatedAt: now,
		AutoApplyAfterImport: autoApply,
	}

	profileStore.AddProfile(profile)

	return &profile, nil
}

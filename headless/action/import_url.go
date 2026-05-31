package action

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	"flclash-headless/configbuilder"
	"flclash-headless/model"
	"flclash-headless/storage"
	"flclash-headless/util"
)

func ImportFromURL(profileStore *storage.ProfileStore, url, name string, autoApply bool) (*model.ProfileRecord, error) {
	return ImportFromURLWithProxy(profileStore, url, name, autoApply, "")
}

func ImportFromURLWithProxy(profileStore *storage.ProfileStore, url, name string, autoApply bool, proxyURL string) (*model.ProfileRecord, error) {
	fmt.Println("  正在下载订阅...")

	body, err := downloadProfileURL(url, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("下载失败: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("配置内容为空")
	}

	fmt.Println("  正在校验配置...")
	if err := configbuilder.ValidateRawYAML(body); err != nil {
		return nil, fmt.Errorf("配置校验未通过: %w", err)
	}

	now := time.Now()
	id := now.UnixMilli()
	profilePath, err := configbuilder.WriteRawYAML(id, body)
	if err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}

	if name == "" {
		name = util.AutoName(url)
	}

	profile := model.ProfileRecord{
		ID:                   id,
		Name:                 name,
		Type:                 model.ProfileTypeURL,
		Source:               url,
		FilePath:             profilePath,
		CreatedAt:            now,
		UpdatedAt:            now,
		AutoApplyAfterImport: autoApply,
	}

	profileStore.AddProfile(profile)
	fmt.Println("  正在保存配置信息...")

	return &profile, nil
}

func downloadProfileURL(sourceURL, proxyURL string) ([]byte, error) {
	clients := make([]*http.Client, 0, 2)
	if proxyURL != "" {
		client, err := newProxyHTTPClient(proxyURL)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	clients = append(clients, &http.Client{Timeout: 30 * time.Second})

	var lastErr error
	for _, client := range clients {
		body, err := downloadProfileURLWithClient(client, sourceURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func newProxyHTTPClient(proxyURL string) (*http.Client, error) {
	parsed, err := neturl.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("代理地址无效: %w", err)
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyURL(parsed)
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}, nil
}

func downloadProfileURLWithClient(client *http.Client, sourceURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "FlClash")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器返回 %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	return body, nil
}

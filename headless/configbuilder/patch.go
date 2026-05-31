package configbuilder

type PatchOptions struct {
	Mode               string
	MixedPort          int
	ExternalController string
	LogLevel           string
	TunEnabled         bool
	ProfileID          int64
}

func PatchConfig(config map[string]interface{}, opts PatchOptions) map[string]interface{} {
	if config == nil {
		config = make(map[string]interface{})
	}

	config["mixed-port"] = opts.MixedPort
	config["external-controller"] = opts.ExternalController
	config["allow-lan"] = true
	config["mode"] = opts.Mode
	config["log-level"] = opts.LogLevel
	config["ipv6"] = false
	config["tcp-concurrent"] = true
	config["unified-delay"] = true
	config["find-process-mode"] = "always"
	config["keep-alive-interval"] = 30

	tun := map[string]interface{}{}
	if existingTun, ok := config["tun"].(map[string]interface{}); ok {
		for k, v := range existingTun {
			tun[k] = v
		}
	}
	tun["enable"] = opts.TunEnabled
	tun["device"] = tunDeviceName
	tun["stack"] = "system"
	tun["auto-route"] = true
	tun["auto-detect-interface"] = true
	tun["strict-route"] = true
	tun["dns-hijack"] = []string{"any:53"}
	config["tun"] = tun

	if _, ok := config["dns"]; !ok {
		config["dns"] = map[string]interface{}{
			"enable":             true,
			"listen":             "0.0.0.0:1053",
			"enhanced-mode":      "fake-ip",
			"fake-ip-range":      "198.18.0.1/16",
			"default-nameserver": []string{"223.5.5.5"},
			"nameserver":         []string{"https://doh.pub/dns-query", "https://dns.alidns.com/dns-query"},
			"fallback":           []string{"tls://8.8.4.4", "tls://1.1.1.1"},
			"fake-ip-filter":     []string{"*.lan", "localhost.ptlogin2.qq.com"},
		}
	}

	if _, ok := config["geox-url"]; !ok {
		config["geox-url"] = map[string]interface{}{
			"mmdb":    "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.metadb",
			"geoip":   "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat",
			"geosite": "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat",
			"asn":     "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/GeoLite2-ASN.mmdb",
		}
	}

	return config
}

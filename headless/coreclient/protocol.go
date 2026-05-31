package coreclient

import "encoding/json"

type ActionMethod string

const (
	ActionInitClash          ActionMethod = "initClash"
	ActionGetIsInit          ActionMethod = "getIsInit"
	ActionShutdown           ActionMethod = "shutdown"
	ActionValidateConfig     ActionMethod = "validateConfig"
	ActionSetupConfig        ActionMethod = "setupConfig"
	ActionUpdateConfig       ActionMethod = "updateConfig"
	ActionGetConfig          ActionMethod = "getConfig"
	ActionGetProxies         ActionMethod = "getProxies"
	ActionChangeProxy        ActionMethod = "changeProxy"
	ActionGetTraffic         ActionMethod = "getTraffic"
	ActionGetTotalTraffic    ActionMethod = "getTotalTraffic"
	ActionResetTraffic       ActionMethod = "resetTraffic"
	ActionStartLog           ActionMethod = "startLog"
	ActionStopLog            ActionMethod = "stopLog"
	ActionStartListener      ActionMethod = "startListener"
	ActionStopListener       ActionMethod = "stopListener"
	ActionGetConnections     ActionMethod = "getConnections"
	ActionCloseConnections   ActionMethod = "closeConnections"
	ActionResetConnections   ActionMethod = "resetConnections"
	ActionGetExternalProviders ActionMethod = "getExternalProviders"
	ActionUpdateExternalProvider ActionMethod = "updateExternalProvider"
	ActionAsyncTestDelay     ActionMethod = "asyncTestDelay"
	ActionGetCountryCode     ActionMethod = "getCountryCode"
	ActionGetMemory          ActionMethod = "getMemory"
	ActionForceGC            ActionMethod = "forceGc"
	ActionCrash              ActionMethod = "crash"
	ActionDeleteFile         ActionMethod = "deleteFile"
	ActionCloseConnection    ActionMethod = "closeConnection"
	ActionSideLoadExternalProvider ActionMethod = "sideLoadExternalProvider"
	ActionUpdateGeoData      ActionMethod = "updateGeoData"
)

type ResultType int

const (
	ResultSuccess ResultType = 0
	ResultError   ResultType = -1
)

type Action struct {
	ID     string        `json:"id"`
	Method ActionMethod  `json:"method"`
	Data   interface{}   `json:"data"`
}

type ActionResult struct {
	ID     string      `json:"id"`
	Method ActionMethod `json:"method"`
	Data   interface{} `json:"data"`
	Code   ResultType  `json:"code"`
}

type SetupParams struct {
	SelectedMap map[string]string `json:"selected-map"`
	TestUrl     string            `json:"test-url"`
}

type InitParams struct {
	HomeDir string `json:"home-dir"`
	Version int    `json:"version"`
}

type ChangeProxyParams struct {
	GroupName string `json:"group-name"`
	ProxyName string `json:"proxy-name"`
}

type UpdateParams struct {
	Tun                *TunParams `json:"tun,omitempty"`
	AllowLan           *bool      `json:"allow-lan,omitempty"`
	MixedPort          *int       `json:"mixed-port,omitempty"`
	FindProcessMode    *string    `json:"find-process-mode,omitempty"`
	Mode               *string    `json:"mode,omitempty"`
	LogLevel           *string    `json:"log-level,omitempty"`
	IPv6               *bool      `json:"ipv6,omitempty"`
	TCPConcurrent      *bool      `json:"tcp-concurrent,omitempty"`
	ExternalController *string    `json:"external-controller,omitempty"`
	UnifiedDelay       *bool      `json:"unified-delay,omitempty"`
}

type TunParams struct {
	Enable       bool      `json:"enable"`
	Device       *string   `json:"device,omitempty"`
	Stack        *string   `json:"stack,omitempty"`
	DNSHijack    *[]string `json:"dns-hijack,omitempty"`
	AutoRoute    *bool     `json:"auto-route,omitempty"`
	RouteAddress *[]string `json:"route-address,omitempty"`
}

type DelayParams struct {
	ProxyName string `json:"proxy-name"`
	Timeout   int    `json:"timeout"`
	TestURL   string `json:"test-url"`
}

type ProxiesData struct {
	Proxies map[string]interface{} `json:"proxies"`
	All     []string               `json:"all"`
}

func NewAction(method ActionMethod, data interface{}) Action {
	return Action{
		ID:     string(method),
		Method: method,
		Data:   data,
	}
}

func MarshalAction(a Action) ([]byte, error) {
	return json.Marshal(a)
}

func UnmarshalResult(data []byte) (*ActionResult, error) {
	var result ActionResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

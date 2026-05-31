package coreclient

type ActionMethod string

const (
	ActionInitClash      ActionMethod = "initClash"
	ActionShutdown       ActionMethod = "shutdown"
	ActionValidateConfig ActionMethod = "validateConfig"
	ActionSetupConfig    ActionMethod = "setupConfig"
	ActionGetConfig      ActionMethod = "getConfig"
	ActionGetProxies     ActionMethod = "getProxies"
	ActionChangeProxy    ActionMethod = "changeProxy"
	ActionGetTraffic     ActionMethod = "getTraffic"
	ActionGetTotalTraffic ActionMethod = "getTotalTraffic"
	ActionStartLog       ActionMethod = "startLog"
	ActionStopLog        ActionMethod = "stopLog"
	ActionStartListener  ActionMethod = "startListener"
	ActionStopListener   ActionMethod = "stopListener"
	ActionAsyncTestDelay ActionMethod = "asyncTestDelay"
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

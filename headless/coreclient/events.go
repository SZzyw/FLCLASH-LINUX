package coreclient

type CoreEventType string

const (
	EventLog     CoreEventType = "log"
	EventDelay   CoreEventType = "delay"
	EventRequest CoreEventType = "request"
	EventLoaded  CoreEventType = "loaded"
	EventCrash   CoreEventType = "crash"
)

type CoreEvent struct {
	Type CoreEventType `json:"type"`
	Data interface{}   `json:"data"`
}

type LogEvent struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type DelayEvent struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Value int    `json:"value"`
}

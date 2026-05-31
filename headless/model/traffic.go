package model

type TrafficSnapshot struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

type TotalTraffic struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

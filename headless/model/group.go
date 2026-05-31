package model

type GroupSummary struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Now         string `json:"now"`
	Total       int    `json:"total"`
	TestUrl     string `json:"test_url"`
}

type GroupNode struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Delay int    `json:"delay"`
	Now   bool   `json:"now"`
}

type GroupDetail struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Now   string      `json:"now"`
	Nodes []GroupNode `json:"nodes"`
}

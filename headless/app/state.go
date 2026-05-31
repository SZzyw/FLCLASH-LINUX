package app

import (
	"sync"
	"time"

	"flclash-headless/coreclient"
	"flclash-headless/model"
)

type RuntimeState struct {
	mu            sync.RWMutex
	CoreStatus    coreclient.CoreStatus
	CoreConnected bool
	CoreStartTime time.Time
	Traffic       model.TrafficSnapshot
	TotalTraffic  model.TotalTraffic
	Groups        []model.GroupSummary
	GroupDetails  map[string]*model.GroupDetail
	Logs          []model.LogEntry
	LastError     string
	MessageText   string
}

func NewRuntimeState() *RuntimeState {
	return &RuntimeState{
		Groups:       make([]model.GroupSummary, 0),
		GroupDetails: make(map[string]*model.GroupDetail),
		Logs:         make([]model.LogEntry, 0, 500),
	}
}

func (rs *RuntimeState) GetCoreStatus() coreclient.CoreStatus {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.CoreStatus
}

func (rs *RuntimeState) SetCoreStatus(s coreclient.CoreStatus) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.CoreStatus = s
}

func (rs *RuntimeState) IsCoreRunning() bool {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.CoreConnected && rs.CoreStatus == coreclient.StatusRunning
}

func (rs *RuntimeState) AddLog(entry model.LogEntry) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.Logs = append(rs.Logs, entry)
	if len(rs.Logs) > 500 {
		rs.Logs = rs.Logs[len(rs.Logs)-500:]
	}
}

func (rs *RuntimeState) GetLogs() []model.LogEntry {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	result := make([]model.LogEntry, len(rs.Logs))
	copy(result, rs.Logs)
	return result
}

func (rs *RuntimeState) SetGroups(groups []model.GroupSummary) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.Groups = groups
}

func (rs *RuntimeState) GetGroups() []model.GroupSummary {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	result := make([]model.GroupSummary, len(rs.Groups))
	copy(result, rs.Groups)
	return result
}

func (rs *RuntimeState) SetGroupDetail(detail *model.GroupDetail) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.GroupDetails[detail.Name] = detail
}

func (rs *RuntimeState) GetGroupDetail(name string) *model.GroupDetail {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.GroupDetails[name]
}

func (rs *RuntimeState) SetTraffic(t model.TrafficSnapshot) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.Traffic = t
}

func (rs *RuntimeState) GetTraffic() model.TrafficSnapshot {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.Traffic
}

func (rs *RuntimeState) SetTotalTraffic(t model.TotalTraffic) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.TotalTraffic = t
}

func (rs *RuntimeState) GetTotalTraffic() model.TotalTraffic {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.TotalTraffic
}

func (rs *RuntimeState) SetMessage(text string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.MessageText = text
}

func (rs *RuntimeState) GetMessage() string {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.MessageText
}

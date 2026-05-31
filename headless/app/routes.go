package app

import "sync"

type Route string

const (
	RouteDashboard   Route = "dashboard"
	RouteProfiles    Route = "profiles"
	RouteGroupList   Route = "group_list"
	RouteGroupDetail Route = "group_detail"
	RouteLogs        Route = "logs"
	RouteTraffic     Route = "traffic"
	RouteImportURL   Route = "import_url"
	RouteImportFile  Route = "import_file"
	RouteModeSwitch  Route = "mode_switch"
	RouteConfigDetail Route = "config_detail"
	RouteGlobalExit   Route = "global_exit"
	RouteMessage     Route = "message"
)

type PageStack struct {
	mu    sync.RWMutex
	pages []Route
}

func NewPageStack() *PageStack {
	return &PageStack{
		pages: []Route{RouteDashboard},
	}
}

func (s *PageStack) Current() Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.pages) == 0 {
		return RouteDashboard
	}
	return s.pages[len(s.pages)-1]
}

func (s *PageStack) Push(r Route) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages = append(s.pages, r)
}

func (s *PageStack) Pop() Route {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pages) <= 1 {
		return RouteDashboard
	}
	s.pages = s.pages[:len(s.pages)-1]
	if len(s.pages) == 0 {
		return RouteDashboard
	}
	return s.pages[len(s.pages)-1]
}

func (s *PageStack) PopTo(r Route) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for len(s.pages) > 1 && s.pages[len(s.pages)-1] != r {
		if len(s.pages) <= 1 {
			return
		}
		s.pages = s.pages[:len(s.pages)-1]
	}
}

func (s *PageStack) Replace(r Route) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pages) > 0 {
		s.pages[len(s.pages)-1] = r
	} else {
		s.pages = append(s.pages, r)
	}
}

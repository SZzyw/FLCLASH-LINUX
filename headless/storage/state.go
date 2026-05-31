package storage

import (
	"encoding/json"
	"os"
	"sync"

	"flclash-headless/model"
)

type StateStore struct {
	mu     sync.RWMutex
	prefs  *model.RuntimePrefs
	dirty  bool
}

func NewStateStore() *StateStore {
	return &StateStore{
		prefs: model.NewRuntimePrefs(),
	}
}

func (s *StateStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(StateFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			s.prefs = model.NewRuntimePrefs()
			return nil
		}
		return err
	}

	var prefs model.RuntimePrefs
	if err := json.Unmarshal(data, &prefs); err != nil {
		s.prefs = model.NewRuntimePrefs()
		return nil
	}
	if prefs.SelectedMap == nil {
		prefs.SelectedMap = make(map[string]string)
	}
	if prefs.Mode == "" {
		prefs.Mode = model.ModeRule
	}
	if prefs.MixedPort == 0 {
		prefs.MixedPort = 7890
	}
	if prefs.ExternalController == "" {
		prefs.ExternalController = "127.0.0.1:9090"
	}
	if prefs.LogLevel == "" {
		prefs.LogLevel = "info"
	}
	s.prefs = &prefs
	return nil
}

func (s *StateStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.prefs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StateFilePath(), data, 0644)
}

func (s *StateStore) SaveIfDirty() error {
	s.mu.RLock()
	if !s.dirty {
		s.mu.RUnlock()
		return nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.dirty = false
	s.mu.Unlock()

	return s.Save()
}

func (s *StateStore) Get() *model.RuntimePrefs {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.prefs
}

func (s *StateStore) SetCurrentProfileID(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.CurrentProfileID = id
	s.dirty = true
}

func (s *StateStore) SetMode(mode model.Mode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.Mode = mode
	s.dirty = true
}

func (s *StateStore) SetTunEnabled(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.TunEnabled = enabled
	s.dirty = true
}

func (s *StateStore) SetSystemProxy(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.SystemProxy = enabled
	s.dirty = true
}

func (s *StateStore) SetSelectedMap(selectedMap map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.SelectedMap = selectedMap
	s.dirty = true
}

func (s *StateStore) SetLastRunning(running bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefs.LastRunning = running
	s.dirty = true
}

package storage

import (
	"encoding/json"
	"os"
	"sync"

	"flclash-headless/model"
)

type ProfileStore struct {
	mu       sync.RWMutex
	manifest *model.ProfilesManifest
}

func NewProfileStore() *ProfileStore {
	return &ProfileStore{
		manifest: &model.ProfilesManifest{
			Profiles: []model.ProfileRecord{},
		},
	}
}

func (ps *ProfileStore) Load() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	data, err := os.ReadFile(ProfilesFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			ps.manifest = &model.ProfilesManifest{
				Profiles: []model.ProfileRecord{},
			}
			return nil
		}
		return err
	}

	var manifest model.ProfilesManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		ps.manifest = &model.ProfilesManifest{
			Profiles: []model.ProfileRecord{},
		}
		return nil
	}
	if manifest.Profiles == nil {
		manifest.Profiles = []model.ProfileRecord{}
	}
	ps.manifest = &manifest
	return nil
}

func (ps *ProfileStore) Save() error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	data, err := json.MarshalIndent(ps.manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ProfilesFilePath(), data, 0644)
}

func (ps *ProfileStore) GetManifest() *model.ProfilesManifest {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.manifest
}

func (ps *ProfileStore) AddProfile(p model.ProfileRecord) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.manifest.AddProfile(p)
}

func (ps *ProfileStore) RemoveProfile(id int64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.manifest.RemoveProfile(id)
}

func (ps *ProfileStore) SetCurrent(id int64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.manifest.SetCurrent(id)
}

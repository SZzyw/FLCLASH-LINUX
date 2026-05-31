package model

import "time"

type ProfileType string

const (
	ProfileTypeURL  ProfileType = "url"
	ProfileTypeFile ProfileType = "file"
)

type ProfileRecord struct {
	ID              int64       `json:"id"`
	Name            string      `json:"name"`
	Type            ProfileType `json:"type"`
	Source          string      `json:"source"`
	FilePath        string      `json:"file_path"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	AutoApplyAfterImport bool   `json:"auto_apply_after_import"`
}

type ProfilesManifest struct {
	CurrentProfileID int64           `json:"current_profile_id"`
	Profiles         []ProfileRecord `json:"profiles"`
}

func (m *ProfilesManifest) GetProfile(id int64) *ProfileRecord {
	for i := range m.Profiles {
		if m.Profiles[i].ID == id {
			return &m.Profiles[i]
		}
	}
	return nil
}

func (m *ProfilesManifest) GetCurrentProfile() *ProfileRecord {
	if m.CurrentProfileID == 0 {
		return nil
	}
	return m.GetProfile(m.CurrentProfileID)
}

func (m *ProfilesManifest) AddProfile(p ProfileRecord) {
	m.Profiles = append(m.Profiles, p)
}

func (m *ProfilesManifest) RemoveProfile(id int64) {
	for i := range m.Profiles {
		if m.Profiles[i].ID == id {
			m.Profiles = append(m.Profiles[:i], m.Profiles[i+1:]...)
			if m.CurrentProfileID == id {
				m.CurrentProfileID = 0
			}
			return
		}
	}
}

func (m *ProfilesManifest) SetCurrent(id int64) {
	m.CurrentProfileID = id
}

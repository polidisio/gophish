package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
)

type EntraIDSettings struct {
	ID               int64     `json:"id" gorm:"column:id; primary_key:yes"`
	Enabled          bool      `json:"enabled"`
	ClientID         string    `json:"client_id"`
	ClientSecret     string    `json:"client_secret"`
	TenantID         string    `json:"tenant_id"`
	RedirectURI      string    `json:"redirect_uri"`
	Scopes           string    `json:"scopes"`
	AdminOnly        bool      `json:"admin_only"`
	AutoCreateUsers  bool      `json:"auto_create_users"`
	SyncDepartments  bool      `json:"sync_departments"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (e *EntraIDSettings) Validate() error {
	if e.Enabled {
		if e.ClientID == "" {
			return errors.New("Client ID is required")
		}
		if e.ClientSecret == "" {
			return errors.New("Client Secret is required")
		}
		if e.TenantID == "" {
			return errors.New("Tenant ID is required")
		}
		if e.RedirectURI == "" {
			return errors.New("Redirect URI is required")
		}
	}
	return nil
}

func GetEntraIDSettings() (*EntraIDSettings, error) {
	s := &EntraIDSettings{}
	err := db.First(s).Error
	if err != nil {
		if err.Error() == "record not found" {
			s.Scopes = "openid profile email User.Read"
			return s, nil
		}
		return nil, err
	}
	return s, nil
}

func SaveEntraIDSettings(s *EntraIDSettings) error {
	err := s.Validate()
	if err != nil {
		log.Error(err)
		return err
	}

	existing, err := GetEntraIDSettings()
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if existing.ID == 0 {
		s.CreatedAt = time.Now()
		s.UpdatedAt = time.Now()
		err = db.Create(s).Error
	} else {
		s.ID = existing.ID
		s.CreatedAt = existing.CreatedAt
		s.UpdatedAt = time.Now()
		err = db.Save(s).Error
	}
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

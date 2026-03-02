package models

import (
	"time"
)

// OAuthProvider represents an OAuth2/OIDC provider configuration
type OAuthProvider struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"` // e.g., "Microsoft", "Google"
	Provider    string    `json:"provider"` // "azuread", "google", "okta"
	ClientID    string    `json:"client_id"`
	ClientSecret string   `json:"-"`
	TenantID    string    `json:"tenant_id"` // Azure AD Tenant ID
	RedirectURL string    `json:"redirect_url"`
	Scopes      string    `json:"scopes"` // comma-separated
	Enabled     bool     `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetScopes returns the scopes as a slice
func (o *OAuthProvider) GetScopes() []string {
	return []string{"openid", "profile", "email", "User.Read"}
}

// AzureADUser represents a user synced from Azure AD
type AzureADUser struct {
	ID            string    `json:"id" gorm:"primary_key"` // Azure AD Object ID
	Email         string    `json:"email" gorm:"unique_index"`
	DisplayName   string    `json:"display_name"`
	GivenName     string    `json:"given_name"`
	Surname       string    `json:"surname"`
	JobTitle      string    `json:"job_title"`
	Department    string    `json:"department"`
	OfficeLocation string  `json:"office_location"`
	Photo         []byte    `json:"photo"`
	ManagerID     string    `json:"manager_id"`
	LastSyncedAt  time.Time `json:"last_synced_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SyncStatus represents the status of an Azure AD sync
type SyncStatus struct {
	ID             int64     `json:"id" gorm:"primary_key"`
	ProviderID     int64     `json:"provider_id"`
	StartedAt      time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	TotalUsers     int       `json:"total_users"`
	NewUsers       int       `json:"new_users"`
	UpdatedUsers   int       `json:"updated_users"`
	DeletedUsers   int       `json:"deleted_users"`
	FailedUsers    int       `json:"failed_users"`
	Status        string    `json:"status"` // "running", "completed", "failed"
	ErrorMessage  string    `json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
}

// IsComplete returns true if the sync is complete
func (s *SyncStatus) IsComplete() bool {
	return s.CompletedAt != nil
}

// GetOAuthProvider returns the OAuth provider by ID
func GetOAuthProvider(id int64) (OAuthProvider, error) {
	var provider OAuthProvider
	err := db.Where("id = ?", id).First(&provider).Error
	return provider, err
}

// GetOAuthProviders returns all OAuth providers
func GetOAuthProviders() ([]OAuthProvider, error) {
	var providers []OAuthProvider
	err := db.Find(&providers).Error
	return providers, err
}

// GetEnabledOAuthProvider returns the enabled OAuth provider
func GetEnabledOAuthProvider() (OAuthProvider, error) {
	var provider OAuthProvider
	err := db.Where("enabled = ?", true).First(&provider).Error
	return provider, err
}

// SaveOAuthProvider saves an OAuth provider
func SaveOAuthProvider(provider *OAuthProvider) error {
	return db.Save(provider).Error
}

// DeleteOAuthProvider deletes an OAuth provider
func DeleteOAuthProvider(id int64) error {
	return db.Delete(&OAuthProvider{}, "id = ?", id).Error
}

// GetAzureADUser returns a user by Azure AD ID
func GetAzureADUser(id string) (AzureADUser, error) {
	var user AzureADUser
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

// GetAzureADUserByEmail returns a user by email
func GetAzureADUserByEmail(email string) (AzureADUser, error) {
	var user AzureADUser
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

// GetAzureADUsers returns all synced Azure AD users
func GetAzureADUsers() ([]AzureADUser, error) {
	var users []AzureADUser
	err := db.Find(&users).Error
	return users, err
}

// SaveAzureADUser saves or updates an Azure AD user
func SaveAzureADUser(user *AzureADUser) error {
	existing, err := GetAzureADUser(user.ID)
	if err == nil {
		// Update existing
		user.ID = existing.ID
		user.CreatedAt = existing.CreatedAt
	}
	user.UpdatedAt = time.Now()
	return db.Save(user).Error
}

// DeleteAzureADUser deletes an Azure AD user
func DeleteAzureADUser(id string) error {
	return db.Delete(&AzureADUser{}, "id = ?", id).Error
}

// CreateSyncStatus creates a new sync status record
func CreateSyncStatus(providerID int64) (SyncStatus, error) {
	sync := SyncStatus{
		ProviderID: providerID,
		StartedAt:  time.Now(),
		Status:     "running",
		CreatedAt:  time.Now(),
	}
	err := db.Save(&sync).Error
	return sync, err
}

// UpdateSyncStatus updates a sync status record
func UpdateSyncStatus(sync *SyncStatus) error {
	return db.Save(sync).Error
}

// GetSyncStatuses returns sync history for a provider
func GetSyncStatuses(providerID int64, limit int) ([]SyncStatus, error) {
	var statuses []SyncStatus
	err := db.Where("provider_id = ?", providerID).Order("created_at DESC").Limit(limit).Find(&statuses).Error
	return statuses, err
}

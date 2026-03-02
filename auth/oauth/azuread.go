package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/csrf"
	"github.com/gorilla/m"
	"github.com/gorilla/sessions"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	// Session keys
	sessionKeyOAuthState = "oauth2_state"
	sessionKeyUserID     = "user_id"
	sessionKeyRedirect   = "redirect_url"

	// DefaultScopes for Azure AD
	DefaultScopes = []string{"openid", "profile", "email", "User.Read", "Directory.Read.All"}
)

// AzureADConfig holds the Azure AD OAuth2 configuration
type AzureADConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	RedirectURL  string
	Scopes       []string
}

// AzureADProvider handles Azure AD OAuth2 authentication
type AzureADProvider struct {
	config      AzureADConfig
	verifier    *oidc.IDTokenVerifier
	provider    *oidc.Provider
	oauth2Cfg   oauth2.Config
	sessionStore sessions.Store
}

// NewAzureADProvider creates a new Azure AD OAuth provider
func NewAzureADProvider(cfg AzureADConfig, sessionStore sessions.Store) (*AzureADProvider, error) {
	ctx := context.Background()

	// Determine tenant URL
	tenantURL := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", cfg.TenantID)

	provider, err := oidc.NewProvider(ctx, tenantURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Create verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	// Create OAuth2 config
	oauth2Cfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     microsoft.Endpoint(cfg.TenantID),
		Scopes:       cfg.Scopes,
	}

	return &AzureADProvider{
		config:      cfg,
		provider:    provider,
		verifier:    verifier,
		oauth2Cfg:   oauth2Cfg,
		sessionStore: sessionStore,
	}, nil
}

// generateState generates a random state string for CSRF protection
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Login initiates the Azure AD OAuth2 flow
func (a *AzureADProvider) Login(w http.ResponseWriter, r *http.Request) {
	// Generate state
	state, err := generateState()
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=state_generation_failed", http.StatusTemporaryRedirect)
		return
	}

	// Save state in session
	session, _ := a.sessionStore.Get(r, "gophish")
	session.Values[sessionKeyOAuthState] = state
	if err := session.Save(r, w); err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=session_save_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get redirect URL from query or default to dashboard
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}
	session.Values[sessionKeyRedirect] = redirect
	session.Save(r, w)

	// Generate auth code URL
	authURL := a.oauth2Cfg.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.ApprovalForce)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Callback handles the Azure AD OAuth2 callback
func (a *AzureADProvider) Callback(w http.ResponseWriter, r *http.Request) {
	// Get state from session
	session, err := a.sessionStore.Get(r, "gophish")
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=session_error", http.StatusTemporaryRedirect)
		return
	}

	// Verify state
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Error("missing state parameter")
		http.Redirect(w, r, "/login?error=missing_state", http.StatusTemporaryRedirect)
		return
	}

	expectedState, ok := session.Values[sessionKeyOAuthState].(string)
	if !ok || state != expectedState {
		log.Error("invalid state")
		http.Redirect(w, r, "/login?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Get error from callback
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Errorf("OAuth error: %s - %s", errMsg, r.URL.Query().Get("error_description"))
		http.Redirect(w, r, "/login?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Error("missing code parameter")
		http.Redirect(w, r, "/login?error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	token, err := a.oauth2Cfg.Exchange(r.Context(), code)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get user info from ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Error("no id_token in response")
		http.Redirect(w, r, "/login?error=no_id_token", http.StatusTemporaryRedirect)
		return
	}

	idToken, err := a.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=token_verification_failed", http.StatusTemporaryRedirect)
		return
	}

	// Extract claims
	var claims struct {
		Email         string `json:"email"`
		PreferredUsername string `json:"preferred_username"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		OID           string `json:"oid"`
		TID           string `json:"tid"`
	}

	if err := idToken.Claims(&claims); err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=claims_failed", http.StatusTemporaryRedirect)
		return
	}

	// Use email from claims or preferred_username
	email := claims.Email
	if email == "" {
		email = claims.PreferredUsername
	}

	// Find or create user in Gophish
	user, err := a.findOrCreateUser(email, claims.Name, claims.GivenName, claims.FamilyName)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=user_creation_failed", http.StatusTemporaryRedirect)
		return
	}

	// Save user ID in session
	session.Values[sessionKeyUserID] = user.Id
	delete(session.Values, sessionKeyOAuthState)

	// Get redirect URL
	redirect, _ := session.Values[sessionKeyRedirect].(string)
	if redirect == "" {
		redirect = "/"
	}

	if err := session.Save(r, w); err != nil {
		log.Error(err)
		http.Redirect(w, r, "/login?error=session_save_failed", http.StatusTemporaryRedirect)
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// findOrCreateUser finds or creates a user in Gophish
func (a *AzureADProvider) findOrCreateUser(email, name, givenName, familyName string) (*models.User, error) {
	// Try to find existing user
	user, err := models.GetUserByEmail(email)
	if err == nil {
		return user, nil
	}

	// Check if user exists with different email
	users, err := models.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}

	// Check if Azure AD user already synced
	azureUser, err := models.GetAzureADUserByEmail(email)
	if err == nil {
		// Link to existing Gophish user or create new one
		username := strings.Split(email, "@")[0]
		newUser := models.User{
			Username: username,
			Email:    email,
			ApiKey:   models.GenerateSecureKey(models.APIKeyLength),
		}
		err = models.PostUser(&newUser, 1)
		if err != nil {
			return nil, err
		}
		return &newUser, nil
	}

	// Create new user from Azure AD
	username := strings.Split(email, "@")[0]
	newUser := models.User{
		Username: username,
		Email:    email,
		ApiKey:   models.GenerateSecureKey(models.APIKeyLength),
	}

	err = models.PostUser(&newUser, 1)
	if err != nil {
		return nil, err
	}

	// Save Azure AD user info
	azureUser = models.AzureADUser{
		ID:          email, // Use email as ID if no OID
		Email:       email,
		DisplayName: name,
		GivenName:   givenName,
		Surname:     familyName,
	}
	models.SaveAzureADUser(&azureUser)

	return &newUser, nil
}

// SyncUsersFromAzureAD syncs users from Azure AD
func (a *AzureADProvider) SyncUsersFromAzureAD() error {
	ctx := context.Background()

	// Get access token
	token, err := a.oauth2Cfg.Exchange(ctx, "", oauth2.TokenSourceFunc(func(ctx context.Context) (*oauth2.Token, error) {
		// Use client credentials flow for syncing
		return &oauth2.Token{
			AccessToken: "dummy", // Will be replaced by on-behalf-of flow
		}, nil
	}))
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make API call to get users
	client := a.oauth2Cfg.Client(ctx, token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/users?$select=id,displayName,mail,userPrincipalName,jobTitle,department,officeLocation,manager&$top=999")
	if err != nil {
		return fmt.Errorf("failed to call Graph API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Graph API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var graphResponse struct {
		Value []struct {
			ID              string `json:"id"`
			DisplayName     string `json:"displayName"`
			Mail            string `json:"mail"`
			UserPrincipalName string `json:"userPrincipalName"`
			JobTitle        string `json:"jobTitle"`
			Department      string `json:"department"`
			OfficeLocation  string `json:"officeLocation"`
			Manager         struct {
				ID string `json:"id"`
			} `json:"manager"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&graphResponse); err != nil {
		return fmt.Errorf("failed to parse Graph API response: %w", err)
	}

	// Save users
	for _, u := range graphResponse.Value {
		email := u.Mail
		if email == "" {
			email = u.UserPrincipalName
		}

		azureUser := models.AzureADUser{
			ID:             u.ID,
			Email:          email,
			DisplayName:    u.DisplayName,
			JobTitle:       u.JobTitle,
			Department:     u.Department,
			OfficeLocation: u.OfficeLocation,
			ManagerID:      u.Manager.ID,
			LastSyncedAt:   time.Now(),
		}

		if err := models.SaveAzureADUser(&azureUser); err != nil {
			log.Warningf("Failed to save user %s: %v", email, err)
			continue
		}
	}

	return nil
}

// RequireLogin returns a middleware that requires Azure AD login
func RequireLogin(provider *AzureADProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := provider.sessionStore.Get(r, "gophish")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			userID, ok := session.Values[sessionKeyUserID].(int64)
			if !ok || userID == 0 {
				// Save current URL for redirect after login
				redirect := r.URL.Path
				if r.URL.RawQuery != "" {
					redirect += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, "/oauth/login?redirect="+redirect, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

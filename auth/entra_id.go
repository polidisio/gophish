package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/gophish/gophish/logger"
)

type EntraIDConfig struct {
	ClientID     string
	ClientSecret string
	TenantID    string
	RedirectURI  string
	Scopes      string
}

type EntraIDProvider struct {
	ClientID     string
	ClientSecret string
	TenantID    string
	RedirectURI  string
	Scopes      string
}

func NewEntraIDProvider(config *EntraIDConfig) *EntraIDProvider {
	scopes := config.Scopes
	if scopes == "" {
		scopes = "openid profile email User.Read"
	}
	return &EntraIDProvider{
		ClientID:    config.ClientID,
		ClientSecret: config.ClientSecret,
		TenantID:    config.TenantID,
		RedirectURI:  config.RedirectURI,
		Scopes:      scopes,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type MicrosoftUser struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	Department        string `json:"department"`
}

func (p *EntraIDProvider) GetAuthorizationURL(state string) string {
	return fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s&state=%s",
		p.TenantID,
		p.ClientID,
		p.RedirectURI,
		urlEncode(p.Scopes),
		state,
	)
}

func (p *EntraIDProvider) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	data := fmt.Sprintf(
		"client_id=%s&scope=%s&code=%s&redirect_uri=%s&grant_type=authorization_code&client_secret=%s",
		p.ClientID,
		urlEncode(p.Scopes),
		code,
		p.RedirectURI,
		urlEncode(p.ClientSecret),
	)

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", p.TenantID),
		strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Entra ID token exchange failed: ", string(body))
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var token TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (p *EntraIDProvider) GetUserInfo(ctx context.Context, accessToken string) (*MicrosoftUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://graph.microsoft.com/v1.0/me?$select=id,displayName,mail,userPrincipalName,department",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Microsoft Graph API error: ", string(body))
		return nil, fmt.Errorf("graph API request failed with status %d", resp.StatusCode)
	}

	var user MicrosoftUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *EntraIDProvider) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, *MicrosoftUser, error) {
	data := fmt.Sprintf(
		"client_id=%s&scope=%s&refresh_token=%s&grant_type=refresh_token&client_secret=%s",
		p.ClientID,
		urlEncode(p.Scopes),
		refreshToken,
		urlEncode(p.ClientSecret),
	)

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", p.TenantID),
		strings.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("refresh token failed with status %d", resp.StatusCode)
	}

	var token TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, nil, err
	}

	user, err := p.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, nil, err
	}

	return &token, user, nil
}

func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func urlEncode(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "+", "%2B"), " ", "%20")
}

type OAuthState struct {
	Code        string
	State       string
	ExpiresAt   time.Time
}

var OAuthStateStore = make(map[string]*OAuthState)

func StoreOAuthState(state string, stateData *OAuthState) {
	OAuthStateStore[state] = stateData
	go func() {
		time.Sleep(10 * time.Minute)
		delete(OAuthStateStore, state)
	}()
}

func GetOAuthState(state string) *OAuthState {
	return OAuthStateStore[state]
}

func DeleteOAuthState(state string) {
	delete(OAuthStateStore, state)
}

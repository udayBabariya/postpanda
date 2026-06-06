package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"postpanda/backend-go/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthUserInfo struct {
	ID    string
	Email string
	Name  string
}

// GitHub Provider

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     config.App.GitHubClientID,
			ClientSecret: config.App.GitHubClientSecret,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
			RedirectURL:  config.App.BackendURL + "/auth/github/callback",
		},
	}
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GitHubProvider) ExchangeCode(code string) (*OAuthUserInfo, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := p.config.Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &ghUser); err != nil {
		return nil, err
	}

	if ghUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			emailBody, _ := io.ReadAll(emailResp.Body)
			if json.Unmarshal(emailBody, &emails) == nil {
				for _, e := range emails {
					if e.Primary {
						ghUser.Email = e.Email
						break
					}
				}
			}
		}
	}

	return &OAuthUserInfo{
		ID:    fmt.Sprintf("%d", ghUser.ID),
		Email: ghUser.Email,
		Name:  ghUser.Name,
	}, nil
}

// Google Provider

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     config.App.GoogleClientID,
			ClientSecret: config.App.GoogleClientSecret,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
			RedirectURL:  config.App.BackendURL + "/auth/google/callback",
		},
	}
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GoogleProvider) ExchangeCode(code string) (*OAuthUserInfo, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := p.config.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &gUser); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		ID:    gUser.ID,
		Email: gUser.Email,
		Name:  gUser.Name,
	}, nil
}

// Generic OAuth Provider

type GenericOAuthProvider struct {
	httpClient *http.Client
}

func NewGenericOAuthProvider() *GenericOAuthProvider {
	return &GenericOAuthProvider{httpClient: &http.Client{}}
}

func (p *GenericOAuthProvider) GetAuthURL(state string) string {
	return fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&response_type=code",
		config.App.GenericOAuthURL,
		config.App.GenericOAuthClientID,
		config.App.BackendURL+"/auth/generic/callback",
		state,
	)
}

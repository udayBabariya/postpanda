package social

import (
	"fmt"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type LinkedInProvider struct{}

func init() {
	integrations.Register("linkedin", &LinkedInProvider{})
}

func (p *LinkedInProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{
		ID:          "linkedin",
		Name:        "LinkedIn",
		Type:        "social",
		Description: "Post to LinkedIn profile or company page",
		AuthType:    "oauth2",
	}
}

func (p *LinkedInProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf(
		"https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&scope=openid+profile+email+w_member_social&state=%s",
		config.App.LinkedInClientID,
		config.App.BackendURL+"/integrations/linkedin/callback",
		orgID,
	)
}

func (p *LinkedInProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Exchange code for access token
	// TODO: implement LinkedIn OAuth
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return &integrations.IntegrationInfo{
		InternalID:      "linkedin_user_id",
		Name:            "LinkedIn User",
		Type:            "profile",
		Token:           code,
		TokenExpiration: &expiry,
	}, nil
}

func (p *LinkedInProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *LinkedInProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Post to LinkedIn via REST API
	// TODO: implement
	return &integrations.PublishResult{
		PostID: "linkedin_post_id",
	}, nil
}

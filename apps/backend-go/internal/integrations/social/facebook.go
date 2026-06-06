package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type FacebookProvider struct{}
type InstagramProvider struct{}

func init() {
	integrations.Register("facebook", &FacebookProvider{})
	integrations.Register("instagram", &InstagramProvider{})
}

func (p *FacebookProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "facebook", Name: "Facebook", Type: "social", AuthType: "oauth2"}
}

func (p *FacebookProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/facebook/callback"
	return fmt.Sprintf(
		"https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=pages_manage_posts,pages_read_engagement,instagram_basic,instagram_content_publish&state=%s",
		clientID, redirectURI, orgID,
	)
}

func (p *FacebookProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "fb_page_id", Name: "Facebook Page", Type: "page", Token: code}, nil
}

func (p *FacebookProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *FacebookProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "fb_post_id"}, nil
}

func (p *InstagramProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "instagram", Name: "Instagram", Type: "social", AuthType: "oauth2"}
}

func (p *InstagramProvider) GetOAuthURL(orgID string) string {
	return (&FacebookProvider{}).GetOAuthURL(orgID)
}

func (p *InstagramProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "ig_account_id", Name: "Instagram Account", Type: "profile", Token: code}, nil
}

func (p *InstagramProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *InstagramProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "ig_media_id"}, nil
}

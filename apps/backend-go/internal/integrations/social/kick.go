package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type KickProvider struct{}

func init() {
	integrations.Register("kick", &KickProvider{})
}

func (p *KickProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "kick", Name: "Kick", Type: "social", Description: "Post to Kick streaming platform", AuthType: "oauth2"}
}

func (p *KickProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("KICK_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/kick/callback"
	return fmt.Sprintf("https://kick.com/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=channel:read+channel:write&state=%s", clientID, redirectURI, orgID)
}

func (p *KickProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "kick_channel_id", Name: "Kick Channel", Type: "channel", Token: code}, nil
}

func (p *KickProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(24 * time.Hour)
	return token, &expiry, nil
}

func (p *KickProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "kick_post_id"}, nil
}

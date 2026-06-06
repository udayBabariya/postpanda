package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type VKProvider struct{}

func init() {
	integrations.Register("vk", &VKProvider{})
}

func (p *VKProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "vk", Name: "VK", Type: "social", Description: "Post to VK (VKontakte)", AuthType: "oauth2"}
}

func (p *VKProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("VK_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/vk/callback"
	return fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&display=page&redirect_uri=%s&scope=wall,photos,offline&response_type=code&state=%s", clientID, redirectURI, orgID)
}

func (p *VKProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "vk_user_id", Name: "VK Account", Type: "profile", Token: code}, nil
}

func (p *VKProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *VKProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// VK Wall Post API
	return &integrations.PublishResult{PostID: "vk_post_id"}, nil
}

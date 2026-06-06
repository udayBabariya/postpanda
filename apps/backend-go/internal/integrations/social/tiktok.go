package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type TikTokProvider struct{}

func init() {
	integrations.Register("tiktok", &TikTokProvider{})
}

func (p *TikTokProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "tiktok", Name: "TikTok", Type: "social", AuthType: "oauth2"}
}

func (p *TikTokProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("TIKTOK_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/tiktok/callback"
	return fmt.Sprintf(
		"https://www.tiktok.com/v2/auth/authorize?response_type=code&client_key=%s&redirect_uri=%s&scope=video.publish,video.upload,user.info.basic&state=%s",
		clientID, redirectURI, orgID,
	)
}

func (p *TikTokProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "tiktok_user_id", Name: "TikTok Account", Type: "profile", Token: code}, nil
}

func (p *TikTokProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(24 * time.Hour)
	return token, &expiry, nil
}

func (p *TikTokProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "tiktok_video_id"}, nil
}

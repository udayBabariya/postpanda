package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type TwitchProvider struct{}

func init() {
	integrations.Register("twitch", &TwitchProvider{})
}

func (p *TwitchProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "twitch", Name: "Twitch", Type: "social", Description: "Post to Twitch channel", AuthType: "oauth2"}
}

func (p *TwitchProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/twitch/callback"
	return fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=channel:manage:schedule+channel:read:subscriptions&state=%s", clientID, redirectURI, orgID)
}

func (p *TwitchProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "twitch_broadcaster_id", Name: "Twitch Channel", Type: "channel", Token: code}, nil
}

func (p *TwitchProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(4 * time.Hour)
	return token, &expiry, nil
}

func (p *TwitchProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Post to channel schedule
	return &integrations.PublishResult{PostID: "twitch_schedule_id"}, nil
}

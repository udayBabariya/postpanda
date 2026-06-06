package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type DiscordProvider struct{}

func init() {
	integrations.Register("discord", &DiscordProvider{})
}

func (p *DiscordProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "discord", Name: "Discord", Type: "social", AuthType: "oauth2"}
}

func (p *DiscordProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/discord/callback"
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=webhook.incoming+bot&state=%s",
		clientID, redirectURI, orgID,
	)
}

func (p *DiscordProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "discord_channel_id", Name: "Discord Channel", Type: "social", Token: code}, nil
}

func (p *DiscordProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(7 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *DiscordProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "discord_message_id"}, nil
}

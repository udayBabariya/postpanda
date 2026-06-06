package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type SlackProvider struct{}

func init() {
	integrations.Register("slack", &SlackProvider{})
}

func (p *SlackProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "slack", Name: "Slack", Type: "social", AuthType: "oauth2"}
}

func (p *SlackProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("SLACK_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/slack/callback"
	return fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=chat:write,channels:read&redirect_uri=%s&state=%s",
		clientID, redirectURI, orgID,
	)
}

func (p *SlackProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "slack_channel_id", Name: "Slack Channel", Type: "social", Token: code}, nil
}

func (p *SlackProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *SlackProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "slack_message_ts"}, nil
}

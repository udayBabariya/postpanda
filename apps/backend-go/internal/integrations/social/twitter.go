package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type TwitterProvider struct{}

func init() {
	integrations.Register("x", &TwitterProvider{})
	integrations.Register("twitter", &TwitterProvider{})
}

func (p *TwitterProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{
		ID:          "x",
		Name:        "X (Twitter)",
		Type:        "social",
		Description: "Post to X (formerly Twitter)",
		AuthType:    "oauth2",
	}
}

func (p *TwitterProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("X_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/x/callback"
	return fmt.Sprintf(
		"https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=tweet.read+tweet.write+users.read+offline.access&state=%s&code_challenge=challenge&code_challenge_method=plain",
		clientID, redirectURI, orgID,
	)
}

func (p *TwitterProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Exchange code for token via Twitter API
	// TODO: implement full Twitter OAuth 2.0 flow
	return &integrations.IntegrationInfo{
		InternalID: "twitter_user_id",
		Name:       "Twitter User",
		Type:       "social",
		Token:      code,
	}, nil
}

func (p *TwitterProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	// Twitter token refresh
	expiry := time.Now().Add(2 * time.Hour)
	return token, &expiry, nil
}

func (p *TwitterProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Post tweet via Twitter API v2
	// TODO: implement
	return &integrations.PublishResult{
		PostID: "tweet_id",
		URL:    "https://twitter.com/user/status/tweet_id",
	}, nil
}

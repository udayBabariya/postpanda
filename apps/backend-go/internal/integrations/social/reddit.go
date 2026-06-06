package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type RedditProvider struct{}

func init() {
	integrations.Register("reddit", &RedditProvider{})
}

func (p *RedditProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "reddit", Name: "Reddit", Type: "social", AuthType: "oauth2"}
}

func (p *RedditProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/reddit/callback"
	return fmt.Sprintf(
		"https://www.reddit.com/api/v1/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=submit+identity&state=%s&duration=permanent",
		clientID, redirectURI, orgID,
	)
}

func (p *RedditProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "reddit_user_id", Name: "Reddit User", Type: "profile", Token: code}, nil
}

func (p *RedditProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(1 * time.Hour)
	return token, &expiry, nil
}

func (p *RedditProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "reddit_post_id"}, nil
}

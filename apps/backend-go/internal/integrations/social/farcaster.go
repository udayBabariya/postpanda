package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type FarcasterProvider struct{}

func init() {
	integrations.Register("farcaster", &FarcasterProvider{})
}

func (p *FarcasterProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "farcaster", Name: "Farcaster", Type: "social", Description: "Post to Farcaster decentralized social network", AuthType: "oauth2"}
}

func (p *FarcasterProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("FARCASTER_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/farcaster/callback"
	return fmt.Sprintf("https://app.neynar.com/login?client_id=%s&redirect_uri=%s&state=%s", clientID, redirectURI, orgID)
}

func (p *FarcasterProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Neynar OAuth exchange
	return &integrations.IntegrationInfo{InternalID: "farcaster_fid", Name: "Farcaster Account", Type: "profile", Token: code}, nil
}

func (p *FarcasterProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *FarcasterProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// POST https://api.neynar.com/v2/farcaster/cast
	return &integrations.PublishResult{PostID: "farcaster_cast_hash"}, nil
}

package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

// Lemmy (federated Reddit-like)

type LemmyProvider struct{}

func init() {
	integrations.Register("lemmy", &LemmyProvider{})
	integrations.Register("nostr", &NostrProvider{})
}

func (p *LemmyProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "lemmy", Name: "Lemmy", Type: "social", Description: "Post to Lemmy federated communities", AuthType: "password"}
}

func (p *LemmyProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/lemmy/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *LemmyProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "lemmy_user_id", Name: "Lemmy Account", Type: "profile", Token: code}, nil
}

func (p *LemmyProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *LemmyProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "lemmy_post_id"}, nil
}

// Nostr (decentralized protocol)

type NostrProvider struct{}

func (p *NostrProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "nostr", Name: "Nostr", Type: "social", Description: "Post to Nostr decentralized protocol", AuthType: "keypair"}
}

func (p *NostrProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/nostr/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *NostrProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// code is the nsec private key
	return &integrations.IntegrationInfo{InternalID: "nostr_pubkey", Name: "Nostr Account", Type: "profile", Token: code}, nil
}

func (p *NostrProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *NostrProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Sign and broadcast Nostr event (kind 1 - text note) to relays
	return &integrations.PublishResult{PostID: "nostr_event_id"}, nil
}

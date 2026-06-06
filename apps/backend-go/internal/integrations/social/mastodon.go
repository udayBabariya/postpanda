package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type MastodonProvider struct{}
type BlueskyProvider struct{}
type ThreadsProvider struct{}
type PinterestProvider struct{}
type TelegramProvider struct{}

func init() {
	integrations.Register("mastodon", &MastodonProvider{})
	integrations.Register("bluesky", &BlueskyProvider{})
	integrations.Register("threads", &ThreadsProvider{})
	integrations.Register("pinterest", &PinterestProvider{})
	integrations.Register("telegram", &TelegramProvider{})
}

// Mastodon

func (p *MastodonProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "mastodon", Name: "Mastodon", Type: "social", AuthType: "oauth2"}
}

func (p *MastodonProvider) GetOAuthURL(orgID string) string {
	instance := os.Getenv("MASTODON_INSTANCE")
	if instance == "" {
		instance = "mastodon.social"
	}
	return fmt.Sprintf("https://%s/oauth/authorize?response_type=code&scope=write+read&state=%s", instance, orgID)
}

func (p *MastodonProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "mastodon_account_id", Name: "Mastodon Account", Type: "profile", Token: code}, nil
}

func (p *MastodonProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *MastodonProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "mastodon_status_id"}, nil
}

// Bluesky

func (p *BlueskyProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "bluesky", Name: "Bluesky", Type: "social", AuthType: "password"}
}

func (p *BlueskyProvider) GetOAuthURL(orgID string) string {
	return ""
}

func (p *BlueskyProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "bluesky_did", Name: "Bluesky Account", Type: "profile", Token: code}, nil
}

func (p *BlueskyProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *BlueskyProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "bluesky_post_uri"}, nil
}

// Threads

func (p *ThreadsProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "threads", Name: "Threads", Type: "social", AuthType: "oauth2"}
}

func (p *ThreadsProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("THREADS_CLIENT_ID")
	return fmt.Sprintf("https://threads.net/oauth/authorize?response_type=code&client_id=%s&scope=threads_basic,threads_content_publish&state=%s", clientID, orgID)
}

func (p *ThreadsProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "threads_user_id", Name: "Threads Account", Type: "profile", Token: code}, nil
}

func (p *ThreadsProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *ThreadsProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "threads_post_id"}, nil
}

// Pinterest

func (p *PinterestProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "pinterest", Name: "Pinterest", Type: "social", AuthType: "oauth2"}
}

func (p *PinterestProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("PINTEREST_CLIENT_ID")
	return fmt.Sprintf("https://www.pinterest.com/oauth/?response_type=code&client_id=%s&scope=boards:read,pins:write&state=%s", clientID, orgID)
}

func (p *PinterestProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "pinterest_user_id", Name: "Pinterest Account", Type: "profile", Token: code}, nil
}

func (p *PinterestProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(30 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *PinterestProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "pinterest_pin_id"}, nil
}

// Telegram

func (p *TelegramProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "telegram", Name: "Telegram", Type: "social", AuthType: "bot_token"}
}

func (p *TelegramProvider) GetOAuthURL(orgID string) string {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	return fmt.Sprintf("https://t.me/%s?start=%s", botToken, orgID)
}

func (p *TelegramProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "telegram_chat_id", Name: "Telegram Channel", Type: "channel", Token: code}, nil
}

func (p *TelegramProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *TelegramProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "telegram_message_id"}, nil
}

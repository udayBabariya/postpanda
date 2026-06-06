package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

// Skool, Whop, MeWe community platforms

type SkoolProvider struct{}
type WhopProvider struct{}
type MeWeProvider struct{}

func init() {
	integrations.Register("skool", &SkoolProvider{})
	integrations.Register("whop", &WhopProvider{})
	integrations.Register("mewe", &MeWeProvider{})
}

// Skool

func (p *SkoolProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "skool", Name: "Skool", Type: "community", Description: "Post to Skool community", AuthType: "api_key"}
}

func (p *SkoolProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/skool/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *SkoolProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "skool_community_id", Name: "Skool Community", Type: "community", Token: code}, nil
}

func (p *SkoolProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *SkoolProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "skool_post_id"}, nil
}

// Whop

func (p *WhopProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "whop", Name: "Whop", Type: "community", Description: "Post to Whop community", AuthType: "oauth2"}
}

func (p *WhopProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("WHOP_CLIENT_ID")
	return fmt.Sprintf("https://whop.com/oauth?client_id=%s&redirect_uri=%s&state=%s", clientID, os.Getenv("BACKEND_URL")+"/integrations/whop/callback", orgID)
}

func (p *WhopProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "whop_company_id", Name: "Whop Community", Type: "community", Token: code}, nil
}

func (p *WhopProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(30 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *WhopProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "whop_post_id"}, nil
}

// MeWe

func (p *MeWeProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "mewe", Name: "MeWe", Type: "social", Description: "Post to MeWe social network", AuthType: "oauth2"}
}

func (p *MeWeProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("MEWE_CLIENT_ID")
	return fmt.Sprintf("https://mewe.com/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&response_type=code", clientID, os.Getenv("BACKEND_URL")+"/integrations/mewe/callback", orgID)
}

func (p *MeWeProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "mewe_user_id", Name: "MeWe Account", Type: "profile", Token: code}, nil
}

func (p *MeWeProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(30 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *MeWeProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "mewe_post_id"}, nil
}

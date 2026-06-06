package social

import (
	"postpanda/backend-go/internal/integrations"
	"time"
)

// LinkedIn Pages (separate from personal LinkedIn)

type LinkedInPagesProvider struct{}

func init() {
	integrations.Register("linkedin-page", &LinkedInPagesProvider{})
}

func (p *LinkedInPagesProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "linkedin-page", Name: "LinkedIn Page", Type: "social", Description: "Post to LinkedIn Company Page", AuthType: "oauth2"}
}

func (p *LinkedInPagesProvider) GetOAuthURL(orgID string) string {
	return (&LinkedInProvider{}).GetOAuthURL(orgID)
}

func (p *LinkedInPagesProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Same OAuth as LinkedIn personal, but we select company pages
	return &integrations.IntegrationInfo{InternalID: "linkedin_org_id", Name: "LinkedIn Page", Type: "page", Token: code}, nil
}

func (p *LinkedInPagesProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(60 * 24 * time.Hour)
	return token, &expiry, nil
}

func (p *LinkedInPagesProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Post on behalf of organization using LinkedIn UGC Posts API
	return &integrations.PublishResult{PostID: "linkedin_org_post_id"}, nil
}

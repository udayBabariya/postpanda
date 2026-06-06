package integrations

import "time"

// Provider is the interface all social media integrations must implement.
type Provider interface {
	GetOAuthURL(orgID string) string
	HandleCallback(code, state string) (*IntegrationInfo, error)
	RefreshToken(token string, refreshToken *string) (string, *time.Time, error)
	PublishPost(integration *PostPublishRequest) (*PublishResult, error)
	GetProviderMeta() ProviderMeta
}

type IntegrationInfo struct {
	InternalID      string
	Name            string
	Picture         *string
	Type            string
	Token           string
	RefreshToken    *string
	TokenExpiration *time.Time
	Profile         *string
}

type PostPublishRequest struct {
	Token    string
	Content  string
	Image    []string
	Settings map[string]interface{}
}

type PublishResult struct {
	PostID string
	URL    string
}

type ProviderMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	AuthType    string `json:"authType"`
	IconURL     string `json:"iconUrl"`
}

var registry = map[string]Provider{}

func Register(name string, p Provider) {
	registry[name] = p
}

func GetProvider(name string) Provider {
	return registry[name]
}

func GetAllProvidersMeta() []map[string]interface{} {
	var result []map[string]interface{}
	for _, p := range registry {
		meta := p.GetProviderMeta()
		result = append(result, map[string]interface{}{
			"id":          meta.ID,
			"name":        meta.Name,
			"type":        meta.Type,
			"description": meta.Description,
			"authType":    meta.AuthType,
			"iconUrl":     meta.IconURL,
		})
	}
	return result
}

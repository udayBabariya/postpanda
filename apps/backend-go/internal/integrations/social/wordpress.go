package social

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type WordPressProvider struct{ client *http.Client }

func init() {
	integrations.Register("wordpress", &WordPressProvider{client: &http.Client{}})
}

func (p *WordPressProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "wordpress", Name: "WordPress", Type: "blog", Description: "Publish to WordPress site", AuthType: "oauth2"}
}

func (p *WordPressProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("WORDPRESS_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/wordpress/callback"
	return fmt.Sprintf("https://public-api.wordpress.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=posts&state=%s", clientID, redirectURI, orgID)
}

func (p *WordPressProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Exchange for access token
	clientID := os.Getenv("WORDPRESS_CLIENT_ID")
	clientSecret := os.Getenv("WORDPRESS_CLIENT_SECRET")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/wordpress/callback"

	body, _ := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
		"code":          code,
		"grant_type":    "authorization_code",
	})

	resp, err := p.client.Post("https://public-api.wordpress.com/oauth2/token", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		BlogID      string `json:"blog_id"`
		BlogURL     string `json:"blog_url"`
	}
	json.Unmarshal(b, &tokenResp)

	return &integrations.IntegrationInfo{
		InternalID: tokenResp.BlogID,
		Name:       tokenResp.BlogURL,
		Type:       "blog",
		Token:      tokenResp.AccessToken,
	}, nil
}

func (p *WordPressProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *WordPressProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	blogID := req.Settings["blogId"]
	body, _ := json.Marshal(map[string]interface{}{
		"title":   req.Content[:min2(len(req.Content), 100)],
		"content": req.Content,
		"status":  "publish",
	})
	httpReq, _ := http.NewRequest("POST",
		fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%v/posts/new", blogID),
		bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var result struct {
		ID  int    `json:"ID"`
		URL string `json:"URL"`
	}
	json.Unmarshal(b, &result)
	return &integrations.PublishResult{PostID: fmt.Sprintf("%d", result.ID), URL: result.URL}, nil
}

func min2(a, b int) int {
	if a < b { return a }
	return b
}

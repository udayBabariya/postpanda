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

// Medium

type MediumProvider struct{ client *http.Client }

func init() {
	integrations.Register("medium", &MediumProvider{client: &http.Client{}})
	integrations.Register("devto", &DevToProvider{client: &http.Client{}})
	integrations.Register("hashnode", &HashnodeProvider{client: &http.Client{}})
}

func (p *MediumProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "medium", Name: "Medium", Type: "blog", Description: "Publish to Medium", AuthType: "oauth2"}
}

func (p *MediumProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("MEDIUM_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/medium/callback"
	return fmt.Sprintf("https://medium.com/m/oauth/authorize?response_type=code&client_id=%s&scope=basicProfile,publishPost&redirect_uri=%s&state=%s", clientID, redirectURI, orgID)
}

func (p *MediumProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "medium_user_id", Name: "Medium Account", Type: "blog", Token: code}, nil
}

func (p *MediumProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *MediumProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Get user id
	meResp, err := p.makeRequest("GET", "https://api.medium.com/v1/me", req.Token, nil)
	if err != nil {
		return nil, err
	}
	var me struct{ Data struct{ ID string } }
	json.Unmarshal(meResp, &me)

	body, _ := json.Marshal(map[string]interface{}{
		"title":         req.Content[:min(len(req.Content), 100)],
		"contentFormat": "markdown",
		"content":       req.Content,
		"publishStatus": "public",
	})
	resp, err := p.makeRequest("POST", fmt.Sprintf("https://api.medium.com/v1/users/%s/posts", me.Data.ID), req.Token, body)
	if err != nil {
		return nil, err
	}
	var result struct{ Data struct{ ID string; URL string } }
	json.Unmarshal(resp, &result)
	return &integrations.PublishResult{PostID: result.Data.ID, URL: result.Data.URL}, nil
}

func (p *MediumProvider) makeRequest(method, url, token string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Dev.to

type DevToProvider struct{ client *http.Client }

func (p *DevToProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "devto", Name: "Dev.to", Type: "blog", Description: "Publish to Dev.to", AuthType: "api_key"}
}

func (p *DevToProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/devto/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *DevToProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// code is the API key for dev.to
	req, _ := http.NewRequest("GET", "https://dev.to/api/users/me", nil)
	req.Header.Set("api-key", code)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var user struct{ ID int; Username string; Name string }
	json.Unmarshal(b, &user)
	return &integrations.IntegrationInfo{
		InternalID: fmt.Sprintf("%d", user.ID),
		Name:       user.Name,
		Type:       "blog",
		Token:      code,
	}, nil
}

func (p *DevToProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *DevToProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"article": map[string]interface{}{
			"title":     req.Content[:min(len(req.Content), 100)],
			"body_markdown": req.Content,
			"published": true,
		},
	})
	httpReq, _ := http.NewRequest("POST", "https://dev.to/api/articles", bytes.NewReader(body))
	httpReq.Header.Set("api-key", req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var result struct{ ID int; URL string }
	json.Unmarshal(b, &result)
	return &integrations.PublishResult{PostID: fmt.Sprintf("%d", result.ID), URL: result.URL}, nil
}

// Hashnode

type HashnodeProvider struct{ client *http.Client }

func (p *HashnodeProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "hashnode", Name: "Hashnode", Type: "blog", Description: "Publish to Hashnode", AuthType: "api_key"}
}

func (p *HashnodeProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/hashnode/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *HashnodeProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Verify token via GraphQL
	return &integrations.IntegrationInfo{InternalID: "hashnode_user_id", Name: "Hashnode Account", Type: "blog", Token: code}, nil
}

func (p *HashnodeProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *HashnodeProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	query := `mutation PublishPost($input: PublishPostInput!) { publishPost(input: $input) { post { id url } } }`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":            req.Content[:min(len(req.Content), 100)],
			"contentMarkdown":  req.Content,
			"publicationId":    req.Settings["publicationId"],
		},
	}
	body, _ := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	httpReq, _ := http.NewRequest("POST", "https://gql.hashnode.com/", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &integrations.PublishResult{PostID: "hashnode_post_id"}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

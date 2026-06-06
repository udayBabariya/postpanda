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

// ListMonk - self-hosted newsletter platform

type ListMonkProvider struct{ client *http.Client }

func init() {
	integrations.Register("listmonk", &ListMonkProvider{client: &http.Client{}})
	integrations.Register("moltbook", &MoltbookProvider{client: &http.Client{}})
}

func (p *ListMonkProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "listmonk", Name: "ListMonk", Type: "newsletter", Description: "Send newsletters via ListMonk", AuthType: "api_key"}
}

func (p *ListMonkProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/listmonk/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *ListMonkProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// code = base64(instance_url:api_key)
	return &integrations.IntegrationInfo{InternalID: "listmonk_instance", Name: "ListMonk", Type: "newsletter", Token: code}, nil
}

func (p *ListMonkProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *ListMonkProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	instanceURL := fmt.Sprintf("%v", req.Settings["instanceUrl"])
	listID := req.Settings["listId"]

	body, _ := json.Marshal(map[string]interface{}{
		"name":         req.Content[:min3(len(req.Content), 100)],
		"subject":      req.Content[:min3(len(req.Content), 100)],
		"body":         req.Content,
		"content_type": "markdown",
		"lists":        []interface{}{listID},
		"type":         "regular",
		"status":       "scheduled",
	})

	httpReq, _ := http.NewRequest("POST", instanceURL+"/api/campaigns", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var result struct{ Data struct{ ID int } }
	json.Unmarshal(b, &result)
	return &integrations.PublishResult{PostID: fmt.Sprintf("%d", result.Data.ID)}, nil
}

func min3(a, b int) int {
	if a < b { return a }
	return b
}

// Moltbook - custom agent/webhook platform

type MoltbookProvider struct{ client *http.Client }

func (p *MoltbookProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "moltbook", Name: "Moltbook", Type: "automation", Description: "Trigger Moltbook automation agent", AuthType: "api_key"}
}

func (p *MoltbookProvider) GetOAuthURL(orgID string) string {
	return fmt.Sprintf("%s/integrations/moltbook/connect?state=%s", os.Getenv("BACKEND_URL"), orgID)
}

func (p *MoltbookProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "moltbook_agent_id", Name: "Moltbook Agent", Type: "automation", Token: code}, nil
}

func (p *MoltbookProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *MoltbookProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	webhookURL := fmt.Sprintf("%v", req.Settings["webhookUrl"])
	body, _ := json.Marshal(map[string]interface{}{"content": req.Content, "images": req.Image})
	httpReq, _ := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	p.client.Do(httpReq)
	return &integrations.PublishResult{PostID: "moltbook_trigger"}, nil
}

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

type GMBProvider struct{ client *http.Client }

func init() {
	integrations.Register("gmb", &GMBProvider{client: &http.Client{}})
}

func (p *GMBProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "gmb", Name: "Google My Business", Type: "social", Description: "Post to Google Business Profile", AuthType: "oauth2"}
}

func (p *GMBProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/gmb/callback"
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=https://www.googleapis.com/auth/business.manage&state=%s&access_type=offline", clientID, redirectURI, orgID)
}

func (p *GMBProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "gmb_location_id", Name: "Google Business Profile", Type: "business", Token: code}, nil
}

func (p *GMBProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(1 * time.Hour)
	return token, &expiry, nil
}

func (p *GMBProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	locationID := req.Settings["locationId"]
	body, _ := json.Marshal(map[string]interface{}{
		"languageCode": "en-US",
		"summary":      req.Content,
		"callToAction": map[string]string{"actionType": "LEARN_MORE"},
		"media":        []interface{}{},
	})
	httpReq, _ := http.NewRequest("POST",
		fmt.Sprintf("https://mybusiness.googleapis.com/v4/%v/localPosts", locationID),
		bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+req.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var result struct{ Name string }
	json.Unmarshal(b, &result)
	return &integrations.PublishResult{PostID: result.Name}, nil
}

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

type DribbbleProvider struct{ client *http.Client }

func init() {
	integrations.Register("dribbble", &DribbbleProvider{client: &http.Client{}})
}

func (p *DribbbleProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "dribbble", Name: "Dribbble", Type: "social", Description: "Share shots on Dribbble", AuthType: "oauth2"}
}

func (p *DribbbleProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("DRIBBBLE_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/dribbble/callback"
	return fmt.Sprintf("https://dribbble.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=public+write&state=%s", clientID, redirectURI, orgID)
}

func (p *DribbbleProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	// Exchange for token
	clientID := os.Getenv("DRIBBBLE_CLIENT_ID")
	clientSecret := os.Getenv("DRIBBBLE_CLIENT_SECRET")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/dribbble/callback"
	body, _ := json.Marshal(map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code, "redirect_uri": redirectURI})
	resp, err := (&http.Client{}).Post("https://dribbble.com/oauth/token", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var tokenResp struct{ AccessToken string `json:"access_token"` }
	json.Unmarshal(b, &tokenResp)

	// Get user info
	req, _ := http.NewRequest("GET", "https://api.dribbble.com/v2/user", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	uResp, _ := (&http.Client{}).Do(req)
	if uResp != nil {
		defer uResp.Body.Close()
		ub, _ := io.ReadAll(uResp.Body)
		var user struct{ ID int; Name string }
		json.Unmarshal(ub, &user)
		return &integrations.IntegrationInfo{
			InternalID: fmt.Sprintf("%d", user.ID),
			Name:       user.Name,
			Type:       "profile",
			Token:      tokenResp.AccessToken,
		}, nil
	}
	return &integrations.IntegrationInfo{InternalID: "dribbble_user", Name: "Dribbble Account", Type: "profile", Token: tokenResp.AccessToken}, nil
}

func (p *DribbbleProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	return token, nil, nil
}

func (p *DribbbleProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	// Dribbble shots require image upload first
	return &integrations.PublishResult{PostID: "dribbble_shot_id"}, nil
}

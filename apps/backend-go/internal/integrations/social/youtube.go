package social

import (
	"fmt"
	"os"
	"postpanda/backend-go/internal/integrations"
	"time"
)

type YouTubeProvider struct{}

func init() {
	integrations.Register("youtube", &YouTubeProvider{})
}

func (p *YouTubeProvider) GetProviderMeta() integrations.ProviderMeta {
	return integrations.ProviderMeta{ID: "youtube", Name: "YouTube", Type: "social", AuthType: "oauth2"}
}

func (p *YouTubeProvider) GetOAuthURL(orgID string) string {
	clientID := os.Getenv("YOUTUBE_CLIENT_ID")
	redirectURI := os.Getenv("BACKEND_URL") + "/integrations/youtube/callback"
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=https://www.googleapis.com/auth/youtube.upload+https://www.googleapis.com/auth/youtube.readonly&state=%s&access_type=offline",
		clientID, redirectURI, orgID,
	)
}

func (p *YouTubeProvider) HandleCallback(code, state string) (*integrations.IntegrationInfo, error) {
	return &integrations.IntegrationInfo{InternalID: "yt_channel_id", Name: "YouTube Channel", Type: "channel", Token: code}, nil
}

func (p *YouTubeProvider) RefreshToken(token string, refreshToken *string) (string, *time.Time, error) {
	expiry := time.Now().Add(1 * time.Hour)
	return token, &expiry, nil
}

func (p *YouTubeProvider) PublishPost(req *integrations.PostPublishRequest) (*integrations.PublishResult, error) {
	return &integrations.PublishResult{PostID: "yt_video_id"}, nil
}

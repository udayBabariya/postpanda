package services

import (
	"context"
	"encoding/json"
	"fmt"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/google/uuid"
)

// SettingsService

type SettingsService struct{}

func NewSettingsService() *SettingsService { return &SettingsService{} }

func (s *SettingsService) GetOrgSettings(orgID string) (map[string]interface{}, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, description, api_key, shortlink, allow_trial, is_trailing FROM organizations WHERE id=$1`, orgID)
	var id, name, shortlink string
	var description, apiKey *string
	var allowTrial, isTrialing bool
	if err := row.Scan(&id, &name, &description, &apiKey, &shortlink, &allowTrial, &isTrialing); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": description,
		"apiKey":      apiKey,
		"shortlink":   shortlink,
		"allowTrial":  allowTrial,
		"isTrialing":  isTrialing,
	}, nil
}

func (s *SettingsService) UpdateOrgSettings(orgID string, updates map[string]interface{}) error {
	ctx := context.Background()
	if name, ok := updates["name"]; ok {
		database.DB.Exec(ctx, `UPDATE organizations SET name=$1, updated_at=$2 WHERE id=$3`, name, time.Now(), orgID)
	}
	if desc, ok := updates["description"]; ok {
		database.DB.Exec(ctx, `UPDATE organizations SET description=$1, updated_at=$2 WHERE id=$3`, desc, time.Now(), orgID)
	}
	return nil
}

func (s *SettingsService) RegenerateAPIKey(orgID string) (string, error) {
	ctx := context.Background()
	newKey := uuid.New().String()
	_, err := database.DB.Exec(ctx, `UPDATE organizations SET api_key=$1, updated_at=$2 WHERE id=$3`, newKey, time.Now(), orgID)
	return newKey, err
}

func (s *SettingsService) UpdateShortLink(orgID, preference string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE organizations SET shortlink=$1, updated_at=$2 WHERE id=$3`, preference, time.Now(), orgID)
	return err
}

// SetService

type SetService struct{}

func NewSetService() *SetService { return &SetService{} }

func (s *SetService) List(orgID string) ([]*models.Set, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, name, content, created_at, updated_at FROM sets WHERE organization_id=$1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []*models.Set
	for rows.Next() {
		s := &models.Set{}
		rows.Scan(&s.ID, &s.OrganizationID, &s.Name, &s.Content, &s.CreatedAt, &s.UpdatedAt)
		sets = append(sets, s)
	}
	return sets, nil
}

func (s *SetService) Create(orgID, name, content string) (*models.Set, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO sets (id, organization_id, name, content, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, orgID, name, content, now, now)
	if err != nil {
		return nil, err
	}
	return &models.Set{ID: id, OrganizationID: orgID, Name: name, Content: content, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *SetService) Update(orgID, id, name, content string) (*models.Set, error) {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE sets SET name=$1, content=$2, updated_at=$3 WHERE id=$4 AND organization_id=$5`,
		name, content, now, id, orgID)
	if err != nil {
		return nil, err
	}
	return &models.Set{ID: id, OrganizationID: orgID, Name: name, Content: content, UpdatedAt: now}, nil
}

func (s *SetService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `DELETE FROM sets WHERE id=$1 AND organization_id=$2`, id, orgID)
	return err
}

// SignatureService

type SignatureService struct{}

func NewSignatureService() *SignatureService { return &SignatureService{} }

func (s *SignatureService) List(orgID string) ([]*models.Signature, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, content, auto_add, deleted_at, created_at, updated_at
		FROM signatures WHERE organization_id=$1 AND deleted_at IS NULL ORDER BY created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sigs []*models.Signature
	for rows.Next() {
		sig := &models.Signature{}
		rows.Scan(&sig.ID, &sig.OrganizationID, &sig.Content, &sig.AutoAdd, &sig.DeletedAt, &sig.CreatedAt, &sig.UpdatedAt)
		sigs = append(sigs, sig)
	}
	return sigs, nil
}

func (s *SignatureService) Create(orgID, content string, autoAdd bool) (*models.Signature, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO signatures (id, organization_id, content, auto_add, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, orgID, content, autoAdd, now, now)
	if err != nil {
		return nil, err
	}
	return &models.Signature{ID: id, OrganizationID: orgID, Content: content, AutoAdd: autoAdd, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *SignatureService) Update(orgID, id, content string, autoAdd bool) (*models.Signature, error) {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE signatures SET content=$1, auto_add=$2, updated_at=$3 WHERE id=$4 AND organization_id=$5`,
		content, autoAdd, now, id, orgID)
	if err != nil {
		return nil, err
	}
	return &models.Signature{ID: id, OrganizationID: orgID, Content: content, AutoAdd: autoAdd, UpdatedAt: now}, nil
}

func (s *SignatureService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE signatures SET deleted_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	return err
}

// AutoPostService

type AutoPostService struct{}

func NewAutoPostService() *AutoPostService { return &AutoPostService{} }

func (s *AutoPostService) List(orgID string) ([]*models.AutoPost, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, title, content, on_slot, sync_last, url, last_url, active,
		add_picture, generate_content, integrations, deleted_at, created_at, updated_at
		FROM auto_posts WHERE organization_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var autoPosts []*models.AutoPost
	for rows.Next() {
		ap := &models.AutoPost{}
		rows.Scan(&ap.ID, &ap.OrganizationID, &ap.Title, &ap.Content, &ap.OnSlot, &ap.SyncLast,
			&ap.URL, &ap.LastURL, &ap.Active, &ap.AddPicture, &ap.GenerateContent, &ap.Integrations,
			&ap.DeletedAt, &ap.CreatedAt, &ap.UpdatedAt)
		autoPosts = append(autoPosts, ap)
	}
	return autoPosts, nil
}

func (s *AutoPostService) Create(orgID string, req interface{}) (*models.AutoPost, error) {
	ctx := context.Background()
	b, _ := json.Marshal(req)
	var r struct {
		Title           string   `json:"title"`
		Content         string   `json:"content"`
		OnSlot          bool     `json:"onSlot"`
		SyncLast        bool     `json:"syncLast"`
		URL             string   `json:"url"`
		Active          bool     `json:"active"`
		AddPicture      bool     `json:"addPicture"`
		GenerateContent bool     `json:"generateContent"`
		Integrations    []string `json:"integrations"`
	}
	json.Unmarshal(b, &r)

	id := uuid.New().String()
	now := time.Now()
	intJSON, _ := json.Marshal(r.Integrations)

	_, err := database.DB.Exec(ctx,
		`INSERT INTO auto_posts (id, organization_id, title, content, on_slot, sync_last, url, last_url,
		active, add_picture, generate_content, integrations, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		id, orgID, r.Title, r.Content, r.OnSlot, r.SyncLast, r.URL, r.URL,
		r.Active, r.AddPicture, r.GenerateContent, string(intJSON), now, now)
	if err != nil {
		return nil, err
	}
	return &models.AutoPost{ID: id, OrganizationID: orgID, Title: r.Title, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *AutoPostService) Update(orgID, id string, req map[string]interface{}) (*models.AutoPost, error) {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE auto_posts SET updated_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	if err != nil {
		return nil, err
	}
	return &models.AutoPost{ID: id, OrganizationID: orgID}, nil
}

func (s *AutoPostService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE auto_posts SET deleted_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	return err
}

func (s *AutoPostService) Toggle(orgID, id string, active bool) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE auto_posts SET active=$1, updated_at=$2 WHERE id=$3 AND organization_id=$4`, active, time.Now(), id, orgID)
	return err
}

// AnnouncementService

type AnnouncementService struct{}

func NewAnnouncementService() *AnnouncementService { return &AnnouncementService{} }

func (s *AnnouncementService) List() ([]*models.Announcement, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, title, description, color, created_at FROM announcements ORDER BY created_at DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var announcements []*models.Announcement
	for rows.Next() {
		a := &models.Announcement{}
		rows.Scan(&a.ID, &a.Title, &a.Description, &a.Color, &a.CreatedAt)
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (s *AnnouncementService) Create(title, description, color string) (*models.Announcement, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO announcements (id, title, description, color, created_at) VALUES ($1,$2,$3,$4,$5)`,
		id, title, description, color, now)
	if err != nil {
		return nil, err
	}
	return &models.Announcement{ID: id, Title: title, Description: description, CreatedAt: now}, nil
}

func (s *AnnouncementService) Delete(id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `DELETE FROM announcements WHERE id=$1`, id)
	return err
}

// OAuthAppService

type OAuthAppService struct{}

func NewOAuthAppService() *OAuthAppService { return &OAuthAppService{} }

func (s *OAuthAppService) GetByOrg(orgID string) (*models.OAuthApp, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, organization_id, name, description, picture_id, redirect_url, client_id, deleted_at, created_at, updated_at
		FROM oauth_apps WHERE organization_id=$1 AND deleted_at IS NULL LIMIT 1`, orgID)
	app := &models.OAuthApp{}
	err := row.Scan(&app.ID, &app.OrganizationID, &app.Name, &app.Description, &app.PictureID,
		&app.RedirectURL, &app.ClientID, &app.DeletedAt, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *OAuthAppService) Create(orgID, name, description, redirectURL string) (*models.OAuthApp, error) {
	ctx := context.Background()
	id := uuid.New().String()
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO oauth_apps (id, organization_id, name, description, redirect_url, client_id, client_secret, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		id, orgID, name, description, redirectURL, clientID, clientSecret, now, now)
	if err != nil {
		return nil, err
	}
	return &models.OAuthApp{ID: id, OrganizationID: orgID, Name: name, ClientID: clientID, RedirectURL: redirectURL, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *OAuthAppService) Update(orgID, id, name, description, redirectURL string) (*models.OAuthApp, error) {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE oauth_apps SET name=$1, description=$2, redirect_url=$3, updated_at=$4 WHERE id=$5 AND organization_id=$6`,
		name, description, redirectURL, now, id, orgID)
	if err != nil {
		return nil, err
	}
	return s.GetByOrg(orgID)
}

func (s *OAuthAppService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE oauth_apps SET deleted_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	return err
}

func (s *OAuthAppService) GenerateAuthCode(clientID, redirectURI, state string) (string, error) {
	ctx := context.Background()
	var appID string
	err := database.DB.QueryRow(ctx, `SELECT id FROM oauth_apps WHERE client_id=$1 AND deleted_at IS NULL`, clientID).Scan(&appID)
	if err != nil {
		return "", err
	}
	code := uuid.New().String()
	database.Redis.Set(ctx, "oauth_code:"+code, appID, 10*time.Minute)
	return code, nil
}

func (s *OAuthAppService) ExchangeCode(code, clientID, clientSecret string) (string, error) {
	ctx := context.Background()
	appID, err := database.Redis.Get(ctx, "oauth_code:"+code).Result()
	if err != nil {
		return "", err
	}
	database.Redis.Del(ctx, "oauth_code:"+code)

	var storedSecret string
	database.DB.QueryRow(ctx, `SELECT client_secret FROM oauth_apps WHERE id=$1 AND client_id=$2`, appID, clientID).Scan(&storedSecret)
	if storedSecret != clientSecret {
		return "", errInvalidCredentials
	}

	token := uuid.New().String()
	database.Redis.Set(ctx, "oauth_token:"+token, appID, 24*time.Hour)
	return token, nil
}

var errInvalidCredentials = fmt.Errorf("invalid credentials")

package services

import (
	"context"
	"encoding/json"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/google/uuid"
)

type WebhookService struct{}

func NewWebhookService() *WebhookService { return &WebhookService{} }

func (s *WebhookService) List(orgID string) ([]*models.Webhook, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, organization_id, url, deleted_at, created_at, updated_at
		FROM webhooks WHERE organization_id=$1 AND deleted_at IS NULL ORDER BY created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*models.Webhook
	for rows.Next() {
		w := &models.Webhook{}
		rows.Scan(&w.ID, &w.Name, &w.OrganizationID, &w.URL, &w.DeletedAt, &w.CreatedAt, &w.UpdatedAt)
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}

func (s *WebhookService) Create(orgID, name, url string, integrations []string) (*models.Webhook, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO webhooks (id, name, organization_id, url, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, name, orgID, url, now, now)
	if err != nil {
		return nil, err
	}

	for _, intID := range integrations {
		database.DB.Exec(ctx,
			`INSERT INTO integrations_webhooks (integration_id, webhook_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			intID, id)
	}

	return &models.Webhook{ID: id, Name: name, OrganizationID: orgID, URL: url, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *WebhookService) Update(orgID, id, name, url string, integrations []string) (*models.Webhook, error) {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE webhooks SET name=$1, url=$2, updated_at=$3 WHERE id=$4 AND organization_id=$5`,
		name, url, now, id, orgID)
	if err != nil {
		return nil, err
	}

	database.DB.Exec(ctx, `DELETE FROM integrations_webhooks WHERE webhook_id=$1`, id)
	for _, intID := range integrations {
		database.DB.Exec(ctx,
			`INSERT INTO integrations_webhooks (integration_id, webhook_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			intID, id)
	}

	return &models.Webhook{ID: id, Name: name, OrganizationID: orgID, URL: url, UpdatedAt: now}, nil
}

func (s *WebhookService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE webhooks SET deleted_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	return err
}

func (s *WebhookService) Dispatch(orgID, event string, payload interface{}) {
	ctx := context.Background()
	rows, _ := database.DB.Query(ctx,
		`SELECT w.url FROM webhooks w WHERE w.organization_id=$1 AND w.deleted_at IS NULL`, orgID)
	if rows == nil {
		return
	}
	defer rows.Close()

	body, _ := json.Marshal(map[string]interface{}{"event": event, "data": payload})
	_ = body // TODO: dispatch HTTP requests to webhook URLs
}

package services

import (
	"context"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService { return &NotificationService{} }

func (s *NotificationService) List(orgID string) ([]*models.Notification, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, content, link, deleted_at, created_at, updated_at
		FROM notifications WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 50`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		rows.Scan(&n.ID, &n.OrganizationID, &n.Content, &n.Link, &n.DeletedAt, &n.CreatedAt, &n.UpdatedAt)
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (s *NotificationService) MarkRead(userID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET last_read_notifications=$1 WHERE id=$2`, time.Now(), userID)
	return err
}

func (s *NotificationService) GetUnreadCount(orgID, userID string) (int, error) {
	ctx := context.Background()
	var lastRead time.Time
	database.DB.QueryRow(ctx, `SELECT last_read_notifications FROM users WHERE id=$1`, userID).Scan(&lastRead)

	var count int
	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE organization_id=$1 AND deleted_at IS NULL AND created_at > $2`,
		orgID, lastRead).Scan(&count)
	return count, nil
}

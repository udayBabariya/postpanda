package services

import (
	"context"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/google/uuid"
)

type CommentService struct{}

func NewCommentService() *CommentService { return &CommentService{} }

func (s *CommentService) ListByPost(postID string) ([]*models.Comment, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, content, organization_id, post_id, user_id, deleted_at, created_at, updated_at
		FROM comments WHERE post_id=$1 AND deleted_at IS NULL ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.Comment
	for rows.Next() {
		cm := &models.Comment{}
		rows.Scan(&cm.ID, &cm.Content, &cm.OrganizationID, &cm.PostID, &cm.UserID, &cm.DeletedAt, &cm.CreatedAt, &cm.UpdatedAt)
		comments = append(comments, cm)
	}
	return comments, nil
}

func (s *CommentService) Create(orgID, userID, postID, content string) (*models.Comment, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO comments (id, content, organization_id, post_id, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		id, content, orgID, postID, userID, now, now)
	if err != nil {
		return nil, err
	}
	return &models.Comment{ID: id, Content: content, OrganizationID: orgID, PostID: postID, UserID: userID, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *CommentService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE comments SET deleted_at=$1 WHERE id=$2 AND organization_id=$3`, time.Now(), id, orgID)
	return err
}

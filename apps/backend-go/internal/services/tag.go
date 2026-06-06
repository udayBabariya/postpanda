package services

import (
	"context"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/google/uuid"
)

type TagService struct{}

func NewTagService() *TagService { return &TagService{} }

func (s *TagService) List(orgID string) ([]*models.Tag, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, color, org_id, deleted_at, created_at, updated_at
		FROM tags WHERE org_id = $1 AND deleted_at IS NULL ORDER BY name ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		t := &models.Tag{}
		rows.Scan(&t.ID, &t.Name, &t.Color, &t.OrgID, &t.DeletedAt, &t.CreatedAt, &t.UpdatedAt)
		tags = append(tags, t)
	}
	return tags, nil
}

func (s *TagService) Create(orgID, name, color string) (*models.Tag, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO tags (id, name, color, org_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, name, color, orgID, now, now)
	if err != nil {
		return nil, err
	}
	return &models.Tag{ID: id, Name: name, Color: color, OrgID: orgID, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *TagService) Update(orgID, id, name, color string) (*models.Tag, error) {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE tags SET name=$1, color=$2, updated_at=$3 WHERE id=$4 AND org_id=$5 AND deleted_at IS NULL`,
		name, color, now, id, orgID)
	if err != nil {
		return nil, err
	}
	return &models.Tag{ID: id, Name: name, Color: color, OrgID: orgID, UpdatedAt: now}, nil
}

func (s *TagService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE tags SET deleted_at=$1 WHERE id=$2 AND org_id=$3`, time.Now(), id, orgID)
	return err
}

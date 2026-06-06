package services

import (
	"context"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type MediaService struct{}

func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) List(orgID, page, limit, mediaType string) ([]*models.Media, int, error) {
	ctx := context.Background()

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = 20
	}
	offset := (p - 1) * l

	query := `SELECT id, name, original_name, path, organization_id, file_size, type, thumbnail, alt,
		thumbnail_timestamp, deleted_at, created_at, updated_at
		FROM media WHERE organization_id = $1 AND deleted_at IS NULL`
	args := []interface{}{orgID}
	argIdx := 2

	if mediaType != "" {
		query += " AND type = $" + strconv.Itoa(argIdx)
		args = append(args, mediaType)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, l, offset)

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var media []*models.Media
	for rows.Next() {
		m := &models.Media{}
		if err := rows.Scan(
			&m.ID, &m.Name, &m.OriginalName, &m.Path, &m.OrganizationID,
			&m.FileSize, &m.Type, &m.Thumbnail, &m.Alt, &m.ThumbnailTimestamp,
			&m.DeletedAt, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		media = append(media, m)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM media WHERE organization_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{orgID}
	if mediaType != "" {
		countQuery += " AND type = $2"
		countArgs = append(countArgs, mediaType)
	}
	database.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)

	return media, total, nil
}

func (s *MediaService) Create(orgID, originalName, name, path string, size int, mediaType string) (*models.Media, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()

	_, err := database.DB.Exec(ctx,
		`INSERT INTO media (id, name, original_name, path, organization_id, file_size, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		id, name, originalName, path, orgID, size, mediaType, now, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.Media{
		ID:             id,
		Name:           name,
		OriginalName:   &originalName,
		Path:           path,
		OrganizationID: orgID,
		FileSize:       size,
		Type:           mediaType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (s *MediaService) Delete(orgID, id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE media SET deleted_at = $1 WHERE id = $2 AND organization_id = $3`,
		now, id, orgID,
	)
	return err
}

func (s *MediaService) UpdateAlt(orgID, id, alt string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE media SET alt = $1, updated_at = $2 WHERE id = $3 AND organization_id = $4`,
		alt, time.Now(), id, orgID,
	)
	return err
}

func (s *MediaService) SetUserPicture(userID, mediaID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET picture_id = $1 WHERE id = $2`,
		mediaID, userID,
	)
	return err
}

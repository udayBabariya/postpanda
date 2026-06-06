package services

import (
	"context"
	"encoding/json"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

type PostCreateRequest struct {
	Date         string
	Integrations []string
	Posts        []interface{}
	Tags         []string
}

func (s *PostService) Create(orgID, userID string, req interface{}) ([]*models.Post, error) {
	ctx := context.Background()

	type createReq struct {
		Date         string   `json:"date"`
		Integrations []string `json:"integrations"`
		Posts        []struct {
			Content  string      `json:"content"`
			Image    interface{} `json:"image"`
			Settings interface{} `json:"settings"`
		} `json:"posts"`
		Tags     []string    `json:"tags"`
		Settings interface{} `json:"settings"`
	}

	b, _ := json.Marshal(req)
	var r createReq
	json.Unmarshal(b, &r)

	publishDate, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		return nil, err
	}

	groupID := uuid.New().String()
	var createdPosts []*models.Post

	for i, postContent := range r.Posts {
		if i >= len(r.Integrations) {
			break
		}

		integrationID := r.Integrations[i]
		postID := uuid.NewString()
		settingsJSON := "{}"
		if postContent.Settings != nil {
			b, _ := json.Marshal(postContent.Settings)
			settingsJSON = string(b)
		}

		now := time.Now()
		_, err := database.DB.Exec(ctx,
			`INSERT INTO posts (id, state, publish_date, organization_id, integration_id, content,
			delay, "group", approved_submit_for_order, creation_method, settings, created_at, updated_at)
			VALUES ($1, 'QUEUE', $2, $3, $4, $5, 0, $6, 'NO', 'WEB', $7, $8, $9)`,
			postID, publishDate, orgID, integrationID, postContent.Content, groupID, settingsJSON, now, now,
		)
		if err != nil {
			return nil, err
		}

		for _, tagID := range r.Tags {
			database.DB.Exec(ctx,
				`INSERT INTO tags_posts (post_id, tag_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
				postID, tagID, now, now,
			)
		}

		createdPosts = append(createdPosts, &models.Post{
			ID:             postID,
			State:          models.StateQueue,
			PublishDate:    publishDate,
			OrganizationID: orgID,
			IntegrationID:  integrationID,
			Content:        postContent.Content,
			Group:          groupID,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	return createdPosts, nil
}

func (s *PostService) List(orgID, startDate, endDate, integrationID string) ([]*models.Post, error) {
	ctx := context.Background()

	query := `SELECT id, state, publish_date, organization_id, integration_id, content, delay, "group",
		title, description, parent_post_id, settings, image, creation_method, error, deleted_at, created_at, updated_at
		FROM posts
		WHERE organization_id = $1 AND deleted_at IS NULL`
	args := []interface{}{orgID}
	argIdx := 2

	if startDate != "" {
		query += " AND publish_date >= $" + string(rune('0'+argIdx))
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += " AND publish_date <= $" + string(rune('0'+argIdx))
		args = append(args, endDate)
		argIdx++
	}
	if integrationID != "" {
		query += " AND integration_id = $" + string(rune('0'+argIdx))
		args = append(args, integrationID)
	}

	query += " ORDER BY publish_date ASC"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		if err := rows.Scan(
			&p.ID, &p.State, &p.PublishDate, &p.OrganizationID, &p.IntegrationID,
			&p.Content, &p.Delay, &p.Group, &p.Title, &p.Description, &p.ParentPostID,
			&p.Settings, &p.Image, &p.CreationMethod, &p.Error, &p.DeletedAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *PostService) GetByID(orgID, postID string) (*models.Post, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, state, publish_date, organization_id, integration_id, content, delay, "group",
		title, description, parent_post_id, settings, image, creation_method, error, deleted_at, created_at, updated_at
		FROM posts WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		postID, orgID)

	p := &models.Post{}
	err := row.Scan(
		&p.ID, &p.State, &p.PublishDate, &p.OrganizationID, &p.IntegrationID,
		&p.Content, &p.Delay, &p.Group, &p.Title, &p.Description, &p.ParentPostID,
		&p.Settings, &p.Image, &p.CreationMethod, &p.Error, &p.DeletedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostService) Update(orgID, postID string, req interface{}) (*models.Post, error) {
	ctx := context.Background()
	type updateReq struct {
		Date    string `json:"date"`
		Content string `json:"content"`
	}
	b, _ := json.Marshal(req)
	var r updateReq
	json.Unmarshal(b, &r)

	publishDate, _ := time.Parse(time.RFC3339, r.Date)
	_, err := database.DB.Exec(ctx,
		`UPDATE posts SET content = $1, publish_date = $2, updated_at = $3
		WHERE id = $4 AND organization_id = $5 AND deleted_at IS NULL`,
		r.Content, publishDate, time.Now(), postID, orgID,
	)
	if err != nil {
		return nil, err
	}
	return s.GetByID(orgID, postID)
}

func (s *PostService) Delete(orgID, postID string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`UPDATE posts SET deleted_at = $1 WHERE id = $2 AND organization_id = $3`,
		now, postID, orgID,
	)
	return err
}

func (s *PostService) GetGroup(orgID, groupID string) ([]*models.Post, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, state, publish_date, organization_id, integration_id, content, delay, "group",
		title, description, parent_post_id, settings, image, creation_method, error, deleted_at, created_at, updated_at
		FROM posts WHERE "group" = $1 AND organization_id = $2 AND deleted_at IS NULL
		ORDER BY created_at ASC`,
		groupID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		rows.Scan(
			&p.ID, &p.State, &p.PublishDate, &p.OrganizationID, &p.IntegrationID,
			&p.Content, &p.Delay, &p.Group, &p.Title, &p.Description, &p.ParentPostID,
			&p.Settings, &p.Image, &p.CreationMethod, &p.Error, &p.DeletedAt,
			&p.CreatedAt, &p.UpdatedAt,
		)
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *PostService) DeleteGroup(orgID, groupID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE posts SET deleted_at = $1 WHERE "group" = $2 AND organization_id = $3`,
		time.Now(), groupID, orgID,
	)
	return err
}

func (s *PostService) Reschedule(orgID, postID, date string) (*models.Post, error) {
	ctx := context.Background()
	publishDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}
	_, err = database.DB.Exec(ctx,
		`UPDATE posts SET publish_date = $1, updated_at = $2
		WHERE id = $3 AND organization_id = $4`,
		publishDate, time.Now(), postID, orgID,
	)
	if err != nil {
		return nil, err
	}
	return s.GetByID(orgID, postID)
}

func (s *PostService) SubmitForApproval(orgID, postID, orderID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE posts SET submitted_for_order_id = $1, approved_submit_for_order = 'WAITING_CONFIRMATION',
		updated_at = $2 WHERE id = $3 AND organization_id = $4`,
		orderID, time.Now(), postID, orgID,
	)
	return err
}

func (s *PostService) Approve(orgID, postID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE posts SET approved_submit_for_order = 'YES', updated_at = $1
		WHERE id = $2 AND organization_id = $3`,
		time.Now(), postID, orgID,
	)
	return err
}

func (s *PostService) GetPublicPost(postID string) (*models.Post, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, state, publish_date, organization_id, integration_id, content, delay, "group",
		title, description, parent_post_id, settings, image, creation_method, error, deleted_at, created_at, updated_at
		FROM posts WHERE id = $1 AND deleted_at IS NULL`,
		postID)

	p := &models.Post{}
	err := row.Scan(
		&p.ID, &p.State, &p.PublishDate, &p.OrganizationID, &p.IntegrationID,
		&p.Content, &p.Delay, &p.Group, &p.Title, &p.Description, &p.ParentPostID,
		&p.Settings, &p.Image, &p.CreationMethod, &p.Error, &p.DeletedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostService) GetAnalytics(orgID, postID string) (interface{}, error) {
	// Return analytics data - to be implemented with actual analytics queries
	return map[string]interface{}{
		"postId": postID,
		"views":  0,
		"likes":  0,
		"shares": 0,
	}, nil
}

func (s *PostService) GetStatistics(orgID, postID string) (map[string]interface{}, error) {
	ctx := context.Background()
	var publishDate time.Time
	var state, integrationID string
	database.DB.QueryRow(ctx,
		`SELECT publish_date, state, integration_id FROM posts WHERE id=$1 AND organization_id=$2`,
		postID, orgID).Scan(&publishDate, &state, &integrationID)
	return map[string]interface{}{
		"postId":        postID,
		"state":         state,
		"publishDate":   publishDate,
		"integrationId": integrationID,
		"views":         0,
		"likes":         0,
		"shares":        0,
	}, nil
}

func (s *PostService) FindSlot(orgID, integrationID string) ([]string, error) {
	ctx := context.Background()
	// Get posting times for integration
	var postingTimes string
	database.DB.QueryRow(ctx,
		`SELECT posting_times FROM integrations WHERE id=$1 AND organization_id=$2`,
		integrationID, orgID).Scan(&postingTimes)

	// Find next available slots
	now := time.Now()
	var slots []string
	for i := 1; i <= 5; i++ {
		slot := now.Add(time.Duration(i*24) * time.Hour).Truncate(time.Hour)
		slots = append(slots, slot.Format(time.RFC3339))
	}
	return slots, nil
}

func (s *PostService) GetOldPosts(orgID, page string) ([]*models.Post, int, error) {
	ctx := context.Background()
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	offset := (p - 1) * 20

	rows, err := database.DB.Query(ctx,
		`SELECT id, state, publish_date, organization_id, integration_id, content, delay, "group",
		title, description, parent_post_id, settings, image, creation_method, error, deleted_at, created_at, updated_at
		FROM posts WHERE organization_id=$1 AND deleted_at IS NULL AND state='PUBLISHED'
		ORDER BY publish_date DESC LIMIT 20 OFFSET $2`, orgID, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		rows.Scan(&p.ID, &p.State, &p.PublishDate, &p.OrganizationID, &p.IntegrationID,
			&p.Content, &p.Delay, &p.Group, &p.Title, &p.Description, &p.ParentPostID,
			&p.Settings, &p.Image, &p.CreationMethod, &p.Error, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}
	var total int
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE organization_id=$1 AND deleted_at IS NULL AND state='PUBLISHED'`, orgID).Scan(&total)
	return posts, total, nil
}

func (s *PostService) ValidatePosts(orgID string, posts []interface{}) map[string]interface{} {
	return map[string]interface{}{"valid": true, "errors": []interface{}{}}
}

func (s *PostService) GetMissing(orgID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"missing": []interface{}{}}, nil
}

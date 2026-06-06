package services

import (
	"context"
	"postpanda/backend-go/internal/database"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

type DashboardData struct {
	TotalPosts     int         `json:"totalPosts"`
	PublishedPosts int         `json:"publishedPosts"`
	ScheduledPosts int         `json:"scheduledPosts"`
	FailedPosts    int         `json:"failedPosts"`
	PostsByDay     interface{} `json:"postsByDay"`
	TopPlatforms   interface{} `json:"topPlatforms"`
}

func (s *AnalyticsService) GetDashboard(orgID, startDate, endDate string) (*DashboardData, error) {
	ctx := context.Background()

	data := &DashboardData{}

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND deleted_at IS NULL`,
		orgID,
	).Scan(&data.TotalPosts)

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND state = 'PUBLISHED' AND deleted_at IS NULL`,
		orgID,
	).Scan(&data.PublishedPosts)

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND state = 'QUEUE' AND deleted_at IS NULL`,
		orgID,
	).Scan(&data.ScheduledPosts)

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND state = 'ERROR' AND deleted_at IS NULL`,
		orgID,
	).Scan(&data.FailedPosts)

	// Posts per day
	rows, err := database.DB.Query(ctx,
		`SELECT DATE(publish_date), COUNT(*) FROM posts
		WHERE organization_id = $1 AND deleted_at IS NULL
		AND publish_date >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(publish_date) ORDER BY DATE(publish_date)`,
		orgID,
	)
	if err == nil {
		defer rows.Close()
		var byDay []map[string]interface{}
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			byDay = append(byDay, map[string]interface{}{"date": date, "count": count})
		}
		data.PostsByDay = byDay
	}

	return data, nil
}

func (s *AnalyticsService) GetIntegrationStats(orgID, integrationID, startDate, endDate string) (interface{}, error) {
	ctx := context.Background()

	var total, published, failed int
	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND integration_id = $2 AND deleted_at IS NULL`,
		orgID, integrationID,
	).Scan(&total)

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND integration_id = $2 AND state = 'PUBLISHED' AND deleted_at IS NULL`,
		orgID, integrationID,
	).Scan(&published)

	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE organization_id = $1 AND integration_id = $2 AND state = 'ERROR' AND deleted_at IS NULL`,
		orgID, integrationID,
	).Scan(&failed)

	return map[string]interface{}{
		"integrationId": integrationID,
		"total":         total,
		"published":     published,
		"failed":        failed,
	}, nil
}

func (s *AnalyticsService) GetPostStats(orgID, postID string) (interface{}, error) {
	// In production this would pull from each social platform's API
	return map[string]interface{}{
		"postId":    postID,
		"views":     0,
		"likes":     0,
		"comments":  0,
		"shares":    0,
		"reach":     0,
		"impressions": 0,
	}, nil
}

func (s *AnalyticsService) GetChannelOverview(orgID string) (interface{}, error) {
	ctx := context.Background()

	rows, err := database.DB.Query(ctx,
		`SELECT i.provider_identifier, COUNT(p.id) as post_count
		FROM integrations i
		LEFT JOIN posts p ON p.integration_id = i.id AND p.deleted_at IS NULL
		WHERE i.organization_id = $1 AND i.deleted_at IS NULL
		GROUP BY i.provider_identifier`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []map[string]interface{}
	for rows.Next() {
		var provider string
		var count int
		rows.Scan(&provider, &count)
		channels = append(channels, map[string]interface{}{
			"provider":  provider,
			"postCount": count,
		})
	}

	return channels, nil
}

func (s *AnalyticsService) GetScheduleAnalytics(orgID, startDate, endDate string) (interface{}, error) {
	ctx := context.Background()

	query := `SELECT DATE(publish_date) as date, state, COUNT(*) as count
		FROM posts WHERE organization_id = $1 AND deleted_at IS NULL`
	args := []interface{}{orgID}

	if startDate != "" {
		query += " AND publish_date >= $2"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND publish_date <= $3"
		args = append(args, endDate)
	}
	query += " GROUP BY DATE(publish_date), state ORDER BY date"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var date, state string
		var count int
		rows.Scan(&date, &state, &count)
		data = append(data, map[string]interface{}{
			"date":  date,
			"state": state,
			"count": count,
		})
	}

	return data, nil
}

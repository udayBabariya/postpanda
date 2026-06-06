package services

import (
	"context"
	"postpanda/backend-go/internal/auth"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AdminService struct{}

func NewAdminService() *AdminService { return &AdminService{} }

func (s *AdminService) GetErrors(platform, page string) (interface{}, int, error) {
	ctx := context.Background()
	p, _ := strconv.Atoi(page)
	if p < 1 { p = 1 }
	offset := (p - 1) * 20

	query := `SELECT id, message, platform, organization_id, post_id, body, created_at FROM errors WHERE 1=1`
	args := []interface{}{}
	if platform != "" {
		query += " AND platform=$1"
		args = append(args, platform)
	}
	query += " ORDER BY created_at DESC LIMIT 20 OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, offset)

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()

	var errors []map[string]interface{}
	for rows.Next() {
		var id, message, plat, orgID, postID, body string
		var createdAt time.Time
		rows.Scan(&id, &message, &plat, &orgID, &postID, &body, &createdAt)
		errors = append(errors, map[string]interface{}{"id": id, "message": message, "platform": plat, "organizationId": orgID, "postId": postID, "body": body, "createdAt": createdAt})
	}

	var total int
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM errors`).Scan(&total)
	return errors, total, nil
}

func (s *AdminService) GetErrorPlatforms() ([]string, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx, `SELECT DISTINCT platform FROM errors ORDER BY platform`)
	if err != nil { return nil, err }
	defer rows.Close()
	var platforms []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		platforms = append(platforms, p)
	}
	return platforms, nil
}

func (s *AdminService) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()
	stats := map[string]interface{}{}
	var totalUsers, totalOrgs, totalPosts, totalIntegrations, activeSubscriptions int
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM organizations`).Scan(&totalOrgs)
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL`).Scan(&totalPosts)
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM integrations WHERE deleted_at IS NULL`).Scan(&totalIntegrations)
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM subscriptions WHERE deleted_at IS NULL`).Scan(&activeSubscriptions)
	stats["totalUsers"] = totalUsers
	stats["totalOrgs"] = totalOrgs
	stats["totalPosts"] = totalPosts
	stats["totalIntegrations"] = totalIntegrations
	stats["activeSubscriptions"] = activeSubscriptions
	return stats, nil
}

func (s *AdminService) GetUsers(page, search string) (interface{}, int, error) {
	ctx := context.Background()
	p, _ := strconv.Atoi(page)
	if p < 1 { p = 1 }
	offset := (p - 1) * 20

	query := `SELECT id, email, name, provider_name, is_super_admin, activated, created_at FROM users WHERE 1=1`
	args := []interface{}{}
	if search != "" {
		query += " AND (email ILIKE $1 OR name ILIKE $1)"
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY created_at DESC LIMIT 20 OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, offset)

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var users []map[string]interface{}
	for rows.Next() {
		var id, email, provider string
		var name *string
		var isSuperAdmin, activated bool
		var createdAt time.Time
		rows.Scan(&id, &email, &name, &provider, &isSuperAdmin, &activated, &createdAt)
		users = append(users, map[string]interface{}{"id": id, "email": email, "name": name, "provider": provider, "isSuperAdmin": isSuperAdmin, "activated": activated, "createdAt": createdAt})
	}
	var total int
	database.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	return users, total, nil
}

func (s *AdminService) AddSubscription(orgID, tier, period string, channels int) error {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO subscriptions (id, organization_id, subscription_tier, period, total_channels, is_lifetime, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,false,$6,$7) ON CONFLICT (organization_id) DO UPDATE SET subscription_tier=$3, period=$4, total_channels=$5, deleted_at=NULL, updated_at=$7`,
		id, orgID, tier, period, channels, now, now)
	return err
}

func (s *AdminService) CancelSubscription(orgID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE subscriptions SET deleted_at=$1 WHERE organization_id=$2`, time.Now(), orgID)
	return err
}

func (s *AdminService) ImpersonateUser(userID string) (string, error) {
	userSvc := NewUserService()
	org, err := userSvc.GetPrimaryOrg(userID)
	if err != nil {
		return "", err
	}
	return auth.GenerateToken(userID, org.ID)
}

// ensure models package is used
var _ *models.Organization

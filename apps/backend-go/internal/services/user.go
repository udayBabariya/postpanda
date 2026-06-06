package services

import (
	"context"
	"fmt"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) FindByEmail(email string, provider models.Provider) (*models.User, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, email, password, provider_name, name, last_name, is_super_admin, bio, audience,
		picture_id, provider_id, timezone, last_read_notifications, invite_id, activated, account,
		connected_account, last_online, send_success_emails, send_failure_emails, send_streak_emails,
		created_at, updated_at
		FROM users WHERE email = $1 AND provider_name = $2`, email, provider)

	user := &models.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.ProviderName, &user.Name, &user.LastName,
		&user.IsSuperAdmin, &user.Bio, &user.Audience, &user.PictureID, &user.ProviderID,
		&user.Timezone, &user.LastReadNotifications, &user.InviteID, &user.Activated, &user.Account,
		&user.ConnectedAccount, &user.LastOnline, &user.SendSuccessEmails, &user.SendFailureEmails,
		&user.SendStreakEmails, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByID(id string) (*models.User, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, email, provider_name, name, last_name, is_super_admin, bio, audience,
		picture_id, timezone, last_read_notifications, activated, account, connected_account,
		last_online, send_success_emails, send_failure_emails, send_streak_emails, created_at, updated_at
		FROM users WHERE id = $1`, id)

	user := &models.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.ProviderName, &user.Name, &user.LastName,
		&user.IsSuperAdmin, &user.Bio, &user.Audience, &user.PictureID,
		&user.Timezone, &user.LastReadNotifications, &user.Activated, &user.Account,
		&user.ConnectedAccount, &user.LastOnline, &user.SendSuccessEmails,
		&user.SendFailureEmails, &user.SendStreakEmails, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateWithOrganization(email, hashedPassword, name string, provider models.Provider) (*models.User, *models.Organization, error) {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	userID := uuid.New().String()
	orgID := uuid.New().String()
	now := time.Now()

	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, email, password, provider_name, name, timezone, last_read_notifications,
		activated, last_online, send_success_emails, send_failure_emails, send_streak_emails, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, $6, true, $7, true, true, true, $8, $9)`,
		userID, email, hashedPassword, provider, name, now, now, now, now,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO organizations (id, name, shortlink, allow_trial, is_trailing, created_at, updated_at)
		VALUES ($1, $2, 'ASK', false, false, $3, $4)`,
		orgID, name+"'s Organization", now, now,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create organization: %w", err)
	}

	uoID := uuid.New().String()
	_, err = tx.Exec(ctx,
		`INSERT INTO user_organizations (id, user_id, organization_id, role, disabled, created_at, updated_at)
		VALUES ($1, $2, $3, 'ADMIN', false, $4, $5)`,
		uoID, userID, orgID, now, now,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user-org link: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	user := &models.User{
		ID:           userID,
		Email:        email,
		ProviderName: provider,
		Activated:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	org := &models.Organization{
		ID:        orgID,
		Name:      name + "'s Organization",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, org, nil
}

func (s *UserService) FindOrCreateOAuthUser(email, name, providerID string, provider models.Provider) (*models.User, *models.Organization, error) {
	user, err := s.FindByEmail(email, provider)
	if err == nil && user != nil {
		org, err := s.GetPrimaryOrg(user.ID)
		if err != nil {
			return nil, nil, err
		}
		return user, org, nil
	}

	// Create user without password for OAuth
	return s.createOAuthUser(email, name, providerID, provider)
}

func (s *UserService) createOAuthUser(email, name, providerID string, provider models.Provider) (*models.User, *models.Organization, error) {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	userID := uuid.New().String()
	orgID := uuid.New().String()
	now := time.Now()

	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, email, provider_name, provider_id, name, timezone, last_read_notifications,
		activated, last_online, send_success_emails, send_failure_emails, send_streak_emails, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, $6, true, $7, true, true, true, $8, $9)`,
		userID, email, provider, providerID, name, now, now, now, now,
	)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO organizations (id, name, shortlink, allow_trial, is_trailing, created_at, updated_at)
		VALUES ($1, $2, 'ASK', false, false, $3, $4)`,
		orgID, name+"'s Organization", now, now,
	)
	if err != nil {
		return nil, nil, err
	}

	uoID := uuid.New().String()
	_, err = tx.Exec(ctx,
		`INSERT INTO user_organizations (id, user_id, organization_id, role, disabled, created_at, updated_at)
		VALUES ($1, $2, $3, 'ADMIN', false, $4, $5)`,
		uoID, userID, orgID, now, now,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	user := &models.User{ID: userID, Email: email, ProviderName: provider}
	org := &models.Organization{ID: orgID}
	return user, org, nil
}

func (s *UserService) GetPrimaryOrg(userID string) (*models.Organization, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT o.id, o.name, o.description, o.api_key, o.payment_id, o.streak_since, o.shortlink,
		o.allow_trial, o.is_trailing, o.created_at, o.updated_at
		FROM organizations o
		JOIN user_organizations uo ON uo.organization_id = o.id
		WHERE uo.user_id = $1 AND uo.disabled = false
		ORDER BY uo.created_at ASC LIMIT 1`, userID)

	org := &models.Organization{}
	err := row.Scan(
		&org.ID, &org.Name, &org.Description, &org.APIKey, &org.PaymentID,
		&org.StreakSince, &org.ShortLink, &org.AllowTrial, &org.IsTrialing,
		&org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *UserService) FindOrgByID(orgID string) (*models.Organization, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, description, api_key, payment_id, streak_since, shortlink,
		allow_trial, is_trailing, created_at, updated_at
		FROM organizations WHERE id = $1`, orgID)

	org := &models.Organization{}
	err := row.Scan(
		&org.ID, &org.Name, &org.Description, &org.APIKey, &org.PaymentID,
		&org.StreakSince, &org.ShortLink, &org.AllowTrial, &org.IsTrialing,
		&org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *UserService) GetOrgMembership(userID, orgID string) (*models.UserOrganization, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, organization_id, disabled, role, created_at, updated_at
		FROM user_organizations WHERE user_id = $1 AND organization_id = $2 AND disabled = false`,
		userID, orgID)

	uo := &models.UserOrganization{}
	err := row.Scan(&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Disabled, &uo.Role, &uo.CreatedAt, &uo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return uo, nil
}

func (s *UserService) UpdateProfile(userID string, updates map[string]interface{}) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET name = $1, last_name = $2, bio = $3, timezone = $4, updated_at = $5
		WHERE id = $6`,
		updates["name"], updates["lastName"], updates["bio"], updates["timezone"], time.Now(), userID,
	)
	return err
}

func (s *UserService) SendPasswordReset(email string) error {
	// TODO: implement with email service
	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	ctx := context.Background()
	// token lookup from Redis
	userID, err := database.Redis.Get(ctx, "reset:"+token).Result()
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	_, err = database.DB.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", string(hashed), userID)
	database.Redis.Del(ctx, "reset:"+token)
	return err
}

func (s *UserService) ActivateAccount(code string) error {
	ctx := context.Background()
	userID, err := database.Redis.Get(ctx, "activate:"+code).Result()
	if err != nil {
		return fmt.Errorf("invalid code")
	}
	_, err = database.DB.Exec(ctx, "UPDATE users SET activated = true WHERE id = $1", userID)
	database.Redis.Del(ctx, "activate:"+code)
	return err
}

func (s *UserService) UpdateEmailPreferences(userID string, success, failure, streak bool) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET send_success_emails=$1, send_failure_emails=$2, send_streak_emails=$3, updated_at=$4 WHERE id=$5`,
		success, failure, streak, time.Now(), userID)
	return err
}

func (s *UserService) GetOrgMembers(orgID string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT u.id, u.email, u.name, u.last_name, u.picture_id, uo.role, uo.disabled
		FROM users u JOIN user_organizations uo ON uo.user_id=u.id
		WHERE uo.organization_id=$1 ORDER BY uo.created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []map[string]interface{}
	for rows.Next() {
		var id, email, role string
		var name, lastName, pictureID *string
		var disabled bool
		rows.Scan(&id, &email, &name, &lastName, &pictureID, &role, &disabled)
		members = append(members, map[string]interface{}{
			"id": id, "email": email, "name": name, "lastName": lastName,
			"pictureId": pictureID, "role": role, "disabled": disabled,
		})
	}
	return members, nil
}

func (s *UserService) InviteMember(orgID, email, role string) error {
	// TODO: send invite email and create pending user entry
	return nil
}

func (s *UserService) RemoveMember(orgID, memberID string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE user_organizations SET disabled=true, updated_at=$1 WHERE user_id=$2 AND organization_id=$3`,
		time.Now(), memberID, orgID)
	return err
}

func (s *UserService) UpdateMemberRole(orgID, memberID, role string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE user_organizations SET role=$1, updated_at=$2 WHERE user_id=$3 AND organization_id=$4`,
		role, time.Now(), memberID, orgID)
	return err
}

func (s *UserService) CreateOrg(userID, name string) (*models.Organization, error) {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	orgID := uuid.New().String()
	uoID := uuid.New().String()
	now := time.Now()

	_, err = tx.Exec(ctx,
		`INSERT INTO organizations (id, name, shortlink, allow_trial, is_trailing, created_at, updated_at)
		VALUES ($1, $2, 'ASK', false, false, $3, $4)`, orgID, name, now, now)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_organizations (id, user_id, organization_id, role, disabled, created_at, updated_at)
		VALUES ($1, $2, $3, 'ADMIN', false, $4, $5)`, uoID, userID, orgID, now, now)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Organization{ID: orgID, Name: name, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *UserService) ResendActivation(email string) {
	// TODO: send activation email
}

func (s *UserService) CheckOAuthExists(providerID string, provider models.Provider) bool {
	ctx := context.Background()
	var count int
	database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE provider_id=$1 AND provider_name=$2`, providerID, provider).Scan(&count)
	return count > 0
}

func (s *UserService) ListOrgs(userID string) ([]*models.Organization, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT o.id, o.name, o.description, o.shortlink, o.allow_trial, o.is_trailing, o.created_at, o.updated_at
		FROM organizations o
		JOIN user_organizations uo ON uo.organization_id = o.id
		WHERE uo.user_id = $1 AND uo.disabled = false
		ORDER BY uo.created_at ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*models.Organization
	for rows.Next() {
		org := &models.Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.ShortLink,
			&org.AllowTrial, &org.IsTrialing, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

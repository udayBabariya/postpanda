package services

import (
	"context"
	"encoding/json"
	"fmt"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/integrations"
	"postpanda/backend-go/internal/models"
	"time"
)

type IntegrationService struct{}

func NewIntegrationService() *IntegrationService {
	return &IntegrationService{}
}

func (s *IntegrationService) List(orgID string) ([]*models.Integration, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, internal_id, organization_id, name, picture, provider_identifier, type,
		disabled, token_expiration, profile, deleted_at, in_between_steps, refresh_needed,
		posting_times, custom_instance_details, customer_id, root_internal_id, additional_settings,
		created_at, updated_at
		FROM integrations
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Integration
	for rows.Next() {
		i := &models.Integration{}
		if err := rows.Scan(
			&i.ID, &i.InternalID, &i.OrganizationID, &i.Name, &i.Picture,
			&i.ProviderIdentifier, &i.Type, &i.Disabled, &i.TokenExpiration,
			&i.Profile, &i.DeletedAt, &i.InBetweenSteps, &i.RefreshNeeded,
			&i.PostingTimes, &i.CustomInstanceDetails, &i.CustomerID,
			&i.RootInternalID, &i.AdditionalSettings, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func (s *IntegrationService) GetByID(orgID, id string) (*models.Integration, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, internal_id, organization_id, name, picture, provider_identifier, type,
		disabled, token_expiration, profile, deleted_at, in_between_steps, refresh_needed,
		posting_times, custom_instance_details, customer_id, root_internal_id, additional_settings,
		created_at, updated_at
		FROM integrations
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	)

	i := &models.Integration{}
	err := row.Scan(
		&i.ID, &i.InternalID, &i.OrganizationID, &i.Name, &i.Picture,
		&i.ProviderIdentifier, &i.Type, &i.Disabled, &i.TokenExpiration,
		&i.Profile, &i.DeletedAt, &i.InBetweenSteps, &i.RefreshNeeded,
		&i.PostingTimes, &i.CustomInstanceDetails, &i.CustomerID,
		&i.RootInternalID, &i.AdditionalSettings, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *IntegrationService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE integrations SET deleted_at = $1 WHERE id = $2 AND organization_id = $3`,
		time.Now(), id, orgID,
	)
	return err
}

func (s *IntegrationService) Disable(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE integrations SET disabled = true WHERE id = $1 AND organization_id = $2`,
		id, orgID,
	)
	return err
}

func (s *IntegrationService) Enable(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE integrations SET disabled = false WHERE id = $1 AND organization_id = $2`,
		id, orgID,
	)
	return err
}

func (s *IntegrationService) UpdatePostingTimes(orgID, id, postingTimes string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx,
		`UPDATE integrations SET posting_times = $1 WHERE id = $2 AND organization_id = $3`,
		postingTimes, id, orgID,
	)
	return err
}

func (s *IntegrationService) RefreshToken(orgID, id string) error {
	integration, err := s.GetByID(orgID, id)
	if err != nil {
		return err
	}

	provider := integrations.GetProvider(integration.ProviderIdentifier)
	if provider == nil {
		return fmt.Errorf("unknown provider: %s", integration.ProviderIdentifier)
	}

	newToken, expiry, err := provider.RefreshToken(integration.Token, integration.RefreshToken)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = database.DB.Exec(ctx,
		`UPDATE integrations SET token = $1, token_expiration = $2, refresh_needed = false
		WHERE id = $3 AND organization_id = $4`,
		newToken, expiry, id, orgID,
	)
	return err
}

func (s *IntegrationService) UpdateSettings(orgID, id string, settings map[string]interface{}) error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = database.DB.Exec(ctx,
		`UPDATE integrations SET additional_settings = $1 WHERE id = $2 AND organization_id = $3`,
		string(b), id, orgID,
	)
	return err
}

func (s *IntegrationService) GetAvailableProviders() []map[string]interface{} {
	return integrations.GetAllProvidersMeta()
}

func (s *IntegrationService) GetOAuthURL(orgID, providerName string) (string, error) {
	provider := integrations.GetProvider(providerName)
	if provider == nil {
		return "", fmt.Errorf("unknown provider: %s", providerName)
	}
	return provider.GetOAuthURL(orgID), nil
}

func (s *IntegrationService) HandleCallback(orgID, providerName, code, state string) (*models.Integration, error) {
	provider := integrations.GetProvider(providerName)
	if provider == nil {
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}

	info, err := provider.HandleCallback(code, state)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	now := time.Now()

	var existing models.Integration
	err = database.DB.QueryRow(ctx,
		`SELECT id FROM integrations WHERE organization_id = $1 AND internal_id = $2`,
		orgID, info.InternalID,
	).Scan(&existing.ID)

	if err == nil {
		// Update existing
		database.DB.Exec(ctx,
			`UPDATE integrations SET token = $1, refresh_token = $2, token_expiration = $3,
			name = $4, picture = $5, disabled = false, refresh_needed = false, updated_at = $6
			WHERE id = $7`,
			info.Token, info.RefreshToken, info.TokenExpiration,
			info.Name, info.Picture, now, existing.ID,
		)
		return s.GetByID(orgID, existing.ID)
	}

	// Create new
	newID := fmt.Sprintf("int_%s", now.Format("20060102150405"))
	_, err = database.DB.Exec(ctx,
		`INSERT INTO integrations (id, internal_id, organization_id, name, picture, provider_identifier,
		type, token, refresh_token, token_expiration, profile, disabled, in_between_steps, refresh_needed,
		posting_times, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false, false, false,
		'[{"time":120}, {"time":400}, {"time":700}]', $12, $13)`,
		newID, info.InternalID, orgID, info.Name, info.Picture, providerName,
		info.Type, info.Token, info.RefreshToken, info.TokenExpiration, info.Profile,
		now, now,
	)
	if err != nil {
		return nil, err
	}

	return s.GetByID(orgID, newID)
}

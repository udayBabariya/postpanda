package services

import (
	"context"
	"postpanda/backend-go/internal/database"
	"time"

	"github.com/google/uuid"
)

type PlugsService struct{}

func NewPlugsService() *PlugsService { return &PlugsService{} }

type Plug struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	IntegrationID  string `json:"integrationId"`
	PlugFunction   string `json:"plugFunction"`
	Data           string `json:"data"`
	Activated      bool   `json:"activated"`
}

func (s *PlugsService) ListByOrg(orgID string) ([]*Plug, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, integration_id, plug_function, data, activated FROM plugs WHERE organization_id=$1`, orgID)
	if err != nil { return nil, err }
	defer rows.Close()
	var plugs []*Plug
	for rows.Next() {
		p := &Plug{}
		rows.Scan(&p.ID, &p.OrganizationID, &p.IntegrationID, &p.PlugFunction, &p.Data, &p.Activated)
		plugs = append(plugs, p)
	}
	return plugs, nil
}

func (s *PlugsService) ListByIntegration(orgID, integrationID string) ([]*Plug, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, organization_id, integration_id, plug_function, data, activated FROM plugs WHERE organization_id=$1 AND integration_id=$2`, orgID, integrationID)
	if err != nil { return nil, err }
	defer rows.Close()
	var plugs []*Plug
	for rows.Next() {
		p := &Plug{}
		rows.Scan(&p.ID, &p.OrganizationID, &p.IntegrationID, &p.PlugFunction, &p.Data, &p.Activated)
		plugs = append(plugs, p)
	}
	return plugs, nil
}

func (s *PlugsService) Upsert(orgID, integrationID, plugFunction, data string) (*Plug, error) {
	ctx := context.Background()
	id := uuid.New().String()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO plugs (id, organization_id, integration_id, plug_function, data, activated)
		VALUES ($1,$2,$3,$4,$5,true)
		ON CONFLICT (plug_function, integration_id) DO UPDATE SET data=$5, activated=true`,
		id, orgID, integrationID, plugFunction, data)
	if err != nil { return nil, err }
	return &Plug{ID: id, OrganizationID: orgID, IntegrationID: integrationID, PlugFunction: plugFunction, Data: data, Activated: true}, nil
}

func (s *PlugsService) Activate(id string, activated bool) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE plugs SET activated=$1 WHERE id=$2`, activated, id)
	return err
}

// Customer service

type CustomerService struct{}

func NewCustomerService() *CustomerService { return &CustomerService{} }

type Customer struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	OrgID     string  `json:"orgId"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (s *CustomerService) List(orgID string) ([]*Customer, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, org_id, deleted_at, created_at FROM customers WHERE org_id=$1 AND deleted_at IS NULL ORDER BY name ASC`, orgID)
	if err != nil { return nil, err }
	defer rows.Close()
	var customers []*Customer
	for rows.Next() {
		c := &Customer{}
		rows.Scan(&c.ID, &c.Name, &c.OrgID, &c.DeletedAt, &c.CreatedAt)
		customers = append(customers, c)
	}
	return customers, nil
}

func (s *CustomerService) Create(orgID, name string) (*Customer, error) {
	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	_, err := database.DB.Exec(ctx,
		`INSERT INTO customers (id, name, org_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`,
		id, name, orgID, now, now)
	if err != nil { return nil, err }
	return &Customer{ID: id, Name: name, OrgID: orgID, CreatedAt: now}, nil
}

func (s *CustomerService) Delete(orgID, id string) error {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, `UPDATE customers SET deleted_at=$1 WHERE id=$2 AND org_id=$3`, time.Now(), id, orgID)
	return err
}

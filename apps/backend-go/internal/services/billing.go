package services

import (
	"context"
	"fmt"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"time"

	"github.com/stripe/stripe-go/v81"
	portalSession "github.com/stripe/stripe-go/v81/billingportal/session"
	stripesession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
)

type BillingService struct{}

func NewBillingService() *BillingService {
	stripe.Key = config.App.StripeSecretKey
	return &BillingService{}
}

func (s *BillingService) GetSubscription(orgID string) (*models.Subscription, error) {
	ctx := context.Background()
	row := database.DB.QueryRow(ctx,
		`SELECT id, organization_id, subscription_tier, identifier, cancel_at, period,
		total_channels, is_lifetime, deleted_at, created_at, updated_at
		FROM subscriptions WHERE organization_id=$1 AND deleted_at IS NULL`, orgID)

	sub := &models.Subscription{}
	err := row.Scan(
		&sub.ID, &sub.OrganizationID, &sub.Tier, &sub.Identifier, &sub.CancelAt,
		&sub.Period, &sub.TotalChannels, &sub.IsLifetime, &sub.DeletedAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *BillingService) GetOrCreateCustomer(orgID, email string) (string, error) {
	ctx := context.Background()
	var paymentID string
	database.DB.QueryRow(ctx, `SELECT payment_id FROM organizations WHERE id=$1`, orgID).Scan(&paymentID)

	if paymentID != "" {
		return paymentID, nil
	}

	params := &stripe.CustomerParams{Email: stripe.String(email)}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}

	database.DB.Exec(ctx, `UPDATE organizations SET payment_id=$1 WHERE id=$2`, c.ID, orgID)
	return c.ID, nil
}

func (s *BillingService) CreateCheckoutSession(orgID, userID, plan, period string) (string, error) {
	ctx := context.Background()
	var email string
	database.DB.QueryRow(ctx, `SELECT u.email FROM users u
		JOIN user_organizations uo ON uo.user_id=u.id
		WHERE uo.organization_id=$1 AND uo.role='ADMIN' LIMIT 1`, orgID).Scan(&email)

	customerID, err := s.GetOrCreateCustomer(orgID, email)
	if err != nil {
		return "", err
	}

	priceID := config.App.StripeMonthlyPriceID
	if period == "yearly" {
		priceID = config.App.StripeYearlyPriceID
	}

	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(customerID),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(config.App.FrontendURL + "/billing?success=true"),
		CancelURL:  stripe.String(config.App.FrontendURL + "/billing?canceled=true"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"orgId": orgID,
			"plan":  plan,
		},
	}

	sess, err := stripesession.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

func (s *BillingService) CreatePortalSession(orgID string) (string, error) {
	ctx := context.Background()
	var paymentID string
	database.DB.QueryRow(ctx, `SELECT payment_id FROM organizations WHERE id=$1`, orgID).Scan(&paymentID)

	if paymentID == "" {
		return "", fmt.Errorf("no stripe customer found")
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(paymentID),
		ReturnURL: stripe.String(config.App.FrontendURL + "/billing"),
	}

	sess, err := portalSession.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

func (s *BillingService) HandleWebhookEvent(event stripe.Event) error {
	ctx := context.Background()
	switch event.Type {
	case "checkout.session.completed":
		// Handle subscription creation
	case "customer.subscription.updated":
		// Handle subscription update
	case "customer.subscription.deleted":
		// Handle subscription deletion - downgrade to free
		var paymentID string
		// Extract customer ID from event
		_ = paymentID
		database.DB.Exec(ctx,
			`UPDATE subscriptions SET deleted_at=$1 WHERE identifier=$2`, time.Now(), "")
	}
	return nil
}

func (s *BillingService) CancelSubscription(orgID string) error {
	ctx := context.Background()
	var subID string
	database.DB.QueryRow(ctx, `SELECT identifier FROM subscriptions WHERE organization_id=$1 AND deleted_at IS NULL`, orgID).Scan(&subID)

	// Cancel via Stripe
	_ = subID

	_, err := database.DB.Exec(ctx,
		`UPDATE subscriptions SET cancel_at=$1 WHERE organization_id=$2`, time.Now().Add(30*24*time.Hour), orgID)
	return err
}

func (s *BillingService) GetPlans() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": "standard", "name": "Standard", "channels": 5, "monthlyPrice": 29, "yearlyPrice": 290},
		{"id": "pro", "name": "Pro", "channels": 15, "monthlyPrice": 59, "yearlyPrice": 590},
		{"id": "team", "name": "Team", "channels": 50, "monthlyPrice": 149, "yearlyPrice": 1490},
		{"id": "ultimate", "name": "Ultimate", "channels": -1, "monthlyPrice": 399, "yearlyPrice": 3990},
	}
}

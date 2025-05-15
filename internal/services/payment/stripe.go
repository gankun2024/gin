package payment

import (
	"encoding/json"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"

	"github.com/gankun2024/gin-demo-project/internal/config"
	"github.com/gankun2024/gin-demo-project/internal/db/models"
)

// CheckoutSession represents a Stripe checkout session
type CheckoutSession struct {
	ID  string
	URL string
}

// Initialize Stripe with API key
func init() {
	cfg, err := config.Load()
	if err != nil {
		// Log error but don't panic - we might be in a test environment
		return
	}
	stripe.Key = cfg.Stripe.SecretKey
}

// CreateCheckoutSession creates a new Stripe checkout session
func CreateCheckoutSession(
	userID string,
	priceID string,
	successURL string,
	cancelURL string,
	customerName string,
	customerEmail string,
) (*CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		ClientReferenceID: stripe.String(userID),
	}

	// Add customer details if provided
	if customerEmail != "" {
		params.CustomerEmail = stripe.String(customerEmail)
	}

	// Create session
	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	// Store session in database
	paymentSession := &models.PaymentSession{
		ID:          s.ID,
		UserID:      userID,
		PriceID:     priceID,
		Status:      string(s.Status),
		PaymentType: string(s.Mode),
		CreatedAt:   s.Created,
		// ExpiresAt:    s.Expires,
		SuccessURL:   successURL,
		CancelURL:    cancelURL,
		CustomerName: customerName,
	}

	if err := models.CreatePaymentSession(paymentSession); err != nil {
		return nil, err
	}

	return &CheckoutSession{
		ID:  s.ID,
		URL: s.URL,
	}, nil
}

// GetSession retrieves a Stripe checkout session
func GetSession(sessionID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{}
	s, err := session.Get(sessionID, params)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ProcessWebhook processes a Stripe webhook event
func ProcessWebhook(payload []byte, signature string) (*stripe.Event, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, signature, cfg.Stripe.WebhookSecret)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// HandleSessionCompleted handles the checkout.session.completed event
func HandleSessionCompleted(event *stripe.Event) error {
	var checkoutSession stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &checkoutSession)
	if err != nil {
		return err
	}

	// Update session status in database
	err = models.UpdatePaymentSessionStatus(checkoutSession.ID, string(checkoutSession.Status))
	if err != nil {
		return err
	}

	// Create a record of the payment
	payment := &models.Payment{
		SessionID:  checkoutSession.ID,
		UserID:     checkoutSession.ClientReferenceID,
		Amount:     checkoutSession.AmountTotal,
		Currency:   string(checkoutSession.Currency),
		Status:     "completed",
		PaymentID:  checkoutSession.PaymentIntent.ID,
		CreatedAt:  checkoutSession.Created,
		CustomerID: checkoutSession.Customer.ID,
	}

	return models.CreatePayment(payment)
}

// HandleSubscriptionCreated handles the customer.subscription.created event
func HandleSubscriptionCreated(event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	// Here you would update the user's subscription status in your database
	// For simplicity, we're just returning nil
	return nil
}

// HandleSubscriptionUpdated handles the customer.subscription.updated event
func HandleSubscriptionUpdated(event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	// Here you would update the user's subscription status in your database
	// For simplicity, we're just returning nil
	return nil
}

// HandleSubscriptionCanceled handles the customer.subscription.deleted event
func HandleSubscriptionCanceled(event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	// Here you would update the user's subscription status in your database
	// For simplicity, we're just returning nil
	return nil
}

// HandleInvoicePaid handles the invoice.paid event
func HandleInvoicePaid(event *stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return err
	}

	// Here you would record the payment in your database
	// For simplicity, we're just returning nil
	return nil
}

// HandleInvoicePaymentFailed handles the invoice.payment_failed event
func HandleInvoicePaymentFailed(event *stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return err
	}

	// Here you would update the payment status in your database
	// For simplicity, we're just returning nil
	return nil
}

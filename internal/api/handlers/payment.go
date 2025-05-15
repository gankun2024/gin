package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gankun2024/gin-demo-project/internal/services/payment"
)

// CreateCheckoutSessionRequest represents the request for creating a checkout session
type CreateCheckoutSessionRequest struct {
	PriceID       string `json:"price_id" binding:"required"`
	SuccessURL    string `json:"success_url" binding:"required"`
	CancelURL     string `json:"cancel_url" binding:"required"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
}

// CreateCheckoutSession creates a new Stripe checkout session
func CreateCheckoutSession(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create checkout session
	session, err := payment.CreateCheckoutSession(
		userID,
		req.PriceID,
		req.SuccessURL,
		req.CancelURL,
		req.CustomerName,
		req.CustomerEmail,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"checkout_url": session.URL,
	})
}

// PaymentSuccess handles successful payments
func PaymentSuccess(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session_id"})
		return
	}

	// Get payment session details
	session, err := payment.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"session": session,
	})
}

// PaymentCancel handles canceled payments
func PaymentCancel(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "canceled",
	})
}

// StripeWebhook handles Stripe webhook events
func StripeWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Get the Stripe-Signature header
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
		return
	}

	// Process the webhook
	event, err := payment.ProcessWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		// Payment is successful
		err = payment.HandleSessionCompleted(event)
	case "customer.subscription.created":
		// Subscription created
		err = payment.HandleSubscriptionCreated(event)
	case "customer.subscription.updated":
		// Subscription updated
		err = payment.HandleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		// Subscription canceled
		err = payment.HandleSubscriptionCanceled(event)
	case "invoice.paid":
		// Invoice paid
		err = payment.HandleInvoicePaid(event)
	case "invoice.payment_failed":
		// Payment failed
		err = payment.HandleInvoicePaymentFailed(event)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

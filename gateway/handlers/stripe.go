package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
	"gorm.io/gorm"
)

func createCheckoutSession(walletAddress string) (*stripe.CheckoutSession, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(frontendURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1OehLu2mES9P7kjwSQS45ZKq"), // comes from Stripe Product Dashboard
				Quantity: stripe.Int64(1),
				AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
					Enabled: stripe.Bool(true),
				},
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"walletAddress": walletAddress,
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	params.AddMetadata("walletAddress", walletAddress)
	result, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func StripeCreateCheckoutSessionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			http.Error(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		session, err := createCheckoutSession(user.WalletAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
	}
}

func StripeFullfillmentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")

		event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "payment_intent.succeeded":
			var paymentIntent stripe.PaymentIntent
			err := json.Unmarshal(event.Data.Raw, &paymentIntent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing payment intent: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			walletAddress, ok := paymentIntent.Metadata["walletAddress"]
			if !ok {
				fmt.Fprintf(os.Stderr, "Wallet address not found in payment intent metadata\n")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Fetch the User by walletAddress
			var user models.User
			result := db.Where("wallet_address ILIKE ?", walletAddress).First(&user)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					fmt.Fprintf(os.Stderr, "User with wallet address %s not found\n", walletAddress)
					w.WriteHeader(http.StatusNotFound)
				} else {
					fmt.Fprintf(os.Stderr, "Error querying user: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			// Convert Stripe amount (in cents) to a float64 representation
			amount := float64(paymentIntent.Amount) / 100.0

			// Create a new Transaction
			transaction := models.Transaction{
				ID:          uuid.New().String(),
				Amount:      amount,
				IsDebit:     false, // Assuming payment intents are always credits
				UserID:      user.WalletAddress,
				Description: "Stripe payment intent succeeded",
			}

			// Save the Transaction to the database
			if result := db.Create(&transaction); result.Error != nil {
				fmt.Fprintf(os.Stderr, "Failed to save transaction: %v\n", result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Printf("PaymentIntent succeeded, Amount: %v, WalletAddress: %v\n", paymentIntent.Amount, walletAddress)

		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}

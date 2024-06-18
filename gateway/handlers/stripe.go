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
	"github.com/labdao/plex/gateway/utils"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentmethod"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/webhook"
	"gorm.io/gorm"
)

func setupStripeClient() error {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return errors.New("STRIPE_SECRET_KEY environment variable not set")
	}
	stripe.Key = apiKey
	return nil
}

func createStripeCustomer(walletAddress string) (string, error) {
	err := setupStripeClient()
	if err != nil {
		return "", err
	}

	params := &stripe.CustomerParams{
		Name: stripe.String(walletAddress),
	}
	customer, err := customer.New(params)
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func getCustomerPaymentMethod(stripeUserID string) (*stripe.PaymentMethod, error) {
	err := setupStripeClient()
	if err != nil {
		return nil, err
	}

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeUserID),
		Type:     stripe.String(string(stripe.PaymentMethodTypeCard)),
	}
	i := paymentmethod.List(params)

	for i.Next() {
		return i.PaymentMethod(), nil
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func createCheckoutSession(stripeUserID string, computeCost int) (*stripe.CheckoutSession, error) {
	err := setupStripeClient()
	if err != nil {
		return nil, err
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	priceParams := &stripe.PriceParams{
		UnitAmount: stripe.Int64(int64(computeCost * 10)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Product:    stripe.String(os.Getenv("STRIPE_PRODUCT_ID")),
	}
	price, err := price.New(priceParams)
	if err != nil {
		return nil, err
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeUserID),
		// TODO: success url needs to be accessible to user, not just the backend
		SuccessURL: stripe.String(frontendURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(price.ID),
				Quantity: stripe.Int64(1),
				AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
					Enabled: stripe.Bool(false),
				},
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"Stripe User ID": stripeUserID,
			},
			SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOnSession)),
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	params.AddMetadata("Stripe User ID", stripeUserID)

	session, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func StripeCreateCheckoutSessionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			utils.SendJSONError(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		// TODO: pass in tool

		// toolID := r.URL.Query().Get("toolID")
		// if toolID == "" {
		// 	utils.SendJSONError(w, "Tool ID not provided", http.StatusBadRequest)
		// 	return
		// }

		// var tool models.Tool
		// result := db.Where("cid = ?", toolID).First(&tool)
		// if result.Error != nil {
		// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 		utils.SendJSONError(w, "Tool not found", http.StatusNotFound)
		// 	} else {
		// 		utils.SendJSONError(w, fmt.Sprintf("Error fetching Tool: %v", result.Error), http.StatusInternalServerError)
		// 	}
		// 	return
		// }

		// TODO: modify so we're not passing in hardcoded value
		session, err := createCheckoutSession(user.StripeUserID, 10)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println("Checkout URL:", session.URL)

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

			transaction := models.Transaction{
				ID:          uuid.New().String(),
				Amount:      amount,
				IsDebit:     false, // Assuming payment intents are always credits (money in)
				UserID:      user.WalletAddress,
				Description: "Stripe payment intent succeeded",
			}

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

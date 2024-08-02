package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meter"
	"github.com/stripe/stripe-go/v78/billing/metereventsummary"
	billingportal "github.com/stripe/stripe-go/v78/billingportal/session"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/product"
	"github.com/stripe/stripe-go/v78/subscription"
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

func createCheckoutSession(stripeUserID, walletAddress, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	err := setupStripeClient()
	if err != nil {
		return nil, err
	}

	priceID := os.Getenv("STRIPE_PRICE_ID")
	if priceID == "" {
		return nil, errors.New("STRIPE_PRICE_ID environment variable not set")
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeUserID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(priceID),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"Wallet Address": walletAddress,
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		Metadata: map[string]string{
			"Wallet Address": walletAddress,
		},
	}

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

		// thresholdStr := os.Getenv("TIER_THRESHOLD")
		// threshold, _ := strconv.Atoi(thresholdStr)

		// if user.ComputeTally < threshold || user.SubscriptionStatus == "active" {
		// 	utils.SendJSONError(w, "User does not need a subscription at this time", http.StatusBadRequest)
		// 	return
		// }

		var requestBody struct {
			SuccessURL string `json:"success_url"`
			CancelURL  string `json:"cancel_url"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			utils.SendJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		session, err := createCheckoutSession(user.StripeUserID, user.WalletAddress, requestBody.SuccessURL, requestBody.CancelURL)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println("Checkout URL:", session.URL)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
	}
}

func StripeFulfillmentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")

		event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), endpointSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
			var subscription stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &subscription)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing subscription: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			walletAddress, ok := subscription.Metadata["Wallet Address"]
			if !ok {
				fmt.Fprintf(os.Stderr, "Wallet Address not found in subscription metadata\n")
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

			if subscription.Status == "active" {
				user.SubscriptionStatus = "active"
			} else {
				user.SubscriptionStatus = string(subscription.Status)
			}
			user.SubscriptionID = &subscription.ID
			result = db.Save(&user)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error updating user subscription status: %v\n", result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Printf("Subscription %s for user %s updated to %s\n", subscription.ID, walletAddress, user.SubscriptionStatus)

		// PR#1010 this might not be relevant anymore as we are not providing a trial period
		case "customer.subscription.trial_will_end":
			// Handle trial ending soon (e.g., send notification to user)
			// This event occurs 3 days before the trial ends
			// You might want to implement logic to notify the user

		case "invoice.paid", "invoice.payment_failed":
			// Handle successful or failed payments
			// You might want to update the user's payment status or send notifications

		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func StripeGetPlanDetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := setupStripeClient()
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error setting up Stripe client: %v", err), http.StatusInternalServerError)
			return
		}

		// Fetch the default price ID from environment variables
		priceID := os.Getenv("STRIPE_PRICE_ID")
		if priceID == "" {
			utils.SendJSONError(w, "STRIPE_PRICE_ID environment variable not set", http.StatusInternalServerError)
			return
		}

		priceParams := &stripe.PriceParams{
			Expand: []*string{stripe.String("tiers")},
		}

		// Get the price details
		priceDetails, err := price.Get(priceID, priceParams)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting price details: %v", err), http.StatusInternalServerError)
			return
		}

		// Get the associated product details
		productDetails, err := product.Get(priceDetails.Product.ID, nil)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting product details: %v", err), http.StatusInternalServerError)
			return
		}

		var flatFee float64
		tiers := make([]map[string]interface{}, 0)

		if priceDetails.TiersMode == "graduated" && len(priceDetails.Tiers) > 0 {
			for _, tier := range priceDetails.Tiers {
				tierInfo := map[string]interface{}{
					"up_to":       tier.UpTo,
					"unit_amount": tier.UnitAmountDecimal / 100, // Using UnitAmountDecimal for precision
				}
				tiers = append(tiers, tierInfo)
			}

			// Set included credits and overage charge based on tiers
			if len(tiers) > 0 {
				flatFee = float64(priceDetails.Tiers[0].FlatAmount) / 100 // Set the flat fee
			}
		} else if priceDetails.UnitAmount != 0 {
			// Handle non-tiered flat pricing
			flatFee = float64(priceDetails.UnitAmount) / 100
		} else {
			fmt.Println("Tiers or flat pricing not configured properly")
		}

		// Prepare the response
		response := map[string]interface{}{
			"plan_name":        productDetails.Name,
			"plan_amount":      flatFee,
			"plan_currency":    priceDetails.Currency,
			"plan_interval":    priceDetails.Recurring.Interval,
			"included_credits": 500,  // Replace with your logic
			"overage_charge":   0.01, // Replace with your logic
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func StripeGetSubscriptionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			utils.SendJSONError(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		if user.SubscriptionID == nil {
			utils.SendJSONError(w, "User does not have an active subscription", http.StatusBadRequest)
			return
		}

		err := setupStripeClient()
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error setting up Stripe client: %v", err), http.StatusInternalServerError)
			return
		}

		subscription, err := subscription.Get(*user.SubscriptionID, nil)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting subscription: %v", err), http.StatusInternalServerError)
			return
		}

		if len(subscription.Items.Data) == 0 {
			utils.SendJSONError(w, "No subscription items found", http.StatusInternalServerError)
			return
		}

		item := subscription.Items.Data[0]
		plan := item.Plan

		product, err := product.Get(plan.Product.ID, nil)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting product: %v", err), http.StatusInternalServerError)
			return
		}

		// Fetch price details
		priceParams := &stripe.PriceParams{
			Expand: []*string{stripe.String("tiers")},
		}
		priceDetails, err := price.Get(plan.ID, priceParams)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting price details: %v", err), http.StatusInternalServerError)
			return
		}

		var includedCredits int
		var overageCharge float64
		var flatFee float64
		tiers := make([]map[string]interface{}, 0)

		if priceDetails.TiersMode == "graduated" && len(priceDetails.Tiers) > 0 {
			for _, tier := range priceDetails.Tiers {
				tierInfo := map[string]interface{}{
					"up_to":       tier.UpTo,
					"unit_amount": tier.UnitAmountDecimal / 100, // Using UnitAmountDecimal for precision
				}
				tiers = append(tiers, tierInfo)
			}

			// Set included credits and overage charge based on tiers
			if len(tiers) > 0 {
				includedCredits = int(priceDetails.Tiers[0].UpTo)
				if len(priceDetails.Tiers) > 1 {
					overageCharge = float64(priceDetails.Tiers[1].UnitAmountDecimal) / 100
				}
				flatFee = float64(priceDetails.Tiers[0].FlatAmount) / 100 // Set the flat fee
			}
		} else if priceDetails.UnitAmount != 0 {
			// Handle non-tiered flat pricing
			flatFee = float64(priceDetails.UnitAmount) / 100
		} else {
			fmt.Println("Tiers or flat pricing not configured properly")
		}

		// Fetch used credits
		usedCredits, err := fetchUsedCredits(user.StripeUserID, "compute_units", subscription)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting usage records: %v", err), http.StatusInternalServerError)
			return
		}

		// Prepare the response
		response := map[string]interface{}{
			"plan_name":            product.Name,
			"plan_amount":          flatFee, // Use the calculated flat fee
			"plan_currency":        priceDetails.Currency,
			"plan_interval":        plan.Interval,
			"current_period_start": time.Unix(subscription.CurrentPeriodStart, 0),
			"current_period_end":   time.Unix(subscription.CurrentPeriodEnd, 0),
			"next_due":             time.Unix(subscription.CurrentPeriodEnd, 0).Format("2006-01-02"),
			"status":               subscription.Status,
			"included_credits":     includedCredits,
			"used_credits":         usedCredits,
			"overage_charge":       overageCharge,
			"cancel_at_period_end": subscription.CancelAtPeriodEnd,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func fetchUsedCredits(stripeCustomerID string, eventName string, subscription *stripe.Subscription) (int, error) {
	err := setupStripeClient()
	if err != nil {
		return 0, fmt.Errorf("failed to set up Stripe client: %v", err)
	}

	meterID, err := fetchMeterIDByEventName(eventName)
	if err != nil {
		return 0, err
	}

	startTime := time.Unix(subscription.CurrentPeriodStart, 0)
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Unix(subscription.CurrentPeriodEnd, 0)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.UTC)

	fmt.Printf("Fetching usage records from %v to %v\n", startTime, endTime)

	params := &stripe.BillingMeterEventSummaryListParams{
		Customer:            stripe.String(stripeCustomerID),
		StartTime:           stripe.Int64(startTime.Unix()),
		EndTime:             stripe.Int64(endTime.Unix()),
		ValueGroupingWindow: stripe.String("day"), // Aggregate by day
		ID:                  stripe.String(meterID),
	}

	iterator := metereventsummary.List(params)

	usedCredits := 0
	for iterator.Next() {
		summary := iterator.BillingMeterEventSummary()
		if summary == nil {
			fmt.Println("Empty summary, skipping...")
			continue
		}
		fmt.Printf("Aggregated value on %v: %v\n", time.Unix(summary.StartTime, 0), summary.AggregatedValue)
		usedCredits += int(summary.AggregatedValue)
	}

	if err := iterator.Err(); err != nil {
		return 0, fmt.Errorf("error getting usage records: %v", err)
	}

	return usedCredits, nil
}

func fetchMeterIDByEventName(eventName string) (string, error) {
	params := &stripe.BillingMeterListParams{}
	iterator := meter.List(params)

	for iterator.Next() {
		m := iterator.BillingMeter()
		if m.EventName == eventName {
			return m.ID, nil
		}
	}

	if err := iterator.Err(); err != nil {
		return "", fmt.Errorf("error listing billing meters: %v", err)
	}

	return "", fmt.Errorf("meter with event name '%s' not found", eventName)
}

func StripeCheckSubscriptionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			utils.SendJSONError(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		if user.SubscriptionID == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"isSubscribed": false})
			return
		}

		err := setupStripeClient()
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error setting up Stripe client: %v", err), http.StatusInternalServerError)
			return
		}

		subscription, err := subscription.Get(*user.SubscriptionID, nil)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error getting subscription: %v", err), http.StatusInternalServerError)
			return
		}

		isSubscribed := subscription.Status == "active"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"isSubscribed": isSubscribed})
	}
}

func StripeCreateBillingPortalSessionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			utils.SendJSONError(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		if user.StripeUserID == "" {
			utils.SendJSONError(w, "User does not have a Stripe Customer ID", http.StatusBadRequest)
			return
		}

		err := setupStripeClient()
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error setting up Stripe client: %v", err), http.StatusInternalServerError)
			return
		}

		var requestBody struct {
			ReturnURL string `json:"returnURL"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			utils.SendJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		params := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(user.StripeUserID),
			ReturnURL: stripe.String(requestBody.ReturnURL),
		}

		session, err := billingportal.New(params)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating billing portal session: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
	}
}

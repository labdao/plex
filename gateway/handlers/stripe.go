package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
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

func createCheckoutSession(stripeUserID, walletAddress string) (*stripe.CheckoutSession, error) {
	err := setupStripeClient()
	if err != nil {
		return nil, err
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
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
			TrialPeriodDays: stripe.Int64(7),
			Metadata: map[string]string{
				"Wallet Address": walletAddress,
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		SuccessURL:         stripe.String(frontendURL + "/subscription-success"),
		CancelURL:          stripe.String(frontendURL + "/subscription-canceled"),
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

// func createCheckoutSession(stripeUserID, walletAddress, modelID, scatteringMethod, kwargs string) (*stripe.CheckoutSession, error) {
// 	err := setupStripeClient()
// 	if err != nil {
// 		return nil, err
// 	}

// 	frontendURL := os.Getenv("FRONTEND_URL")
// 	if frontendURL == "" {
// 		frontendURL = "http://localhost:3000"
// 	}

// 	productID := os.Getenv("STRIPE_PRODUCT_ID")
// 	if productID == "" {
// 		return nil, errors.New("STRIPE_PRODUCT_ID environment variable not set")
// 	}

// 	priceID := os.Getenv("STRIPE_PRICE_ID")
// 	if priceID == "" {
// 		return nil, errors.New("STRIPE_PRICE_ID environment variable not set")
// 	}

// 	priceObj, err := price.Get(priceID, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching price: %v", err)
// 	}

// 	params := &stripe.CheckoutSessionParams{
// 		Customer:   stripe.String(stripeUserID),
// 		SuccessURL: stripe.String(frontendURL),
// 		LineItems: []*stripe.CheckoutSessionLineItemParams{
// 			{
// 				Price:    stripe.String(priceObj.ID),
// 				Quantity: stripe.Int64(1),
// 				AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
// 					Enabled: stripe.Bool(false),
// 				},
// 			},
// 		},
// 		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
// 		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
// 			Metadata: map[string]string{
// 				"Wallet Address":    walletAddress,
// 				"Model ID":          modelID,
// 				"Scattering Method": scatteringMethod,
// 				"Kwargs":            kwargs,
// 			},
// 			SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOnSession)),
// 		},
// 		SavedPaymentMethodOptions: &stripe.CheckoutSessionSavedPaymentMethodOptionsParams{
// 			PaymentMethodSave: stripe.String(string(stripe.CheckoutSessionSavedPaymentMethodOptionsPaymentMethodSaveEnabled)),
// 		},
// 		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
// 	}

// 	params.AddMetadata("Stripe User ID", stripeUserID)

// 	session, err := session.New(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return session, nil
// }

func StripeCreateCheckoutSessionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			utils.SendJSONError(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		thresholdStr := os.Getenv("TIER_THRESHOLD")
		threshold, _ := strconv.Atoi(thresholdStr)

		if user.ComputeTally < threshold || user.SubscriptionStatus == "active" {
			utils.SendJSONError(w, "User does not need a subscription at this time", http.StatusBadRequest)
			return
		}

		session, err := createCheckoutSession(user.StripeUserID, user.WalletAddress)
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

			if subscription.Status == "trialing" || subscription.Status == "active" {
				user.SubscriptionStatus = "active"
			} else {
				user.SubscriptionStatus = string(subscription.Status)
			}
			user.SubscriptionID = subscription.ID
			result = db.Save(&user)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error updating user subscription status: %v\n", result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Printf("Subscription %s for user %s updated to %s\n", subscription.ID, walletAddress, user.SubscriptionStatus)

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

// func StripeFulfillmentHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		payload, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
// 			w.WriteHeader(http.StatusServiceUnavailable)
// 			return
// 		}

// 		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")

// 		event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), endpointSecret, webhook.ConstructEventOptions{
// 			IgnoreAPIVersionMismatch: true,
// 		})
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		switch event.Type {
// 		case "payment_intent.succeeded":
// 			var paymentIntent stripe.PaymentIntent
// 			err := json.Unmarshal(event.Data.Raw, &paymentIntent)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error parsing payment intent: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			walletAddress := paymentIntent.Metadata["Wallet Address"]
// 			modelID := paymentIntent.Metadata["Model ID"]
// 			scatteringMethod := paymentIntent.Metadata["Scattering Method"]
// 			kwargsRaw := paymentIntent.Metadata["Kwargs"]

// 			var user models.User
// 			result := db.Where("wallet_address ILIKE ?", walletAddress).First(&user)
// 			if result.Error != nil {
// 				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 					fmt.Fprintf(os.Stderr, "User with wallet address %s not found\n", walletAddress)
// 					w.WriteHeader(http.StatusNotFound)
// 				} else {
// 					fmt.Fprintf(os.Stderr, "Error querying user: %v\n", result.Error)
// 					w.WriteHeader(http.StatusInternalServerError)
// 				}
// 				return
// 			}

// 			var model models.Model
// 			result = db.Where("id = ?", modelID).First(&model)
// 			if result.Error != nil {
// 				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 					fmt.Fprintf(os.Stderr, "Model with ID %d not found\n", model.ID)
// 					w.WriteHeader(http.StatusNotFound)
// 				} else {
// 					fmt.Fprintf(os.Stderr, "Error querying model: %v\n", result.Error)
// 					w.WriteHeader(http.StatusInternalServerError)
// 				}
// 				return
// 			}

// 			var kwargs map[string][]interface{}
// 			err = json.Unmarshal([]byte(kwargsRaw), &kwargs)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error parsing kwargs JSON: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			ioList, err := ipwl.InitializeIo(model.S3URI, scatteringMethod, kwargs, db)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error initializing IO list: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			experiment := models.Experiment{
// 				WalletAddress: user.WalletAddress,
// 				Name:          "Experiment created by Stripe webhook",
// 				CreatedAt:     time.Now().UTC(),
// 				Public:        false,
// 			}

// 			result = db.Create(&experiment)
// 			if result.Error != nil {
// 				fmt.Fprintf(os.Stderr, "Error creating experiment: %v\n", result.Error)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			for _, ioItem := range ioList {
// 				inputsJSON, err := json.Marshal(ioItem.Inputs)
// 				if err != nil {
// 					fmt.Fprintf(os.Stderr, "Error marshaling job inputs: %v\n", err)
// 					w.WriteHeader(http.StatusInternalServerError)
// 					return
// 				}

// 				job := models.Job{
// 					ModelID:       model.ID,
// 					ExperimentID:  experiment.ID,
// 					WalletAddress: user.WalletAddress,
// 					Inputs:        datatypes.JSON(inputsJSON),
// 					CreatedAt:     time.Now().UTC(),
// 					Public:        false,
// 				}

// 				result = db.Create(&job)
// 				if result.Error != nil {
// 					fmt.Fprintf(os.Stderr, "Error creating job: %v\n", result.Error)
// 					w.WriteHeader(http.StatusInternalServerError)
// 					return
// 				}

// 				requestTracker := models.InferenceEvent{
// 					JobID:      job.ID,
// 					RetryCount: 0,
// 					EventTime:  time.Now().UTC(),
// 					EventType:  models.EventTypeJobCreated,
// 				}
// 				if err := db.Create(&requestTracker).Error; err != nil {
// 					fmt.Fprintf(os.Stderr, "Error creating request tracker: %v\n", err)
// 					w.WriteHeader(http.StatusInternalServerError)
// 					return
// 				}

// 				for _, input := range ioItem.Inputs {
// 					var idsToAdd []string
// 					switch v := input.(type) {
// 					case string:
// 						strInput, ok := input.(string)
// 						if !ok {
// 							continue
// 						}
// 						if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
// 							split := strings.SplitN(strInput, "/", 2)
// 							id := split[0]
// 							idsToAdd = append(idsToAdd, id)
// 						}
// 					case []interface{}:
// 						fmt.Println("found slice, checking each for 'Qm' prefix")
// 						for _, elem := range v {
// 							strInput, ok := elem.(string)
// 							if !ok {
// 								continue
// 							}
// 							if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
// 								split := strings.SplitN(strInput, "/", 2)
// 								id := split[0]
// 								idsToAdd = append(idsToAdd, id)
// 							}
// 						}
// 					default:
// 						continue
// 					}
// 					for _, id := range idsToAdd {
// 						var file models.File
// 						result := db.First(&file, "id = ?", id)
// 						if result.Error != nil {
// 							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 								fmt.Fprintf(os.Stderr, "File with ID %v not found\n", id)
// 								w.WriteHeader(http.StatusNotFound)
// 								return
// 							} else {
// 								fmt.Fprintf(os.Stderr, "Error looking up File: %v\n", result.Error)
// 								w.WriteHeader(http.StatusInternalServerError)
// 								return
// 							}
// 						}
// 						job.InputFiles = append(job.InputFiles, file)
// 					}
// 				}
// 				result = db.Save(&job)
// 				if result.Error != nil {
// 					fmt.Fprintf(os.Stderr, "Error updating job with input data: %v\n", result.Error)
// 					w.WriteHeader(http.StatusInternalServerError)
// 					return
// 				}
// 			}

// 			user.ComputeTally += model.ComputeCost
// 			result = db.Save(user)
// 			if result.Error != nil {
// 				fmt.Fprintf(os.Stderr, "Error updating user compute tally: %v\n", result.Error)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			thresholdStr := os.Getenv("TIER_THRESHOLD")
// 			if thresholdStr == "" {
// 				fmt.Fprintf(os.Stderr, "TIER_THRESHOLD environment variable is not set\n")
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			threshold, err := strconv.Atoi(thresholdStr)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error converting TIER_THRESHOLD to integer: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			err = UpdateUserTier(db, user.WalletAddress, threshold)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error updating user tier: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			fmt.Printf("PaymentIntent succeeded, Amount: %v, WalletAddress: %v\n", paymentIntent.Amount, walletAddress)

// 			w.WriteHeader(http.StatusOK)

// 		case "customer.subscription.created", "customer.subscription.updated":
// 			var subscription stripe.Subscription
// 			err := json.Unmarshal(event.Data.Raw, &subscription)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Error parsing subscription: %v\n", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			walletAddress, ok := subscription.Metadata["Wallet Address"]
// 			if !ok {
// 				fmt.Fprintf(os.Stderr, "Wallet Address not found in subscription metadata\n")
// 				w.WriteHeader(http.StatusBadRequest)
// 				return
// 			}

// 			var user models.User
// 			result := db.Where("wallet_address ILIKE ?", walletAddress).First(&user)
// 			if result.Error != nil {
// 				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 					fmt.Fprintf(os.Stderr, "User with wallet address %s not found\n", walletAddress)
// 					w.WriteHeader(http.StatusNotFound)
// 				} else {
// 					fmt.Fprintf(os.Stderr, "Error querying user: %v\n", result.Error)
// 					w.WriteHeader(http.StatusInternalServerError)
// 				}
// 				return
// 			}

// 			user.SubscriptionStatus = string(subscription.Status)
// 			user.SubscriptionID = subscription.ID
// 			result = db.Save(&user)
// 			if result.Error != nil {
// 				fmt.Fprintf(os.Stderr, "Error updating user subscription status: %v\n", result.Error)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			fmt.Printf("Subscription %s for user %s updated to %s\n", subscription.ID, walletAddress, subscription.Status)
// 			w.WriteHeader(http.StatusOK)

// 		default:
// 			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
// 		}
// 	}
// }

func RecordUsage(stripeCustomerID string, usage int64) error {
	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String("usage_recorded"),
		Payload: map[string]string{
			"value":              strconv.FormatInt(usage, 10),
			"stripe_customer_id": stripeCustomerID,
		},
		Identifier: stripe.String(fmt.Sprintf("usage-%d", time.Now().Unix())),
	}
	_, err := meterevent.New(params)
	return err
}

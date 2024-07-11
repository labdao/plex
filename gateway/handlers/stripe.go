package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentmethod"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/webhook"
	"gorm.io/datatypes"
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

func createCheckoutSession(stripeUserID, walletAddress, modelID, scatteringMethod, kwargs string) (*stripe.CheckoutSession, error) {
	err := setupStripeClient()
	if err != nil {
		return nil, err
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	productID := os.Getenv("STRIPE_PRODUCT_ID")
	if productID == "" {
		return nil, errors.New("STRIPE_PRODUCT_ID environment variable not set")
	}

	priceID := os.Getenv("STRIPE_PRICE_ID")
	if priceID == "" {
		return nil, errors.New("STRIPE_PRICE_ID environment variable not set")
	}

	priceObj, err := price.Get(priceID, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching price: %v", err)
	}

	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(stripeUserID),
		SuccessURL: stripe.String(frontendURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceObj.ID),
				Quantity: stripe.Int64(1),
				AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
					Enabled: stripe.Bool(false),
				},
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"Wallet Address":    walletAddress,
				"Model ID":          modelID,
				"Scattering Method": scatteringMethod,
				"Kwargs":            kwargs,
			},
			SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOnSession)),
		},
		SavedPaymentMethodOptions: &stripe.CheckoutSessionSavedPaymentMethodOptionsParams{
			PaymentMethodSave: stripe.String(string(stripe.CheckoutSessionSavedPaymentMethodOptionsPaymentMethodSaveEnabled)),
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

		var requestData struct {
			ModelID          string `json:"model"`
			ScatteringMethod string `json:"scatteringMethod"`
			Kwargs           string `json:"kwargs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
			return
		}

		// TODO: pass in model

		// modelID := r.URL.Query().Get("model")
		// if modelID == "" {
		// 	utils.SendJSONError(w, "Model ID not provided", http.StatusBadRequest)
		// 	return
		// }

		// var model models.Model
		// result := db.Where("id = ?", modelID).First(&model)
		// if result.Error != nil {
		// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 		utils.SendJSONError(w, "Model not found", http.StatusNotFound)
		// 	} else {
		// 		utils.SendJSONError(w, fmt.Sprintf("Error fetching Model: %v", result.Error), http.StatusInternalServerError)
		// 	}
		// 	return
		// }

		// TODO: modify so we're not passing in hardcoded value
		session, err := createCheckoutSession(user.StripeUserID, user.WalletAddress, requestData.ModelID, requestData.ScatteringMethod, requestData.Kwargs)
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

		fmt.Fprintf(os.Stderr, "Raw webhook payload: %s\n", string(payload))
		fmt.Fprintf(os.Stderr, "Stripe-Signature header: %s\n", r.Header.Get("Stripe-Signature"))
		fmt.Fprintf(os.Stderr, "Endpoint Secret: %s\n", endpointSecret)

		event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), endpointSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
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

			walletAddress := paymentIntent.Metadata["Wallet Address"]
			modelID := paymentIntent.Metadata["Model ID"]
			scatteringMethod := paymentIntent.Metadata["Scattering Method"]
			kwargsRaw := paymentIntent.Metadata["Kwargs"]

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

			var model models.Model
			result = db.Where("id = ?", modelID).First(&model)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					fmt.Fprintf(os.Stderr, "Model with ID %d not found\n", model.ID)
					w.WriteHeader(http.StatusNotFound)
				} else {
					fmt.Fprintf(os.Stderr, "Error querying model: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			var kwargs map[string][]interface{}
			err = json.Unmarshal([]byte(kwargsRaw), &kwargs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing kwargs JSON: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ioList, err := ipwl.InitializeIo(model.S3URI, scatteringMethod, kwargs, db)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing IO list: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			experiment := models.Experiment{
				WalletAddress: user.WalletAddress,
				Name:          "Experiment created by Stripe webhook",
				CreatedAt:     time.Now().UTC(),
				Public:        false,
			}

			result = db.Create(&experiment)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error creating experiment: %v\n", result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, ioItem := range ioList {
				inputsJSON, err := json.Marshal(ioItem.Inputs)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling job inputs: %v\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				job := models.Job{
					ModelID:       model.ID,
					ExperimentID:  experiment.ID,
					WalletAddress: user.WalletAddress,
					Inputs:        datatypes.JSON(inputsJSON),
					CreatedAt:     time.Now().UTC(),
					Public:        false,
				}

				result = db.Create(&job)
				if result.Error != nil {
					fmt.Fprintf(os.Stderr, "Error creating job: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				requestTracker := models.InferenceEvent{
					JobID:      job.ID,
					RetryCount: 0,
					EventTime:  time.Now().UTC(),
					EventType:  models.EventTypeJobCreated,
				}
				if err := db.Create(&requestTracker).Error; err != nil {
					fmt.Fprintf(os.Stderr, "Error creating request tracker: %v\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				for _, input := range ioItem.Inputs {
					var idsToAdd []string
					switch v := input.(type) {
					case string:
						strInput, ok := input.(string)
						if !ok {
							continue
						}
						if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
							split := strings.SplitN(strInput, "/", 2)
							id := split[0]
							idsToAdd = append(idsToAdd, id)
						}
					case []interface{}:
						fmt.Println("found slice, checking each for 'Qm' prefix")
						for _, elem := range v {
							strInput, ok := elem.(string)
							if !ok {
								continue
							}
							if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
								split := strings.SplitN(strInput, "/", 2)
								id := split[0]
								idsToAdd = append(idsToAdd, id)
							}
						}
					default:
						continue
					}
					for _, id := range idsToAdd {
						var file models.File
						result := db.First(&file, "id = ?", id)
						if result.Error != nil {
							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
								fmt.Fprintf(os.Stderr, "File with ID %v not found\n", id)
								w.WriteHeader(http.StatusNotFound)
								return
							} else {
								fmt.Fprintf(os.Stderr, "Error looking up File: %v\n", result.Error)
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
						job.InputFiles = append(job.InputFiles, file)
					}
				}
				result = db.Save(&job)
				if result.Error != nil {
					fmt.Fprintf(os.Stderr, "Error updating job with input data: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			user.ComputeTally += model.ComputeCost
			result = db.Save(user)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error updating user compute tally: %v\n", result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			thresholdStr := os.Getenv("TIER_THRESHOLD")
			if thresholdStr == "" {
				fmt.Fprintf(os.Stderr, "TIER_THRESHOLD environment variable is not set\n")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			threshold, err := strconv.Atoi(thresholdStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting TIER_THRESHOLD to integer: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = UpdateUserTier(db, user.WalletAddress, threshold)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating user tier: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Printf("PaymentIntent succeeded, Amount: %v, WalletAddress: %v\n", paymentIntent.Amount, walletAddress)

			w.WriteHeader(http.StatusOK)

		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}
	}
}

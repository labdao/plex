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

	"github.com/google/uuid"
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

func createCheckoutSession(stripeUserID string, computeCost int, walletAddress, toolCID, scatteringMethod, kwargs string) (*stripe.CheckoutSession, error) {
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
		// SuccessURL: stripe.String(frontendURL),
		SuccessURL: stripe.String("http://localhost:3000"),
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
				"Stripe User ID":    stripeUserID,
				"Wallet Address":    walletAddress,
				"Tool CID":          toolCID,
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
			ToolCID          string `json:"toolCid"`
			ScatteringMethod string `json:"scatteringMethod"`
			Kwargs           string `json:"kwargs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
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
		session, err := createCheckoutSession(user.StripeUserID, 10, user.WalletAddress, requestData.ToolCID, requestData.ScatteringMethod, requestData.Kwargs)
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
			toolCID := paymentIntent.Metadata["Tool CID"]
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

			var tool models.Tool
			result = db.Where("cid = ?", toolCID).First(&tool)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					fmt.Fprintf(os.Stderr, "Tool with CID %s not found\n", toolCID)
					w.WriteHeader(http.StatusNotFound)
				} else {
					fmt.Fprintf(os.Stderr, "Error querying tool: %v\n", result.Error)
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

			ioList, err := ipwl.InitializeIo(toolCID, scatteringMethod, kwargs, db)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing IO list: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			flowUUID := uuid.New().String()

			flow := models.Flow{
				WalletAddress: user.WalletAddress,
				Name:          "Flow created by Stripe webhook",
				StartTime:     time.Now(),
				FlowUUID:      flowUUID,
				Public:        false,
			}

			result = db.Create(&flow)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error creating flow: %v\n", result.Error)
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

				var queue models.QueueType
				if tool.ToolType == "ray" {
					queue = models.QueueTypeRay
				}

				var jobType models.JobType
				if tool.ToolType == "ray" {
					jobType = models.JobTypeRay
				} else {
					jobType = models.JobTypeBacalhau
				}

				job := models.Job{
					ToolID:        ioItem.Tool.S3,
					FlowID:        flow.ID,
					WalletAddress: user.WalletAddress,
					Inputs:        datatypes.JSON(inputsJSON),
					Queue:         queue,
					CreatedAt:     time.Now(),
					Public:        false,
					JobType:       jobType,
				}

				result = db.Create(&job)
				if result.Error != nil {
					fmt.Fprintf(os.Stderr, "Error creating job: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				for _, input := range ioItem.Inputs {
					var cidsToAdd []string
					switch v := input.(type) {
					case string:
						strInput, ok := input.(string)
						if !ok {
							continue
						}
						if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
							split := strings.SplitN(strInput, "/", 2)
							cid := split[0]
							cidsToAdd = append(cidsToAdd, cid)
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
								cid := split[0]
								cidsToAdd = append(cidsToAdd, cid)
							}
						}
					default:
						continue
					}
					for _, cid := range cidsToAdd {
						var dataFile models.DataFile
						result := db.First(&dataFile, "cid = ?", cid)
						if result.Error != nil {
							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
								fmt.Fprintf(os.Stderr, "DataFile with CID %v not found\n", cid)
								w.WriteHeader(http.StatusNotFound)
								return
							} else {
								fmt.Fprintf(os.Stderr, "Error looking up DataFile: %v\n", result.Error)
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
						job.InputFiles = append(job.InputFiles, dataFile)
					}
				}
				result = db.Save(&job)
				if result.Error != nil {
					fmt.Fprintf(os.Stderr, "Error updating job with input data: %v\n", result.Error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			user.ComputeTally += tool.ComputeCost
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

// func StripeFullfillmentHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		payload, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
// 			w.WriteHeader(http.StatusServiceUnavailable)
// 			return
// 		}

// 		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")

// 		event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
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

// 			walletAddress, ok := paymentIntent.Metadata["walletAddress"]
// 			if !ok {
// 				fmt.Fprintf(os.Stderr, "Wallet address not found in payment intent metadata\n")
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

// 			// Convert Stripe amount (in cents) to a float64 representation
// 			amount := float64(paymentIntent.Amount) / 100.0

// 			transaction := models.Transaction{
// 				ID:          uuid.New().String(),
// 				Amount:      amount,
// 				IsDebit:     false, // Assuming payment intents are always credits (money in)
// 				UserID:      user.WalletAddress,
// 				Description: "Stripe payment intent succeeded",
// 			}

// 			if result := db.Create(&transaction); result.Error != nil {
// 				fmt.Fprintf(os.Stderr, "Failed to save transaction: %v\n", result.Error)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			// Create the flow in your database using the extracted information
// 			flow := models.Flow{
// 				WalletAddress:    walletAddress,
// 				ToolCID:          toolCID,
// 				ScatteringMethod: scatteringMethod,
// 				Kwargs:           kwargs,
// 				// Set other necessary fields
// 			}

// 			if result := db.Create(&flow); result.Error != nil {
// 				fmt.Fprintf(os.Stderr, "Failed to save flow: %v\n", result.Error)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}

// 			fmt.Printf("PaymentIntent succeeded, Amount: %v, WalletAddress: %v\n", paymentIntent.Amount, walletAddress)

// 		default:
// 			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}
// }

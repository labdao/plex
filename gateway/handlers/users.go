package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/web3"

	"gorm.io/gorm"
)

func AddUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
			fmt.Println("Received non-POST request for /user endpoint.")
			return
		}

		var requestData struct {
			WalletAddress string `json:"walletAddress"`
		}

		if err := utils.ReadRequestBody(r, &requestData); err != nil {
			utils.SendJSONError(w, "Error parsing request body", http.StatusBadRequest)
			fmt.Println("Error decoding request body:", err)
			return
		}

		fmt.Printf("Received request to create/get user: WalletAddress: %s\n", requestData.WalletAddress)

		isValidAddress := web3.IsValidEthereumAddress(requestData.WalletAddress)
		if !isValidAddress {
			utils.SendJSONError(w, "Invalid wallet address", http.StatusBadRequest)
			fmt.Println("Invalid wallet address:", requestData.WalletAddress)
			return
		}

		did, err := middleware.GetUserDIDFromRequest(r, db)
		if err != nil {
			utils.SendJSONError(w, "Error getting user DID from request", http.StatusInternalServerError)
			fmt.Println("Error getting user DID from request:", err)
			return
		}

		var user models.User
		err = db.Where("wallet_address = ?", requestData.WalletAddress).First(&user).Error

		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create a new one
			stripeUserID, err := createStripeCustomer(requestData.WalletAddress)
			if err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error creating Stripe customer: %v", err), http.StatusInternalServerError)
				fmt.Println("Error creating Stripe customer:", err)
				return
			}

			newUser := models.User{
				WalletAddress:  requestData.WalletAddress,
				DID:            did,
				CreatedAt:      time.Now().UTC(),
				OrganizationID: 1,
				StripeUserID:   stripeUserID,
				SubscriptionID: nil,
			}
			if result := db.Create(&newUser); result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
				fmt.Println("Error creating user in database:", result.Error)
				return
			}
			fmt.Printf("Successfully created user with WalletAddress: %s and StripeUserID: %s\n", newUser.WalletAddress, newUser.StripeUserID)
			user = newUser // Set user to newUser for consistent response
		} else if err != nil {
			utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error:", err)
			return
		} else {
			// User already exists, update DID if necessary
			if user.DID == "" {
				user.DID = did
				if err := db.Save(&user).Error; err != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
					fmt.Println("Error updating user in database:", err)
					return
				}
				fmt.Printf("Successfully updated user with WalletAddress: %s with new DID\n", user.WalletAddress)
			} else {
				fmt.Printf("User already exists with WalletAddress: %s and has a DID\n", user.WalletAddress)
			}
		}

		// Prepare response
		response := struct {
			WalletAddress      string      `json:"walletAddress"`
			DID                string      `json:"did"`
			IsAdmin            bool        `json:"isAdmin"`
			Tier               models.Tier `json:"tier"`
			SubscriptionStatus string      `json:"subscriptionStatus"`
		}{
			WalletAddress:      user.WalletAddress,
			DID:                user.DID,
			IsAdmin:            user.Admin,
			Tier:               user.Tier,
			SubscriptionStatus: user.SubscriptionStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(response)
	}
}

func GetUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusInternalServerError)
			return
		}

		response := struct {
			WalletAddress      string      `json:"walletAddress"`
			DID                string      `json:"did"`
			IsAdmin            bool        `json:"isAdmin"`
			Tier               models.Tier `json:"tier"`
			SubscriptionStatus string      `json:"subscriptionStatus"`
		}{
			WalletAddress:      user.WalletAddress,
			DID:                user.DID,
			IsAdmin:            user.Admin,
			Tier:               user.Tier,
			SubscriptionStatus: user.SubscriptionStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateUserTier(db *gorm.DB, walletAddress string, threshold int) error {
	var user models.User
	err := db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		return err
	}

	if user.ComputeTally >= threshold && user.Tier != models.TierPaid {
		user.Tier = models.TierPaid

		if err := db.Save(&user).Error; err != nil {
			fmt.Printf("Error updating tier for user with WalletAddress: %s: %v\n", walletAddress, err)
			return err
		}
		fmt.Printf("Successfully updated tier for user with WalletAddress: %s\n", walletAddress)
	} else {
		fmt.Printf("No need to update tier for user with WalletAddress: %s\n", walletAddress)
	}

	return nil
}

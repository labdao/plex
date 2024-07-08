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

		fmt.Printf("Received request to create user: WalletAddress: %s\n", requestData.WalletAddress)

		isValidAddress := web3.IsValidEthereumAddress(requestData.WalletAddress)
		if !isValidAddress {
			utils.SendJSONError(w, "Invalid wallet address", http.StatusBadRequest)
			fmt.Println("Invalid wallet address:", requestData.WalletAddress)
			return
		}

		var existingUser models.User
		err := db.Where("wallet_address = ?", requestData.WalletAddress).First(&existingUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error:", err)
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
			newUser := models.User{
				WalletAddress:  requestData.WalletAddress,
				DID:            did,
				CreatedAt:      time.Now().UTC(),
				OrganizationID: 1,
			}
			if result := db.Create(&newUser); result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
				fmt.Println("Error creating user in database:", result.Error)
				return
			}
			fmt.Printf("Successfully created user with WalletAddress: %s\n", newUser.WalletAddress)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newUser)
		} else if err != nil {
			utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error:", err)
			return
		} else {
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		}
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
			WalletAddress string      `json:"walletAddress"`
			DID           string      `json:"did"`
			IsAdmin       bool        `json:"isAdmin"`
			Tier          models.Tier `json:"tier"`
		}{
			WalletAddress: user.WalletAddress,
			DID:           user.DID,
			IsAdmin:       user.Admin,
			Tier:          user.Tier,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getTier(db *gorm.DB, walletAddress string) (models.Tier, error) {
	var user models.User
	err := db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		return models.TierFree, err
	}
	return user.Tier, nil
}

func UpdateUserTier(db *gorm.DB, walletAddress string, threshold int) error {
	var user models.User
	err := db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		return err
	}

	if user.ComputeTally >= threshold && user.Tier != models.TierPaid {
		user.Tier = models.TierPaid

		if user.StripeUserID == "" {
			stripeUserID, err := createStripeCustomer(walletAddress)
			if err != nil {
				fmt.Printf("Error creating Stripe customer for user with WalletAddress: %s: %v\n", walletAddress, err)
				return err
			}
			user.StripeUserID = stripeUserID
		}

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

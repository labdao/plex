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

		// Extract DID from the JWT token in the request
		did, err := middleware.GetUserDIDFromRequest(r, db)
		if err != nil {
			utils.SendJSONError(w, "Error getting user DID from request", http.StatusInternalServerError)
			fmt.Println("Error getting user DID from request:", err)
			return
		}

		// Check if a user exists with the given wallet address
		var user models.User
		err = db.Where("wallet_address = ?", requestData.WalletAddress).First(&user).Error

		if err == gorm.ErrRecordNotFound {
			// No user found by wallet address, create a new user
			newUser := models.User{
				WalletAddress: requestData.WalletAddress,
				DID:           did,
				CreatedAt:     time.Now(),
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
			// Database error other than record not found
			utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error:", err)
			return
		} else {
			// User found by wallet address, update with DID if necessary
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

// func AddUserHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
// 			utils.SendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
// 			fmt.Println("Received non-POST request for /user endpoint.")
// 			return
// 		}

// 		var requestData struct {
// 			WalletAddress string `json:"walletAddress"`
// 		}

// 		if err := utils.ReadRequestBody(r, &requestData); err != nil {
// 			utils.SendJSONError(w, "Error parsing request body", http.StatusBadRequest)
// 			fmt.Println("Error decoding request body:", err)
// 			return
// 		}

// 		fmt.Printf("Received request to create user: WalletAddress: %s\n", requestData.WalletAddress)

// 		isValidAddress := web3.IsValidEthereumAddress(requestData.WalletAddress)
// 		if !isValidAddress {
// 			utils.SendJSONError(w, "Invalid wallet address", http.StatusBadRequest)
// 			fmt.Println("Invalid wallet address:", requestData.WalletAddress)
// 			return
// 		}

// 		var existingUser models.User
// 		if err := db.Where("wallet_address = ?", requestData.WalletAddress).First(&existingUser).Error; err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				did, err := middleware.GetUserDIDFromRequest(r, db)
// 				if err != nil {
// 					utils.SendJSONError(w, "Error fetching user record with the provided DID", http.StatusInternalServerError)
// 					fmt.Printf("Error fetching user record with the provided DID: %v\n", err)
// 					return
// 				}

// 				newUser := models.User{
// 					WalletAddress: requestData.WalletAddress,
// 					DID:           did,
// 					CreatedAt:     time.Now(),
// 				}
// 				if result := db.Create(&newUser); result.Error != nil {
// 					utils.SendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
// 					fmt.Println("Error creating user in database:", result.Error)
// 					return
// 				}
// 				fmt.Printf("Successfully created user with WalletAddress: %s\n", newUser.WalletAddress)
// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusCreated)
// 				json.NewEncoder(w).Encode(newUser)
// 			} else {
// 				utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
// 				fmt.Println("Database error:", err)
// 			}
// 		} else {
// 			if existingUser.DID == "" {
// 				did, err := middleware.GetUserDIDFromRequest(r, db)
// 				if err != nil {
// 					utils.SendJSONError(w, "Error getting user DID from request", http.StatusInternalServerError)
// 					fmt.Println("Error getting user DID from request:", err)
// 					return
// 				}
// 				existingUser.DID = did
// 				if err := db.Save(&existingUser).Error; err != nil {
// 					utils.SendJSONError(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
// 					fmt.Println("Error updating user in database:", err)
// 					return
// 				}
// 			}
// 			fmt.Printf("User already exists with WalletAddress: %s\n", existingUser.WalletAddress)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode(existingUser)
// 		}
// 	}
// }

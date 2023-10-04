package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
			EmailAddress  string `json:"emailAddress"`
		}

		if err := utils.ReadRequestBody(r, &requestData); err != nil {
			utils.SendJSONError(w, "Error parsing request body", http.StatusBadRequest)
			fmt.Println("Error decoding request body:", err)
			return
		}

		fmt.Printf("Received request to create user: WalletAddress: %s, Email: %s,\n", requestData.WalletAddress, requestData.EmailAddress)

		isValidAddress := web3.IsValidEthereumAddress(requestData.WalletAddress)
		if !isValidAddress {
			utils.SendJSONError(w, "Invalid wallet address", http.StatusBadRequest)
			fmt.Println("Invalid wallet address:", requestData.WalletAddress)
			return
		}

		var existingUser models.User
		if err := db.Where("wallet_address = ? AND email_address = ?", requestData.WalletAddress, requestData.EmailAddress).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// User does not exist, create new user
				newUser := models.User{
					WalletAddress: requestData.WalletAddress,
					EmailAddress:  requestData.EmailAddress,
				}
				if result := db.Create(&newUser); result.Error != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
					fmt.Println("Error creating user in database:", result.Error)
					return
				}
				fmt.Printf("Successfully created user with WalletAddress: %s, Email: %s\n", newUser.WalletAddress, newUser.EmailAddress)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(newUser)
			} else {
				// Some other error occurred during the query
				utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error:", err)
			}
		} else {
			// User with given wallet address and email already exists, return that user
			fmt.Printf("User already exists with WalletAddress: %s, Email: %s\n", existingUser.WalletAddress, existingUser.EmailAddress)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingUser)
		}
	}
}

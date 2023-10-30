package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/web3"

	"gorm.io/gorm"
)

func updateUserMemberStatus(db *gorm.DB, walletAddress string) (models.User, error) {
	var user models.User
	if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		return user, err
	}

	isMember, err := utils.CheckNFTOwnership(walletAddress)
	if err != nil {
		return user, err
	}

	if user.IsMember != isMember {
		user.IsMember = isMember
		if err := db.Save(&user).Error; err != nil {
			return user, err
		}
	}

	return user, nil
}

func AddUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
			fmt.Println("Received non-POST request for /user endpoint.")
			return
		}

		var requestData struct {
			WalletAddress string `json:"walletAddress"`
			IsMember      bool   `json:"isMember"`
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
		if err := db.Where("wallet_address = ?", requestData.WalletAddress).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// User does not exist, create new user
				newUser := models.User{
					WalletAddress: requestData.WalletAddress,
				}
				if result := db.Create(&newUser); result.Error != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
					fmt.Println("Error creating user in database:", result.Error)
					return
				}
				// Update membership status for the new user
				updatedUser, err := updateUserMemberStatus(db, newUser.WalletAddress)
				if err != nil {
					utils.SendJSONError(w, "Error updating user membership status", http.StatusInternalServerError)
					fmt.Println("Error updating user membership status:", err)
					return
				}
				fmt.Printf("Successfully created user with WalletAddress: %s\n", updatedUser.WalletAddress)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(updatedUser)
			} else {
				// Some other error occurred during the query
				utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error:", err)
			}
		} else {
			// User already exists, update membership status
			updatedUser, err := updateUserMemberStatus(db, existingUser.WalletAddress)
			if err != nil {
				utils.SendJSONError(w, "Error updating user membership status", http.StatusInternalServerError)
				fmt.Println("Error updating user membership status:", err)
				return
			}
			fmt.Printf("User already exists with WalletAddress: %s\n", updatedUser.WalletAddress)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedUser)
		}
	}
}

func CheckUserMemberStatusHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletAddress := vars["walletAddress"]

		var user models.User
		if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"isMember": user.IsMember})
	}
}

func UpdateUserMemberStatusHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletAddress := vars["walletAddress"]

		var user models.User
		if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		isMember, err := utils.CheckNFTOwnership(walletAddress)
		if err != nil {
			http.Error(w, "Error checking NFT ownership", http.StatusInternalServerError)
			return
		}

		if user.IsMember != isMember {
			user.IsMember = isMember
			if err := db.Save(&user).Error; err != nil {
				http.Error(w, "Error updating user", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"isMember": user.IsMember})
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"gorm.io/gorm"
)

func ListTransactionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			http.Error(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		var transactions []models.Transaction
		result := db.Preload("User").Where("wallet_address = ?", user.WalletAddress).Find(&transactions)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching Transactions: %v", result.Error), http.StatusInternalServerError)
			return
		}

		log.Println("Fetched Transactions from DB: ", transactions)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			http.Error(w, "Error encoding Transactions to JSON", http.StatusInternalServerError)
			return
		}
	}
}

type TransactionSummary struct {
	Tokens  float64 `json:"tokens"`
	Balance float64 `json:"balance"`
}

func SummaryTransactionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		ctxUser := r.Context().Value(middleware.UserContextKey)
		user, ok := ctxUser.(*models.User)
		if !ok {
			http.Error(w, "Unauthorized, user context not passed through auth middleware", http.StatusUnauthorized)
			return
		}

		var totalDebits, totalCredits float64

		db.Model(&models.Transaction{}).
			Where("wallet_address = ? AND is_debit = ?", user.WalletAddress, true).
			Select("sum(amount)").
			Row().
			Scan(&totalDebits)

		db.Model(&models.Transaction{}).
			Where("wallet_address = ? AND is_debit = ?", user.WalletAddress, false).
			Select("sum(amount)").
			Row().
			Scan(&totalCredits)

		summary := TransactionSummary{
			Tokens:  totalDebits,
			Balance: totalCredits - totalDebits,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(summary); err != nil {
			http.Error(w, "Error encoding summary to JSON", http.StatusInternalServerError)
			return
		}
	}
}

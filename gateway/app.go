package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors" // <-- Add this import for CORS

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataFile struct {
	ID            uint   `gorm:"primaryKey"`
	CID           string `gorm:"type:varchar(255);not null"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
	Filename      string `gorm:"type:varchar(255);not null"`
}

func main() {
	// Setup database connection
	dsn := "gorm.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&DataFile{})

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from your React frontend
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	// Health check endpoint
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	http.HandleFunc("/create-datafile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusBadRequest)
			return
		}

		var requestData struct {
			CID           string `json:"cid"`
			WalletAddress string `json:"wallet_address"`
			Filename      string `json:"filename"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		newFile := DataFile{
			CID:           requestData.CID,
			WalletAddress: requestData.WalletAddress,
			Filename:      requestData.Filename,
		}
		if result := db.Create(&newFile); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error creating datafile: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "DataFile created successfully!")
	})

	// Start the server with CORS middleware
	http.ListenAndServe(":8080", corsMiddleware.Handler(http.DefaultServeMux))
}

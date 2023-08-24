package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/rs/cors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataFile struct {
	ID            uint   `gorm:"primaryKey"`
	CID           string `gorm:"type:varchar(255);not null"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
	Filename      string `gorm:"type:varchar(255);not null"`
}

type User struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Username      string `gorm:"type:varchar(255);unique;not null" json:"username"`
	WalletAddress string `gorm:"type:varchar(42);not null" json:"walletAddress"`
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func main() {
	// Setup database connection
	dsn := "gorm.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&DataFile{}, &User{})

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

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
			fmt.Println("Received non-POST request for /user endpoint.")
			return
		}

		var requestData struct {
			Username      string `json:"username"`
			WalletAddress string `json:"walletAddress"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			sendJSONError(w, "Error parsing request body", http.StatusBadRequest)
			fmt.Println("Error decoding request body:", err)
			return
		}

		fmt.Printf("Received request to create user: Username: %s, WalletAddress: %s\n", requestData.Username, requestData.WalletAddress)

		// Regular expression to match Ethereum addresses
		re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
		isValdAddress := re.MatchString(requestData.WalletAddress)
		if !isValdAddress {
			fmt.Println("%s is not a valid Ethereum address", requestData.WalletAddress)
			sendJSONError(w, "Invalid Ethereum address", http.StatusBadRequest)
			return
		}

		newUser := User{
			Username:      requestData.Username,
			WalletAddress: requestData.WalletAddress,
		}
		if result := db.Create(&newUser); result.Error != nil {
			sendJSONError(w, fmt.Sprintf("Error creating user: %v", result.Error), http.StatusInternalServerError)
			fmt.Println("Error creating user in database:", result.Error)
			return
		}

		fmt.Printf("Successfully created user with ID: %d, Username: %s\n", newUser.ID, newUser.Username)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
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
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(http.DefaultServeMux))
}

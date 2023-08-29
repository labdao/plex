package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/rs/cors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataFile struct {
	ID            uint      `gorm:"primaryKey"`
	CID           string    `gorm:"column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	Timestamp     time.Time `gorm:""`
	Public        bool      `gorm:"default:true"`
	Visible       bool      `gorm:"default:true"`
}

type User struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Username      string `gorm:"type:varchar(255);unique;not null" json:"username"`
	WalletAddress string `gorm:"type:varchar(42);not null" json:"walletAddress"`
}

// TODO: Look into directly importing this from existing plex tool.go

type ToolInput struct {
	Type    string   `json:"type"`
	Glob    []string `json:"glob"`
	Default string   `json:"default"`
}

type ToolOutput struct {
	Type string   `json:"type"`
	Item string   `json:"item"`
	Glob []string `json:"glob"`
}

type Tool struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Author      string                `json:"author"`
	BaseCommand []string              `json:"baseCommand"`
	Arguments   []string              `json:"arguments"`
	DockerPull  string                `json:"dockerPull"`
	GpuBool     bool                  `json:"gpuBool"`
	MemoryGB    *int                  `json:"memoryGB"`
	NetworkBool bool                  `json:"networkBool"`
	Inputs      map[string]ToolInput  `json:"inputs"`
	Outputs     map[string]ToolOutput `json:"outputs"`
}

type ToolEntity struct {
	CID      string `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	ToolJSON string `gorm:"type:varchar(255);not null"`
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	}
}

func createUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
		isValidAddress := re.MatchString(requestData.WalletAddress)
		if !isValidAddress {
			fmt.Printf("%s is not a valid Ethereum address\n", requestData.WalletAddress)
			sendJSONError(w, "Invalid Ethereum address", http.StatusBadRequest)
			return
		}

		var existingUser User
		if err := db.Where("username = ? AND wallet_address = ?", requestData.Username, requestData.WalletAddress).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// User does not exist, create new user
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
			} else {
				// Some other error occurred during the query
				sendJSONError(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error:", err)
			}
		} else {
			// User with given username and wallet address already exists, return that user
			fmt.Printf("User already exists with ID: %d, Username: %s\n", existingUser.ID, existingUser.Username)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingUser)
		}
	}
}

func createDataFileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request at /create-datafile")

		if r.Method != http.MethodPost {
			sendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
			fmt.Println("Received non-POST request for /create-datafile endpoint.")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			sendJSONError(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}
		fmt.Println("Parsed multipart form")

		file, _, err := r.FormFile("file")
		if err != nil {
			sendJSONError(w, "Error retrieving file from multipart form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		walletAddress := r.FormValue("walletAddress")
		filename := r.FormValue("filename")
		publicBool := r.FormValue("public")
		visibleBool := r.FormValue("visible")

		fmt.Printf("Received file upload request for file: %s, walletAddress: %s, public: %s, visible: %s\n", filename, walletAddress, publicBool, visibleBool)

		tempFile, err := os.Create(filename)
		if err != nil {
			sendJSONError(w, "Error creating temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(filename)

		_, err = io.Copy(tempFile, file)
		if err != nil {
			sendJSONError(w, "Error writing temp file", http.StatusInternalServerError)
			return
		}
		tempFile.Close()

		cid, err := ipfs.WrapAndPinFile(tempFile.Name())
		if err != nil {
			sendJSONError(w, "Error pinning file to IPFS", http.StatusInternalServerError)
			return
		}

		dataFile := DataFile{
			CID:           cid,
			WalletAddress: walletAddress,
			Filename:      filename,
			Timestamp:     time.Now(),
		}

		if result := db.Create(&dataFile); result.Error != nil {
			sendJSONError(w, fmt.Sprintf("Error saving datafile: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dataFile)
	}
}

func getDataFilesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		var dataFiles []DataFile
		if result := db.Find(&dataFiles); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching datafiles: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dataFiles); err != nil {
			http.Error(w, "Error encoding datafiles to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func addToolHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at /add-tool")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var tool Tool
		err = json.Unmarshal(body, &tool)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// TODO: Validate Tool

		tempFile, err := ioutil.TempFile("", "*.json")
		if err != nil {
			http.Error(w, "Error creating temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write(body)
		if err != nil {
			http.Error(w, "Error writing temp file", http.StatusInternalServerError)
			return
		}
		tempFile.Close()

		cid, err := ipfs.WrapAndPinFile(tempFile.Name())
		if err != nil {
			http.Error(w, "Error adding to IPFS", http.StatusInternalServerError)
			return
		}

		// Serialize Tool to JSON
		toolJSON, err := json.Marshal(tool)
		if err != nil {
			http.Error(w, "Error serializing tool to JSON", http.StatusInternalServerError)
			return
		}

		// Store serialized Tool in DB
		toolEntity := ToolEntity{
			CID:      cid,
			ToolJSON: string(toolJSON),
		}

		if result := db.Create(&toolEntity); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error creating tool entity: %v", result.Error), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"cid": cid}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func getToolsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		var tools []ToolEntity
		if result := db.Find(&tools); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching tools: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tools); err != nil {
			http.Error(w, "Error encoding tools to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	// Setup database connection
	dsn := "gorm.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&DataFile{}, &User{}, &ToolEntity{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	http.HandleFunc("/healthcheck", healthCheckHandler())
	http.HandleFunc("/user", createUserHandler(db))
	http.HandleFunc("/create-datafile", createDataFileHandler(db))
	http.HandleFunc("/get-datafiles", getDataFilesHandler(db))
	http.HandleFunc("/add-tool", addToolHandler(db))
	http.HandleFunc("/get-tools", getToolsHandler(db))

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(http.DefaultServeMux))
}

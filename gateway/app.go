package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/server"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/labdao/plex/internal/web3"
	"github.com/rs/cors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initializeIo(toolPath string, scatteringMethod string, inputVectors map[string][]string) ([]ipwl.IO, error) {
	// Open the file and load its content
	tool, toolInfo, err := ipwl.ReadToolConfig(toolPath)
	if err != nil {
		return nil, err
	}

	// Check if all kwargs are in the tool's inputs
	for inputKey := range inputVectors {
		if _, exists := tool.Inputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", tool.Inputs)
			return nil, fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}

	// Handle scattering methods and create the inputsList
	var inputsList [][]string
	switch scatteringMethod {
	case "dotProduct":
		// check if all lists have the same length
		var vectorLength int
		for _, v := range inputVectors {
			if vectorLength == 0 {
				vectorLength = len(v)
				continue
			}
			if len(v) != vectorLength {
				return nil, fmt.Errorf("all input arguments must have the same length for dot_product scattering method")
			}
			vectorLength = len(v)
		}
		for i := 0; i < vectorLength; i++ {
			tmp := []string{}
			for _, v := range inputVectors {
				tmp = append(tmp, v[i])
			}
			inputsList = append(inputsList, tmp)
		}
	case "crossProduct":
		cartesian := func(arrs ...[]string) [][]string {
			result := [][]string{{}}

			for _, arr := range arrs {
				var temp [][]string

				for _, res := range result {
					for _, str := range arr {
						product := append([]string{}, res...)
						product = append(product, str)
						temp = append(temp, product)
					}
				}

				result = temp
			}

			return result
		}

		keys := make([]string, 0, len(inputVectors))
		for k := range inputVectors {
			keys = append(keys, k)
		}
		arrays := make([][]string, len(inputVectors))
		for i, k := range keys {
			arrays[i] = inputVectors[k]
		}
		inputsList = cartesian(arrays...)
	default:
		return nil, fmt.Errorf("invalid scattering method: %s", scatteringMethod)
	}

	var userId string

	if web3.IsValidEthereumAddress(os.Getenv("RECIPIENT_WALLET")) {
		userId = os.Getenv("RECIPIENT_WALLET")
	} else {
		fmt.Print("Invalid wallet address detected. Using empty string for user ID.\n")
		userId = ""
	}

	// populate ioJSONGraph based on inputsList
	var ioJSONGraph []ipwl.IO
	for _, inputs := range inputsList {
		io := ipwl.IO{
			Tool:    toolInfo,
			Inputs:  make(map[string]ipwl.FileInput),
			Outputs: make(map[string]ipwl.Output),
			State:   "created",
			ErrMsg:  "",
			UserID:  userId,
		}

		inputKeys := make([]string, 0, len(inputVectors))
		for k := range inputVectors {
			inputKeys = append(inputKeys, k)
		}

		for i, inputValue := range inputs {
			inputKey := inputKeys[i]

			if strings.Count(inputValue, "/") == 1 {
				parts := strings.Split(inputValue, "/")
				cid := parts[0]
				fileName := parts[1]
				if !ipfs.IsValidCID(cid) {
					return nil, fmt.Errorf("invalid CID: %s", cid)
				}
				io.Inputs[inputKey] = ipwl.FileInput{
					Class:    tool.Inputs[inputKey].Type,
					FilePath: fileName,
					IPFS:     cid,
				}
			} else {
				cid, err := ipfs.WrapAndPinFile(inputValue) // Pin the file and get the CID
				if err != nil {
					return nil, err
				}
				io.Inputs[inputKey] = ipwl.FileInput{
					Class:    tool.Inputs[inputKey].Type,
					FilePath: filepath.Base(inputValue), // Use the respective input value from inputsList
					IPFS:     cid,                       // Use the CID returned by WrapAndPinFile
				}
			}
		}

		for outputKey := range tool.Outputs {
			io.Outputs[outputKey] = ipwl.FileOutput{
				Class:    tool.Outputs[outputKey].Type,
				FilePath: "", // Assuming filepath is empty, adapt as needed
				IPFS:     "", // Assuming IPFS is not provided, adapt as needed
			}
		}
		ioJSONGraph = append(ioJSONGraph, io)
	}

	return ioJSONGraph, nil
}

// func initJobHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			sendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
// 			return
// 		}

// 		var requestData struct {
// 			Tool             string   `json:"tool"`
// 			Inputs           []string `json:"inputs"`
// 			ScatteringMethod string   `json:"scatteringMethod"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
// 			sendJSONError(w, "Error parsing request body", http.StatusBadRequest)
// 			return
// 		}

// 		log.Printf("Received request to initialize job: Tool: %s, Inputs: %s, ScatteringMethod: %s\n", requestData.Tool, requestData.Inputs, requestData.ScatteringMethod)

// 		if requestData.ScatteringMethod == "" {
// 			requestData.ScatteringMethod = "dotProduct"
// 		}

// 		inputsMap := make(map[string][]string)
// 		for i, input := range requestData.Inputs {
// 			inputsMap[fmt.Sprintf("input%d", i+1)] = []string{input}
// 		}

// 		ioJson, err := initializeIo(requestData.Tool, requestData.ScatteringMethod, inputsMap)
// 		if err != nil {
// 			sendJSONError(w, fmt.Sprintf("Error initializing IO: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		data, err := json.MarshalIndent(ioJson, "", "  ")
// 		if err != nil {
// 			sendJSONError(w, fmt.Sprintf("Error marshalling IO: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		tempFile, err := ioutil.TempFile("", "io.json")
// 		if err != nil {
// 			sendJSONError(w, "Error creating temp file", http.StatusInternalServerError)
// 			return
// 		}

// 		_, err = tempFile.Write(data)
// 		if err != nil {
// 			sendJSONError(w, "Error writing temp file", http.StatusInternalServerError)
// 			return
// 		}

// 		cid, err := ipfs.PinFile(tempFile.Name())
// 		if err != nil {
// 			sendJSONError(w, "Error pinning file to IPFS", http.StatusInternalServerError)
// 			return
// 		}

// 		job := Job{
// 			InitialIoCID:  cid,
// 			InitialIoJson: string(data),
// 		}

// 		if result := db.Create(&job); result.Error != nil {
// 			sendJSONError(w, fmt.Sprintf("Error saving job: %v", result.Error), http.StatusInternalServerError)
// 			return
// 		}

// 		defer tempFile.Close()

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{"cid": cid})
// 	}
// }

func runJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

	dsn := "gorm.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&models.DataFile{}, &models.User{}, &models.ToolEntity{}, &models.Job{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	// http.HandleFunc("/healthcheck", healthCheckHandler())
	// http.HandleFunc("/user", createUserHandler(db))
	// http.HandleFunc("/create-datafile", createDataFileHandler(db))
	// http.HandleFunc("/get-datafiles", getDataFilesHandler(db))
	// http.HandleFunc("/add-tool", addToolHandler(db))
	// http.HandleFunc("/get-tools", getToolsHandler(db))
	// http.HandleFunc("/init-job", initJobHandler(db))

	mux := server.NewServer(db)

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(mux))
}

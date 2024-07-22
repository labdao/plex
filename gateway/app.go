package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/server"
	"github.com/labdao/plex/gateway/utils"

	"github.com/labdao/plex/internal/s3"

	"github.com/rs/cors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ServeWebApp() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Enable color
		},
	)

	endpoint := os.Getenv("BUCKET_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")

	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	s3Client, err := s3.NewS3Client()
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	} else {
		fmt.Println("S3 client created successfully")

	}

	exists, err := s3Client.BucketExists(bucketName)
	if err != nil {
		log.Fatalf("Failed to check if bucket exists: %v", err)
	}
	if !exists {
		err := s3Client.CreateBucket(bucketName)
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		fmt.Println("Bucket created successfully")
	}

	// Setup database connection
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	privyAppId := os.Getenv("NEXT_PUBLIC_PRIVY_APP_ID")
	publicKey := os.Getenv("PRIVY_PUBLIC_KEY")

	privyVerificationKey := fmt.Sprintf(`-----BEGIN PUBLIC KEY-----
%s
-----END PUBLIC KEY-----`, publicKey)

	middleware.SetupConfig(privyAppId, privyVerificationKey)

	// DSN for gorm.Open
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, dbname)

	// URL-encoded DSN for migrate.New
	encodedPassword := url.QueryEscape(password)
	migrateDSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, encodedPassword, host, dbname)

	// Run Raw SQL Migrations First using golang-migrate
	m, err := migrate.New(
		"file://gateway/migrations",
		migrateDSN,
	)
	if err != nil {
		log.Fatalf("Could not create migration: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while migrating the database: %v", err)
	}

	// If needed use log level debug or info. Default set to silent to avoid noisy logs
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to configure database connection pool")
	}
	sqlDB.SetMaxIdleConns(10)           // Set the maximum number of idle connections
	sqlDB.SetMaxOpenConns(25)           // Set the maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Connections are recycled every hour

	// Migrate the schema
	if err := db.AutoMigrate(&models.File{}, &models.User{}, &models.Model{}, &models.Job{}, &models.Tag{}, &models.Transaction{}, &models.InferenceEvent{}, &models.FileEvent{}, &models.UserEvent{}, &models.Organization{}, &models.Design{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Insert default organization if it doesn't exist
	var org models.Organization
	result := db.FirstOrCreate(&org, models.Organization{Name: "no_org"})
	if result.Error != nil {
		log.Printf("Error ensuring default organization exists: %v", result.Error)
	} else {
		log.Println("Default organization ensured in database")
	}

	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")
	log.Printf("Initial STRIPE_WEBHOOK_SECRET_KEY: %s\n", stripeWebhookSecret)

	if stripeWebhookSecret == "" {
		log.Println("STRIPE_WEBHOOK_SECRET_KEY not set, attempting to read from file")
		stripeWebhookSecretBytes, err := os.ReadFile("/var/secrets/stripe/secret.txt")
		if err != nil {
			log.Fatalf("Failed to read Stripe webhook signing secret: %v", err)
		}
		stripeWebhookSecret = strings.TrimSpace(string(stripeWebhookSecretBytes))

		os.Setenv("STRIPE_WEBHOOK_SECRET_KEY", stripeWebhookSecret)
		log.Printf("STRIPE_WEBHOOK_SECRET_KEY set from file: %s\n", stripeWebhookSecret)
	} else {
		log.Println("STRIPE_WEBHOOK_SECRET_KEY is already set")
	}

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL"), "http://localhost:3000", "http://frontend:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
	})

	mux := server.NewServer(db, s3Client)

	maxWorkers := utils.GetEnvAsInt("MAX_WORKERS", 4)

	// Start queue watcher in a separate goroutine
	go func() {
		for {
			if err := utils.StartJobQueues(db, maxWorkers); err != nil {
				fmt.Printf("unexpected error processing job queues: %v\n", err)
				time.Sleep(5 * time.Second) // wait for 5 seconds before retrying
			}
		}
	}()

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(mux))
}

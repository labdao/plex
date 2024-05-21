package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/server"
	"github.com/labdao/plex/gateway/utils"
	"github.com/minio/minio-go/v7"

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
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")
	useSSL := os.Getenv("BUCKET_USE_SSL") == "true"
	bucketName := os.Getenv("BUCKET_NAME")

	minioClient, err := s3.NewMinIOClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	} else {
		fmt.Println("Minio client created successfully")
	}

	exists, err := minioClient.Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("Failed to check if bucket exists: %v", err)
	}
	if !exists {
		err = minioClient.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&models.DataFile{}, &models.User{}, &models.Tool{}, &models.Job{}, &models.Tag{}, &models.Transaction{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL"), "http://localhost:3000", "https://editor.swagger.io", "https://editor-next.swagger.io"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
	})

	mux := server.NewServer(db, minioClient)

	// Start queue watcher in a separate goroutine
	go func() {
		for {
			if err := utils.StartJobQueues(db); err != nil {
				fmt.Printf("unexpected error processing job queues: %v\n", err)
				time.Sleep(5 * time.Second) // wait for 5 seconds before retrying
			}
		}
	}()

	// Start running jobs watcher in a separate goroutine
	go func() {
		for {
			if err := utils.MonitorRunningJobs(db); err != nil {
				fmt.Printf("unexpected error monitoring running jobs: %v\n", err)
				time.Sleep(5 * time.Second) // wait for 5 seconds before retrying
			} else {
				break // exit the loop if no error (optional based on your use case)
			}
		}
	}()

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(mux))
}

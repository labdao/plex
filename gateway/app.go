package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/server"
	"github.com/labdao/plex/gateway/utils"

	"github.com/rs/cors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ServeWebApp() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	// Setup database connection
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

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
	if err := db.AutoMigrate(&models.DataFile{}, &models.User{}, &models.Tool{}, &models.Job{}, &models.Tag{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Set up CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL"), "http://localhost:3000", "https://editor.swagger.io", "https://editor-next.swagger.io"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH"},
	})

	mux := server.NewServer(db)

	// Start queue watcher in a separate goroutine
	go func() {
		if err := utils.StartJobQueues(db); err != nil {
			// There should definitely be log alerts set up around this message
			fmt.Printf("unexpected error processing job queues: %v", err)
		}
	}()

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(mux))
}

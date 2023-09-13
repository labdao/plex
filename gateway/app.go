package gateway

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/server"

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

	mux := server.NewServer(db)

	// Start the server with CORS middleware
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(mux))
}

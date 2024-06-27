package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/handlers"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/internal/s3"

	"gorm.io/gorm"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func createProtectedRouteHandler(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(db)(handler)
	}
}

func createAdminProtectedRouteHandler(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.AdminCheckMiddleware(db)(handler)
	}
}

func NewServer(db *gorm.DB, s3c *s3.S3Client) *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	protected := createProtectedRouteHandler(db)
	adminProtected := createAdminProtectedRouteHandler(db)

	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler())

	router.HandleFunc("/user", handlers.AddUserHandler(db)).Methods("POST")
	router.HandleFunc("/user", protected(handlers.GetUserHandler(db))).Methods("GET")

	router.HandleFunc("/models", protected(adminProtected(handlers.AddModelHandler(db, s3c)))).Methods("POST")
	router.HandleFunc("/models/{id}", protected(handlers.GetModelHandler(db))).Methods("GET")
	router.HandleFunc("/models", protected(handlers.ListModelsHandler(db))).Methods("GET")
	router.HandleFunc("/models/{id}", protected(adminProtected(handlers.UpdateModelHandler(db)))).Methods("PUT")

	router.HandleFunc("/files", protected(handlers.AddFileHandler(db, s3c))).Methods("POST")
	router.HandleFunc("/files/{id}", protected(handlers.GetFileHandler(db))).Methods("GET")
	router.HandleFunc("/files/{id}", protected(handlers.UpdateFileHandler(db))).Methods("PUT")
	router.HandleFunc("/files/{id}/download", protected(handlers.DownloadFileHandler(db, s3c))).Methods("GET")
	router.HandleFunc("/files", protected(handlers.ListFilesHandler(db))).Methods("GET")

	router.HandleFunc("/checkpoints/{experimentID}/get-data", protected(handlers.GetExperimentCheckpointDataHandler(db))).Methods("GET")

	// router.HandleFunc("/experiments", protected(handlers.AddExperimentHandler(db))).Methods("POST")
	// router.HandleFunc("/experiments", protected(handlers.ListExperimentsHandler(db))).Methods("GET")
	// router.HandleFunc("/experiments/{experimentID}", protected(handlers.GetExperimentHandler(db))).Methods("GET")
	// router.HandleFunc("/experiments/{experimentID}", protected(handlers.UpdateExperimentHandler(db))).Methods("PUT")
	// router.HandleFunc("/experiments/{experimentID}/add-job", protected(handlers.AddJobToExperimentHandler(db))).Methods("PUT")

	// router.HandleFunc("/jobs/{jobID}", protected(handlers.GetJobHandler(db))).Methods("GET")
	// // router.HandleFunc("/jobs/{bacalhauJobID}/logs", handlers.StreamJobLogsHandler).Methods("GET")
	// router.HandleFunc("/queue-summary", handlers.GetJobsQueueSummaryHandler(db)).Methods("GET")

	router.HandleFunc("/tags", protected(handlers.AddTagHandler(db))).Methods("POST")
	router.HandleFunc("/tags", protected(handlers.ListTagsHandler(db))).Methods("GET")

	router.HandleFunc("/api-keys", protected(handlers.AddAPIKeyHandler(db))).Methods("POST")
	router.HandleFunc("/api-keys", protected(handlers.ListAPIKeysHandler(db))).Methods("GET")

	// router.HandleFunc("/stripe", handlers.StripeFulfillmentHandler(db)).Methods("POST")
	// router.HandleFunc("/stripe/checkout", protected(handlers.StripeCreateCheckoutSessionHandler(db))).Methods("POST")
	// router.HandleFunc("/transactions", protected(handlers.ListTransactionsHandler(db))).Methods("GET")
	// router.HandleFunc("/transactions-summary", protected(handlers.SummaryTransactionsHandler(db))).Methods("GET")

	return router
}

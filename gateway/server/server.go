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

	router.HandleFunc("/tools", protected(adminProtected(handlers.AddToolHandler(db, s3c)))).Methods("POST")
	router.HandleFunc("/tools/{cid}", protected(handlers.GetToolHandler(db))).Methods("GET")
	router.HandleFunc("/tools", protected(handlers.ListToolsHandler(db))).Methods("GET")
	router.HandleFunc("/tools/{cid}", protected(adminProtected(handlers.UpdateToolHandler(db)))).Methods("PUT")

	router.HandleFunc("/datafiles", protected(handlers.AddDataFileHandler(db, s3c))).Methods("POST")
	router.HandleFunc("/datafiles/{cid}", protected(handlers.GetDataFileHandler(db))).Methods("GET")
	router.HandleFunc("/datafiles/{cid}", protected(handlers.UpdateDataFileHandler(db))).Methods("PUT")
	// pass in minioClient
	router.HandleFunc("/datafiles/{cid}/download", protected(handlers.DownloadDataFileHandler(db, minioClient))).Methods("GET")
	router.HandleFunc("/datafiles", protected(handlers.ListDataFilesHandler(db))).Methods("GET")

	router.HandleFunc("/checkpoints/{flowID}", protected(handlers.ListFlowCheckpointsHandler(db))).Methods("GET")
	router.HandleFunc("/checkpoints/{flowID}/get-data", protected(handlers.GetFlowCheckpointDataHandler(db))).Methods("GET")
	router.HandleFunc("/checkpoints/{flowID}/{jobID}", protected(handlers.ListJobCheckpointsHandler(db))).Methods("GET")
	router.HandleFunc("/checkpoints/{flowID}/{jobID}/get-data", protected(handlers.GetJobCheckpointDataHandler(db))).Methods("GET")

	router.HandleFunc("/flows", protected(handlers.AddFlowHandler(db))).Methods("POST")
	router.HandleFunc("/flows", protected(handlers.ListFlowsHandler(db))).Methods("GET")
	router.HandleFunc("/flows/{flowID}", protected(handlers.GetFlowHandler(db))).Methods("GET")
	router.HandleFunc("/flows/{flowID}", protected(handlers.UpdateFlowHandler(db))).Methods("PUT")
	router.HandleFunc("/flows/{flowID}/add-job", protected(handlers.AddJobToFlowHandler(db))).Methods("PUT")

	router.HandleFunc("/jobs/{jobID}", protected(handlers.GetJobHandler(db))).Methods("GET")
	router.HandleFunc("/jobs/{bacalhauJobID}/logs", handlers.StreamJobLogsHandler).Methods("GET")
	router.HandleFunc("/queue-summary", handlers.GetJobsQueueSummaryHandler(db)).Methods("GET")

	router.HandleFunc("/tags", protected(handlers.AddTagHandler(db))).Methods("POST")
	router.HandleFunc("/tags", protected(handlers.ListTagsHandler(db))).Methods("GET")

	router.HandleFunc("/api-keys", protected(handlers.AddAPIKeyHandler(db))).Methods("POST")
	router.HandleFunc("/api-keys", protected(handlers.ListAPIKeysHandler(db))).Methods("GET")

	router.HandleFunc("/stripe", handlers.StripeFullfillmentHandler(db)).Methods("POST")
	router.HandleFunc("/stripe/checkout", protected(handlers.StripeCreateCheckoutSessionHandler(db))).Methods("GET")
	router.HandleFunc("/transactions", protected(handlers.ListTransactionsHandler(db))).Methods("GET")
	router.HandleFunc("/transactions-summary", protected(handlers.SummaryTransactionsHandler(db))).Methods("GET")

	return router
}

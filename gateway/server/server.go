package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/handlers"
	"github.com/labdao/plex/gateway/middleware"

	"gorm.io/gorm"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func createProtectedRouteHandler(db *gorm.DB, privyPublicKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.JWTMiddleware(db, privyPublicKey)(handler)
	}
}

func NewServer(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	privyVerificationKey := os.Getenv("PRIVY_VERIFICATION_KEY")
	protected := createProtectedRouteHandler(db, privyVerificationKey)

	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler())
	router.HandleFunc("/user", handlers.AddUserHandler(db))

	router.HandleFunc("/tools", protected(handlers.AddToolHandler(db))).Methods("POST")
	router.HandleFunc("/tools/{cid}", protected(handlers.GetToolHandler(db))).Methods("GET")
	router.HandleFunc("/tools", protected(handlers.ListToolsHandler(db))).Methods("GET")

	router.HandleFunc("/datafiles", protected(handlers.AddDataFileHandler(db))).Methods("POST")
	router.HandleFunc("/datafiles/{cid}", protected(handlers.GetDataFileHandler(db))).Methods("GET")
	router.HandleFunc("/datafiles/{cid}/download", protected(handlers.DownloadDataFileHandler(db))).Methods("GET")
	router.HandleFunc("/datafiles", protected(handlers.ListDataFilesHandler(db))).Methods("GET")

	router.HandleFunc("/flows", protected(handlers.AddFlowHandler(db))).Methods("POST")
	router.HandleFunc("/flows", protected(handlers.ListFlowsHandler(db))).Methods("GET")
	router.HandleFunc("/flows/{cid}", protected(handlers.GetFlowHandler(db))).Methods("GET")
	router.HandleFunc("/flows/{cid}", protected(handlers.UpdateFlowHandler(db))).Methods("PATCH")

	router.HandleFunc("/jobs/{bacalhauJobID}", protected(handlers.GetJobHandler(db))).Methods("GET")
	router.HandleFunc("/jobs/{bacalhauJobID}", protected(handlers.UpdateJobHandler(db))).Methods("PATCH")
	router.HandleFunc("/jobs/{bacalhauJobID}/logs", protected(handlers.StreamJobLogsHandler)).Methods("GET")

	router.HandleFunc("/tags", protected(handlers.AddTagHandler(db))).Methods("POST")
	router.HandleFunc("/tags", protected(handlers.ListTagsHandler(db))).Methods("GET")

	return router
}

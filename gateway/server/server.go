package server

import (
	"log"
	"net/http"

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

func NewServer(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	memberOnlyRouter := router.PathPrefix("/member-only").Subrouter()
	memberOnlyRouter.Use(middleware.MemberOnlyMiddleware(db))

	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler())
	router.HandleFunc("/user", handlers.AddUserHandler(db)).Methods("POST")
	router.HandleFunc("/user/{walletAddress}", handlers.CheckUserMemberStatusHandler(db)).Methods("GET")

	router.HandleFunc("/tools", handlers.AddToolHandler(db)).Methods("POST")
	router.HandleFunc("/tools/{cid}", handlers.GetToolHandler(db)).Methods("GET")
	router.HandleFunc("/tools", handlers.ListToolsHandler(db)).Methods("GET")

	router.HandleFunc("/datafiles", handlers.AddDataFileHandler(db)).Methods("POST")
	router.HandleFunc("/datafiles/{cid}", handlers.GetDataFileHandler(db)).Methods("GET")
	router.HandleFunc("/datafiles", handlers.ListDataFilesHandler(db)).Methods("GET")

	memberOnlyRouter.HandleFunc("/flows", handlers.AddFlowHandler(db)).Methods("POST")
	memberOnlyRouter.HandleFunc("/flows", handlers.ListFlowsHandler(db)).Methods("GET")
	memberOnlyRouter.HandleFunc("/flows/{cid}", handlers.GetFlowHandler(db)).Methods("GET")
	memberOnlyRouter.HandleFunc("/flows/{cid}", handlers.UpdateFlowHandler(db)).Methods("PATCH")

	memberOnlyRouter.HandleFunc("/jobs/{bacalhauJobID}", handlers.GetJobHandler(db)).Methods("GET")
	memberOnlyRouter.HandleFunc("/jobs/{bacalhauJobID}", handlers.UpdateJobHandler(db)).Methods("PATCH")
	memberOnlyRouter.HandleFunc("/jobs/{bacalhauJobID}/logs", handlers.StreamJobLogsHandler).Methods("GET")

	return router
}

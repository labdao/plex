package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/handlers"

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

	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler())
	router.HandleFunc("/user", handlers.AddUserHandler(db))

	router.HandleFunc("/tools", handlers.AddToolHandler(db)).Methods("POST")
	router.HandleFunc("/tools", handlers.GetToolsHandler(db)).Methods("GET")
	router.HandleFunc("/tools/{cid}", handlers.GetToolHandler(db)).Methods("GET")

	router.HandleFunc("/datafiles", handlers.AddDataFileHandler(db)).Methods("POST")
	router.HandleFunc("/datafiles", handlers.GetDataFilesHandler(db)).Methods("GET")
	router.HandleFunc("/datafiles/{cid}", handlers.GetDataFileHandler(db)).Methods("GET")

	router.HandleFunc("/flows", handlers.AddFlowHandler(db)).Methods("POST")
	router.HandleFunc("/flows", handlers.ListFlowsHandler(db)).Methods("GET")

	return router
}

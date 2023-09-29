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

	router.HandleFunc("/tool", handlers.AddToolHandler(db)).Methods("POST")
	router.HandleFunc("/get-tools", handlers.GetToolsHandler(db)).Methods("GET")
	router.HandleFunc("/get-tools/{cid}", handlers.GetToolHandler(db)).Methods("GET")

	router.HandleFunc("/add-datafile", handlers.AddDataFileHandler(db)).Methods("POST")
	router.HandleFunc("/get-datafiles", handlers.GetDataFilesHandler(db)).Methods("GET")
	router.HandleFunc("/get-datafiles/{cid}", handlers.GetDataFileHandler(db)).Methods("GET")

	router.HandleFunc("/graph", handlers.AddGraphHandler(db)).Methods("POST")

	// router.HandleFunc("/init-job", handlers.InitJobHandler(db)).Methods("POST")
	// router.HandleFunc("/get-jobs", handlers.GetJobsHandler(db)).Methods("GET")
	// router.HandleFunc("/get-jobs/{cid}", handlers.GetJobHandler(db)).Methods("GET")
	// router.HandleFunc("/run-job", handlers.RunJobHandler(db)).Methods("POST")

	return router
}

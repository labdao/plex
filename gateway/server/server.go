package server

import (
	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/handlers"

	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler())
	router.HandleFunc("/user", handlers.AddUserHandler(db))

	router.HandleFunc("/add-tool", handlers.AddToolHandler(db)).Methods("POST")
	router.HandleFunc("/get-tools", handlers.GetToolsHandler(db)).Methods("GET")
	router.HandleFunc("/get-tools/{id}", handlers.GetToolHandler(db)).Methods("GET")

	router.HandleFunc("/add-datafile", handlers.AddDataFileHandler(db)).Methods("POST")
	router.HandleFunc("/get-datafiles", handlers.GetDataFilesHandler(db)).Methods("GET")
	router.HandleFunc("/get-datafiles/{id}", handlers.GetDataFileHandler(db)).Methods("GET")

	router.HandleFunc("/init-job", handlers.InitJobHandler(db)).Methods("POST")
	// router.HandleFunc("/get-jobs", handlers.GetJobsHandler(db)).Methods("GET")
	// router.HandleFunc("/get-jobs/{id}", handlers.GetJobHandler(db)).Methods("GET")
	// router.HandleFunc("/run-job", handlers.RunJobHandler(db)).Methods("POST")

	return router
}

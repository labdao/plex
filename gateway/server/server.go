package server

import (
	"net/http"

	"github.com/labdao/plex/gateway/handlers"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", handlers.HealthCheckHandler())
	mux.HandleFunc("/user", handlers.AddUserHandler(db))
	mux.HandleFunc("/add-tool", handlers.AddToolHandler(db))
	mux.HandleFunc("/get-tools", handlers.GetToolsHandler(db))
	mux.HandleFunc("/add-datafile", handlers.AddDataFileHandler(db))
	mux.HandleFunc("/get-datafiles", handlers.GetDataFilesHandler(db))
	// mux.HandleFunc("/init-job", handlers.InitJobHandler(db))
	// mux.HandleFunc("/run-job", handlers.RunJobHandler(db))

	return mux
}

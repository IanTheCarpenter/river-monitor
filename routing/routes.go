package router

import (
	"net/http"

	"github.com/IanTheCarpenter/river-monitor/handlers"
)

func Init() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", handlers.GetRoot)
	router.HandleFunc("GET /rivers/{river_name}", handlers.GetFullRiverReport)

	return router
}

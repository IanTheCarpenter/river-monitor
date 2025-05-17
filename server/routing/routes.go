package routing

import (
	"net/http"

	"github.com/IanTheCarpenter/river-monitor/server/handlers"
)

func Init() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /", handlers.GetRoot)
	r.HandleFunc("GET /rivers/{river_name}", handlers.GetFullRiverReport)

	return r
}

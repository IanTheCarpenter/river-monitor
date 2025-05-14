package main

import (
	"net/http"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/forecaster"
	"github.com/IanTheCarpenter/river-monitor/routing"
)

func main() {
	db.Init()
	// river_data.Update()

	go forecaster.Init()

	server := http.Server{
		Addr:    ":3000",
		Handler: routing.Init(),
	}
	server.ListenAndServe()

}

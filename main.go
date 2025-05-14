package main

import (
	"net/http"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/forecaster"
	"github.com/IanTheCarpenter/river-monitor/river_data"
	"github.com/IanTheCarpenter/river-monitor/routing"
)

func main() {
	river_data.Update()
	db.Init()

	go forecaster.Init()

	server := http.Server{
		Addr:    ":3000",
		Handler: routing.Init(),
	}
	server.ListenAndServe()

}

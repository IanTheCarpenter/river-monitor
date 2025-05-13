package main

import (
	"net/http"

	connection "github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/forecaster"
	router "github.com/IanTheCarpenter/river-monitor/routing"
)

func main() {
	connection.Init()

	go forecaster.Init()

	server := http.Server{
		Addr:    ":3000",
		Handler: router.Init(),
	}
	server.ListenAndServe()

}

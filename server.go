package main

import (
	"net/http"

	"github.com/sergeyt/pandora/modules/api"
	"github.com/sergeyt/pandora/modules/config"
	log "github.com/sirupsen/logrus"
)

var server *http.Server

func startServer() {
	log.Printf("listening %s\n", config.ServerAddr)

	server = &http.Server{
		Addr:    config.ServerAddr,
		Handler: api.NewHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func stopServer() {
	log.Println("shutting down")
	server.Shutdown(nil)
}

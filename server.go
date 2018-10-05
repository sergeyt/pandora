package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
)

var server *http.Server

func startServer() {
	dgraph.InitSchema()

	fmt.Printf("listening %s\n", config.ServerAddr)

	server = &http.Server{
		Addr:    config.ServerAddr,
		Handler: makeAPIHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func stopServer() {
	fmt.Println("shutting down")
	server.Shutdown(nil)
}

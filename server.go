package main

import (
	"fmt"
	"log"
	"net/http"
)

var server *http.Server

func startServer() {
	initSchema()

	fmt.Printf("listening %s\n", config.Addr)

	server = &http.Server{
		Addr:    config.Addr,
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

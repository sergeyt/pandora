package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/handlers"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/geoip"
)

func makeAPIHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Group(auth.AuthAPI)
	//r.Group(elasticsearch.SearchAPI)
	r.Group(dataAPI)
	r.Group(fileAPI)
	r.Group(geoip.RegisterAPI)

	return r
}

func Logger(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

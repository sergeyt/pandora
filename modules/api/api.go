package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/handlers"
	"github.com/sergeyt/pandora/modules/auth"
)

// NewHandler makes http.Handler to serve HTTP API
func NewHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logMidleware)
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

	r.Group(auth.RegisterAPI)
	//r.Group(elasticsearch.SearchAPI)
	r.Group(dataAPI)
	r.Group(fileAPI)
	//r.Group(geoip.RegisterAPI)
	r.Group(adminAPI)

	return r
}

func logMidleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

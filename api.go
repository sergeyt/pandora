package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func makeAPIHandler() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(authAPI)

	mux.Group(func(r chi.Router) {
		r = r.With(authMiddleware)
		r.Get("/api/event/stream", GetEventStream)
		r.Get("/api/event/stream/{channel}", GetEventStream)
	})

	mux.Group(uploadAPI)
	mux.Group(dataAPI)

	return mux
}

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

	authAPI(mux)

	mux.Get("/api/event/stream", GetEventStream)
	mux.Get("/api/event/stream/{channel}", GetEventStream)

	uploadAPI(mux)

	mux.Group(dataAPI)

	return mux
}

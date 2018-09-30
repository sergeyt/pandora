package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gocontrib/pubsub/sse"
)

func makeAPIHandler() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(authAPI)

	mux.Group(func(r chi.Router) {
		// TODO configurable api path
		r.Get("/api/event/stream", sse.GetEventStream)
		r.Get("/api/event/stream/{channel}", sse.GetEventStream)
	})

	mux.Group(uploadAPI)
	mux.Group(dataAPI)

	return mux
}

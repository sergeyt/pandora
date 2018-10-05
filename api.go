package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/elasticsearch"
)

func makeAPIHandler() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(auth.AuthAPI)
	mux.Group(elasticsearch.SearchAPI)
	mux.Group(dataAPI)

	return mux
}

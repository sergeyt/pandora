package main

import (
	"github.com/go-chi/chi"
	"github.com/gocontrib/auth"
)

func authAPI(mux chi.Router) {
	config := makeAuthConfig()
	mux.Post("/api/login", auth.LoginHandlerFunc(config))
}

func makeAuthConfig() *auth.Config {
	return &auth.Config{
		UserStore: makeUserStore(),
	}
}

package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocontrib/auth"
)

var (
	authConfig  = makeAuthConfig()
	requireUser = auth.RequireUser(authConfig)
)

func authAPI(mux chi.Router) {
	mux.Post("/api/login", auth.LoginHandlerFunc(authConfig))
}

func makeAuthConfig() *auth.Config {
	return &auth.Config{
		UserStore: makeUserStore(),
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// support local_admin calls
		if r.Header.Get("Authorization") == "local_admin" {
			systemUser := &UserInfo{
				ID:    "0x0",
				Name:  "$system",
				Email: "",
				Admin: true,
			}
			ctx := auth.WithUser(r.Context(), systemUser)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			requireUser(next).ServeHTTP(w, r)
		}
	})
}

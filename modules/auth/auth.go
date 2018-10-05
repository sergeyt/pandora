package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	authbase "github.com/gocontrib/auth"
)

var (
	authConfig  = makeAuthConfig()
	requireUser = authbase.RequireUser(authConfig)
)

func AuthAPI(mux chi.Router) {
	mux.Post("/api/login", authbase.LoginHandlerFunc(authConfig))
}

func makeAuthConfig() *authbase.Config {
	return &authbase.Config{
		UserStore: makeUserStore(),
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// support local_admin calls
		if r.Header.Get("Authorization") == "local_admin" {
			systemUser := &authbase.UserInfo{
				ID:    "system",
				Name:  "system",
				Email: "",
				Admin: true,
			}
			ctx := authbase.WithUser(r.Context(), systemUser)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			requireUser(next).ServeHTTP(w, r)
		}
	})
}

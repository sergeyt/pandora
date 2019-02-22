package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func adminAPI(r chi.Router) {
	r = r.With(adminSecret)
}

func adminSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("secret") != os.Getenv("ADMIN_SECRET") {
			http.Error(w, "bad secret", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

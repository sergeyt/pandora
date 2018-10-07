package auth

import (
	"net/http"
	"strings"

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
		claims := extractClaims(r)
		userID := get(claims, "user_id")
		userName := get(claims, "user_name")
		email := get(claims, "email")
		role := get(claims, "role")

		if len(userID) > 0 && len(userName) > 0 {
			// TODO validate user id
			user := &authbase.UserInfo{
				ID:    userID,
				Name:  userName,
				Email: email,
				Admin: true,
				Claims: map[string]interface{}{
					"email": email,
					"role":  role,
				},
			}
			ctx := authbase.WithUser(r.Context(), user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

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

func extractClaims(r *http.Request) map[string]string {
	result := make(map[string]string)
	prefix := "Token-Claim-"
	for k, v := range r.Header {
		if strings.HasPrefix(k, prefix) {
			name := strings.ToLower(k[len(prefix):])
			result[name] = v[0]
		}
	}
	return result
}

func get(m map[string]string, k string) string {
	s, ok := m[k]
	if ok {
		return s
	}
	return ""
}

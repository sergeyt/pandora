package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	authbase "github.com/gocontrib/auth"
	"github.com/gocontrib/auth/oauth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/vk"
	"github.com/sergeyt/pandora/modules/config"
)

var (
	authConfig  = makeAuthConfig()
	requireUser = authbase.RequireUser(authConfig)
)

// RegisterAPI registers authentication HTTP API
func RegisterAPI(mux chi.Router) {
	mux.Post("/api/login", authbase.LoginHandlerFunc(authConfig))
	mux.Post("/api/register", authbase.RegisterHandlerFunc(authConfig))

	oauth.WithProviders(authConfig, "vk", vk.New, "google", google.New, "facebook", facebook.New)
	oauth.RegisterAPI(mux, authConfig)
}

func makeAuthConfig() *authbase.Config {
	userStore := makeUserStore()
	return &authbase.Config{
		UserStore:   userStore,
		UserStoreEx: userStore,
		ServerURL:   config.ServerURL(),
	}
}

// Middleware is authentication HTTP middleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		claims := extractClaims(r)
		userID := get(claims, "user_id")
		userName := get(claims, "user_name")
		email := get(claims, "email")
		role := get(claims, "role")

		if len(userID) > 0 && len(userName) > 0 {
			var user authbase.User = &authbase.UserInfo{
				ID:    userID,
				Name:  userName,
				Email: email,
				Admin: true,
				Claims: map[string]interface{}{
					"email": email,
					"role":  role,
				},
			}

			ctx := r.Context()

			if userID != "system" {
				user, err = authConfig.UserStore.FindUserByID(ctx, userID)
				if err != nil {
					authbase.SendError(w, authbase.ErrUserNotFound.WithCause(err))
					return
				}
			}

			ctx = authbase.WithUser(ctx, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		requireUser(next).ServeHTTP(w, r)
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

package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	authbase "github.com/gocontrib/auth"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
)

func init() {
	// TODO get host from config
	host := "http://localhost:4201"
	goth.UseProviders(
		facebook.New(
			os.Getenv("FACEBOOK_KEY"),
			os.Getenv("FACEBOOK_SECRET"),
			host+"/api/oauth/callback/facebook",
		),
	)
}

func OAuthAPI(mux chi.Router) {
	mux.Get("/api/oauth/login/:provider", func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating
		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			completeOAuthFlow(w, r, user)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})

	mux.Get("/api/oauth/logout/:provider", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	mux.Get("/api/oauth/callback/:provider", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		completeOAuthFlow(w, r, user)
	})
}

func completeOAuthFlow(w http.ResponseWriter, r *http.Request, account goth.User) {
	ctx := r.Context()
	userStore := makeUserStore()
	user, err := userStore.FindUserByEmail(ctx, account.Email)
	if err != nil {
		data := CreateUserData{
			Name:      account.Name,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Email:     account.Email,
			// TODO login should be unique
			// Login: account.NickName,
		}
		user, err = userStore.CreateUser(ctx, data)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}

	// TODO link external account to the user

	authbase.WriteLoginResponse(w, r, authConfig, user)
}

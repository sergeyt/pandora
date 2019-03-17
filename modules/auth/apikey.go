package auth

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sergeyt/pandora/modules/apiutil"
	log "github.com/sirupsen/logrus"
)

// RequireAPIKey middleware to check API key
func RequireAPIKey(next http.Handler) http.Handler {
	apiKeySecret := os.Getenv("API_KEY_SECRET")
	if len(apiKeySecret) == 0 {
		panic("API_KEY_SECRET is not defined")
	}

	getAPIKey := func(r *http.Request) string {
		apiKey := r.URL.Query().Get("key")
		if len(apiKey) > 0 {
			return apiKey
		}
		return r.Header.Get("X-API-Key")
	}

	appSecret := ""
	getSecret := func(token *jwt.Token) (interface{}, error) {
		v, ok := token.Header["app_secret"]
		if ok {
			s, isStr := v.(string)
			if !isStr {
				return nil, fmt.Errorf("app_secret is string")
			}
			if len(s) == 0 {
				return nil, fmt.Errorf("missing app_secret")
			}
			appSecret = s
			result := []byte(apiKeySecret + s)
			return result, nil
		}
		return nil, fmt.Errorf("app_secret is not defined")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := getAPIKey(r)
		if len(apiKey) == 0 {
			err := fmt.Errorf("missing API key")
			log.Errorf(err.Error())
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}
		parser := new(jwt.Parser)
		parser.SkipClaimsValidation = true
		claims := jwt.MapClaims{}
		token, err := parser.ParseWithClaims(apiKey, claims, getSecret)
		if err != nil {
			log.Errorf("jwt.Parser.ParseWithClaims fail: %v", err)
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			err = fmt.Errorf("bad API key, invalid token")
			log.Errorf(err.Error())
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}
		v, ok := claims["app_id"]
		if !ok {
			err = fmt.Errorf("bad API key, missing app_id")
			log.Errorf(err.Error())
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}
		appID, ok := v.(string)
		if !ok {
			err = fmt.Errorf("bad API key, app_id is not string")
			log.Errorf(err.Error())
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}
		// TODO check app id from database
		knownID := os.Getenv("APP_ID")
		knownSecret := os.Getenv("APP_SECRET")
		if appID != knownID || appSecret != knownSecret {
			err = fmt.Errorf("app %s is not registered", appID)
			log.Errorf(err.Error())
			apiutil.SendError(w, err, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middlewarex

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/factly/x/requestx"
	"github.com/spf13/viper"
)

type ValidationBody struct {
	SecretToken string `json:"secret_token" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
}

// ValidateAPIToken validates the API tokens from kavach server
func ValidateAPIToken(appName string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			base64Token := strings.TrimPrefix(authHeader, "Basic ")

			tokenStr, err := base64.StdEncoding.DecodeString(base64Token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(string(tokenStr), ":")
			tokenBody := ValidationBody{
				AccessToken: tokenParts[0],
				SecretToken: tokenParts[1],
			}

			res, err := requestx.Request("POST", viper.GetString("kavach_url")+"/applications/"+appName+"/validateToken", tokenBody, nil)

			if err != nil || res.StatusCode != http.StatusOK {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

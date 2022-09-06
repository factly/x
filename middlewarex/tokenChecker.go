package middlewarex

import (
	"fmt"
	"net/http"

	"github.com/factly/x/requestx"
	"github.com/spf13/viper"
)

type ValidationBody struct {
	Token string `json:"token" validate:"required"`
}

// ValidateAPIToken validates the API tokens from kavach server
func ValidateAPIToken(header string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sID, err := GetSpace(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authHeader := r.Header.Get(header)
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenBody := ValidationBody{
				Token: authHeader,
			}

			res, err := requestx.Request("POST", viper.GetString("kavach_url")+"/spaces/"+fmt.Sprintf("%d", sID)+"/validateToken", tokenBody, nil)

			if err != nil || res.StatusCode != http.StatusOK {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

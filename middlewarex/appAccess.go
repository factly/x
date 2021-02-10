package middlewarex

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// CheckAccess middleware to check if user can access the application
func CheckAccess(appSlug string, index int, GetOrg func(ctx context.Context) (int, error)) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Split(strings.Trim(r.URL.Path, "/"), "/")[index] != "spaces" {

				uID, err := GetUser(r.Context())
				if err != nil {
					w.Header().Add("Location", viper.GetString("kavach_public_url"))
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}

				oID, err := GetOrg(r.Context())
				if err != nil {
					w.Header().Add("Location", viper.GetString("kavach_public_url"))
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}

				path := fmt.Sprint("/organisations/", oID, "/applications/", appSlug, "/access")
				req, err := http.NewRequest("GET", viper.GetString("kavach_url")+path, nil)
				if err != nil {
					w.Header().Add("Location", viper.GetString("kavach_public_url"))
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User", fmt.Sprint(uID))

				client := &http.Client{}
				resp, err := client.Do(req)

				if err != nil {
					w.Header().Add("Location", viper.GetString("kavach_public_url"))
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}

				if resp.StatusCode > 400 && resp.StatusCode <= 500 {
					w.Header().Add("Location", viper.GetString("kavach_public_url"))
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

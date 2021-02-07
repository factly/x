package middlewarex

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/factly/x/errorx"
	"github.com/spf13/viper"
)

// CheckAccess middleware to check if user can access the application
func CheckAccess(appSlug string, GetOrg func(ctx context.Context) (int, error)) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Split(strings.Trim(r.URL.Path, "/"), "/")[1] != "spaces" {

				uID, err := GetUser(r.Context())
				if err != nil {
					errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
					return
				}

				oID, err := GetOrg(r.Context())
				if err != nil {
					errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
					return
				}

				path := fmt.Sprint("/organisations/", oID, "/applications/", appSlug, "/access")
				req, err := http.NewRequest("GET", viper.GetString("kavach_url")+path, nil)
				if err != nil {
					errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
					return
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User", fmt.Sprint(uID))

				client := &http.Client{}
				resp, err := client.Do(req)

				if err != nil {
					errorx.Render(w, errorx.Parser(errorx.NetworkError()))
					return
				}

				if resp.StatusCode > 400 && resp.StatusCode <= 500 {
					errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

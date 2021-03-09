package middlewarex

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
)

// KetoPolicy model
type KetoPolicy struct {
	ID          string   `json:"id"`
	Subjects    []string `json:"subjects"`
	Actions     []string `json:"actions"`
	Resources   []string `json:"resources"`
	Effect      string   `json:"effect"`
	Description string   `json:"description"`
}

// CheckSuperOrganisation checks weather organisation of user is super org or not
func CheckSuperOrganisation(app string, GetOrganisation func(ctx context.Context) (int, error)) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !viper.GetBool("create_super_organisation") {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			oID, err := GetOrganisation(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			superOrgID, err := GetSuperOrganisationID(app)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if oID != superOrgID {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// GetSuperOrganisationID get superorganisation id from keto policy
func GetSuperOrganisationID(app string) (int, error) {
	req, err := http.NewRequest("GET", viper.GetString("keto_url")+"/engines/acp/ory/regex/policies/app:"+app+":superorg", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusOK {
		var policy KetoPolicy
		err = json.NewDecoder(resp.Body).Decode(&policy)
		if err != nil {
			return 0, err
		}

		if len(policy.Subjects) != 0 {
			orgID, _ := strconv.Atoi(policy.Subjects[0])
			return orgID, nil
		}
	}
	return 0, errors.New("cannot get super organisation id")
}

package middlewarex

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

type KetoRelationTuple struct {
	Namespace string `json:"namespace"`
	Relation  string `json:"relation"`
	Object    string `json:"object"`
	Subject   string `json:"subject_id"`
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
	req, err := http.NewRequest("GET", viper.GetString("keto_url")+"/relation-tuples?namespace=superorganisation", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var relationTuple KetoRelationTuple
		err = json.NewDecoder(resp.Body).Decode(&relationTuple)
		if err != nil {
			return 0, err
		}

		superOrgID := strings.Split(relationTuple.Object, ":")
		superOrgIDInt, err := strconv.Atoi(superOrgID[len(superOrgID)-1])
		if err != nil {
			return 0, err
		}
		return superOrgIDInt, nil
	}
	return 0, errors.New("cannot get super organisation id")
}

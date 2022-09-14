package middlewarex

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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

type KetoResponse struct {
	Tuples []KetoRelationTuple `json:"relation_tuples"`
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
	requestURL, err := url.Parse(viper.GetString("keto_read_api_url"))
	if err != nil {
		return 0, err
	}

	requestURL.Path += "relation-tuples"

	// add Query Parameters
	params := url.Values{}
	params.Add("namespace", "superorganisation")
	params.Add("subject_id", app)
	requestURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", requestURL.String(), nil)
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
		var relationTuple KetoResponse
		err = json.NewDecoder(resp.Body).Decode(&relationTuple)
		if err != nil {
			return 0, err
		}
		if len(relationTuple.Tuples) == 0 {
			return 0, errors.New("super organisation doesn't exists")
		}
		superOrgID := strings.Split(relationTuple.Tuples[0].Object, ":")
		superOrgIDInt, err := strconv.Atoi(superOrgID[len(superOrgID)-1])
		if err != nil {
			return 0, err
		}
		return superOrgIDInt, nil
	}
	return 0, errors.New("cannot get super organisation id")
}

package requestx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// AllApplicationUsers gets all the users for a application from kavach
func AllApplicationUsers(appSlug string, oID, uID uint) (*http.Response, error) {
	if appSlug == "" {
		return nil, errors.New("appSlug parameter cannot be empty")
	}

	resp, err := Request("GET", viper.GetString("kavach_url")+"/users/application?application="+appSlug, nil, map[string]string{
		"X-User":         fmt.Sprint(uID),
		"X-Organisation": fmt.Sprint(oID),
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 500 {
		return nil, errors.New("something went wrong")
	}

	return resp, nil
}

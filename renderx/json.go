package renderx

import (
	"encoding/json"
	"net/http"

	"github.com/opentracing/opentracing-go/log"
)

// JSON - render json
func JSON(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error(err)
	}
}

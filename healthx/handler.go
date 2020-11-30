package healthx

import (
	"net/http"

	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

const (
	// AliveCheckPath is the path where information about the life state of the instance is provided.
	AliveCheckPath = "/health/alive"
	// ReadyCheckPath is the path where information about the ready state of the instance is provided.
	ReadyCheckPath = "/health/ready"
)

// ReadyChecker should return an error if the component is not ready yet.
type ReadyChecker func() error

// ReadyCheckers is a map of ReadyCheckers.
type ReadyCheckers map[string]ReadyChecker

// notReadyErrors response if server is not ready
type notReadyErrors struct {
	Errors map[string]string `json:"errors"`
}

// RegisterRoutes registers health check routes
func RegisterRoutes(r *chi.Mux, readyCheckers ReadyCheckers) {
	r.Get(AliveCheckPath, Alive)
	r.Get(ReadyCheckPath, Ready(readyCheckers))
}

// Alive handler
func Alive(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "ok",
	}
	renderx.JSON(w, http.StatusOK, res)
}

// Ready handler
func Ready(readyCheckers ReadyCheckers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var notReady = notReadyErrors{
			Errors: map[string]string{},
		}

		for n, f := range readyCheckers {
			if err := f(); err != nil {
				notReady.Errors[n] = err.Error()
			}
		}
		if len(notReady.Errors) > 0 {
			renderx.JSON(w, http.StatusServiceUnavailable, notReady)
			return
		}

		renderx.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
	}
}

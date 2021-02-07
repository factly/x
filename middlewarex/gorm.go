package middlewarex

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

// GormRequestID returns middleware to add request_id in gorm context
func GormRequestID(db **gorm.DB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			*db = (*db).WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, requestID))
			h.ServeHTTP(w, r)
		})
	}
}

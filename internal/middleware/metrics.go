package middleware

import (
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/services"
)

func MetricsMiddleware(ms *services.MetricsService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ms.Increment()
			next.ServeHTTP(w, r)
		})
	}
}

package handlers

import (
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/config"
)

func ReadinessHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

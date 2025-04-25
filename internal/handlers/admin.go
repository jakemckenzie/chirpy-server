package handlers

import (
	"fmt"
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/config"
)

func AdminMetricsHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hits := cfg.MetricsService.GetHits()
		htmlTemplate := `
<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
</html>`
		htmlContent := fmt.Sprintf(htmlTemplate, hits)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	}
}

func ResetHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Platform != "dev" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		err := cfg.DBQueries.DeleteAllUsers(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

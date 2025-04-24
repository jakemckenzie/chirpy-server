package handlers

import (
	"fmt"
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/services"
)

func AdminMetricsHandler(ms *services.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hits := ms.GetHits()
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

func ResetHandler(ms *services.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms.Reset()
		w.WriteHeader(http.StatusOK)
	}
}

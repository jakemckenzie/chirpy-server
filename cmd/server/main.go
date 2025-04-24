package main

import (
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/handlers"
	"github.com/jakemckenzie/chirpy-server/internal/middleware"
	"github.com/jakemckenzie/chirpy-server/internal/services"
)

func main() {
	ms := services.NewMetricsService()
	ts := services.NewTextService()

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("../../static"))
	wrappedFileServer := middleware.MetricsMiddleware(ms)(http.StripPrefix("/app", fileServerHandler))
	mux.Handle("/app/", wrappedFileServer)

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/api/healthz", handlers.ReadinessHandler)
	mux.HandleFunc("/admin/metrics", handlers.AdminMetricsHandler(ms))
	mux.HandleFunc("/admin/reset", handlers.ResetHandler(ms))
	mux.HandleFunc("/api/validate_chirp", handlers.ValidateChirpHandler(ts))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

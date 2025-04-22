package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) adminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHits.Load()
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

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	// mux.HandleFunc("/healthz", readinessHandler)

	fileServerHandler := http.FileServer(http.Dir("."))
	wrappedFileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServerHandler))
	mux.Handle("/app/", wrappedFileServer)

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/api/healthz", apiCfg.readinessHandler)
	mux.HandleFunc("/admin/metrics", apiCfg.adminMetricsHandler)
	mux.HandleFunc("/admin/reset", apiCfg.resetHandler)

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
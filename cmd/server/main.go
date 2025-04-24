package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/jakemckenzie/chirpy-server/internal/config"
	"github.com/jakemckenzie/chirpy-server/internal/database"
	"github.com/jakemckenzie/chirpy-server/internal/handlers"
	"github.com/jakemckenzie/chirpy-server/internal/middleware"
	"github.com/jakemckenzie/chirpy-server/internal/services"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in .env")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	ms := services.NewMetricsService()
	ts := services.NewTextService()
	dbQueries := database.New(dbConn)

	cfg := &config.APIConfig{
		MetricsService: ms,
		TextService:    ts,
		DBQueries:      dbQueries,
	}

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("../../static"))
	wrappedFileServer := middleware.MetricsMiddleware(ms)(http.StripPrefix("/app", fileServerHandler))
	mux.Handle("/app/", wrappedFileServer)

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/api/healthz", handlers.ReadinessHandler(cfg))
	mux.HandleFunc("/admin/metrics", handlers.AdminMetricsHandler(cfg))
	mux.HandleFunc("/admin/reset", handlers.ResetHandler(cfg))
	mux.HandleFunc("/api/validate_chirp", handlers.ValidateChirpHandler(cfg))
	mux.HandleFunc("/api/users", handlers.CreateUserHandler(cfg))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on :8080...")
	server.ListenAndServe()
}

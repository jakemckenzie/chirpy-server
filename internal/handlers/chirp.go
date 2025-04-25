package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jakemckenzie/chirpy-server/internal/config"
	"github.com/jakemckenzie/chirpy-server/internal/database"
	"github.com/jakemckenzie/chirpy-server/internal/utils"
)

type ChirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func ChirpsHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Entered ChirpsHandler - Method: %s, Path: %s", r.Method, r.URL.Path)
		path := r.URL.Path
		if path == "/api/chirps" || path == "/api/chirps/" {
			switch r.Method {
			case http.MethodGet:
				log.Println("Calling handleGetAllChirps")
				handleGetAllChirps(w, r, cfg)
			case http.MethodPost:
				log.Println("Calling handleCreateChirp")
				handleCreateChirp(w, r, cfg)
			default:
				log.Println("Method not allowed:", r.Method)
				utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}

		if strings.HasPrefix(path, "/api/chirps/") && len(path) > len("/api/chirps/") {
			if r.Method == http.MethodGet {
				log.Printf("Routing to handleGetChirpByID for path: %s", path)
				handleGetChirpByID(w, r, cfg)
				return
			}
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		utils.RespondWithError(w, http.StatusNotFound, "Not found")
	}
}

func handleGetChirpByID(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) != 4 || pathParts[1] != "api" || pathParts[2] != "chirps" {
		log.Printf("Invalid path format: %v", pathParts)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirpIDStr := pathParts[3]
	log.Printf("Extracted chirpID: %s", chirpIDStr)
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID format")
		return
	}
	
	dbChirp, err := cfg.DBQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
		} else {
			log.Printf("Error retrieving chirp: %s", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
		}
		return
	}

	responseChirp := ChirpResponse{
		ID:        dbChirp.ID.String(),
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseChirp)
}

func handleGetAllChirps(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	dbChirps, err := cfg.DBQueries.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	responseChirps := make([]ChirpResponse, len(dbChirps))
	for i, chirp := range dbChirps {
		responseChirps[i] = ChirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseChirps)
}

func handleCreateChirp(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {

	type chirpRequest struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var req chirpRequest
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cfg.TextService.CleanProfanity(req.Body)

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	params := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}

	dbChirp, err := cfg.DBQueries.CreateChirp(r.Context(), params)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	responseChirp := ChirpResponse{
		ID:        dbChirp.ID.String(),
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID.String(),
	}

	utils.RespondWithJSON(w, http.StatusCreated, responseChirp)
}

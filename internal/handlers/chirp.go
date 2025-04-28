package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jakemckenzie/chirpy-server/internal/config"
	"github.com/jakemckenzie/chirpy-server/internal/database"
	"github.com/jakemckenzie/chirpy-server/internal/utils"
)

type ChirpID string
type UserID string
type ChirpBody string

type ChirpCollection []ChirpResponse

type ChirpResponse struct {
	ID        ChirpID   `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      ChirpBody `json:"body"`
	UserID    UserID    `json:"user_id"`
}

func ChirpsHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Entered ChirpsHandler - Method: %s, Path: %s", r.Method, r.URL.Path)
		path := r.URL.Path
		switch r.Method {
		case http.MethodGet:
			if path == "/api/chirps" || path == "/api/chirps/" {
				log.Println("Calling handleGetAllChirps")
				handleGetAllChirps(w, r, cfg)
				return
			}
			if strings.HasPrefix(path, "/api/chirps/") && len(path) > len("/api/chirps/") {
				log.Printf("Routing to handleGetChirpByID for path: %s", path)
				handleGetChirpByID(w, r, cfg)
				return
			}
			utils.RespondWithError(w, http.StatusNotFound, "Not found")
		case http.MethodPost:
			if path == "/api/chirps" || path == "/api/chirps/" {
				log.Println("Calling handleCreateChirp")
				handleCreateChirp(w, r, cfg)
				return
			}
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		default:
			log.Println("Method not allowed:", r.Method)
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

func parseChirpIDFromPath(path string) (uuid.UUID, error) {
	trimmedPath := strings.TrimSuffix(path, "/")
	pathParts := strings.Split(trimmedPath, "/")
	if len(pathParts) != 4 || pathParts[1] != "api" || pathParts[2] != "chirps" {
		return uuid.UUID{}, errors.New("invalid path format")
	}
	chirpIDStr := pathParts[3]
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return chirpID, nil
}

func handleGetChirpByID(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	chirpID, err := parseChirpIDFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Invalid chirp ID from path: %s", r.URL.Path)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirpByID(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	if err != nil {
		log.Printf("Error retrieving chirp: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
		return
	}

	responseChirp := ChirpResponse{
		ID:        ChirpID(dbChirp.ID.String()),
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      ChirpBody(dbChirp.Body),
		UserID:    UserID(dbChirp.UserID.String()),
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(responseChirp)
}

func handleGetAllChirps(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	dbChirps, err := cfg.DBQueries.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	var collection ChirpCollection
	for _, chirp := range dbChirps {
		responseChirp := ChirpResponse{
			ID:        ChirpID(chirp.ID.String()),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      ChirpBody(chirp.Body),
			UserID:    UserID(chirp.UserID.String()),
		}
		collection = append(collection, responseChirp)
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(collection)
}

type chirpRequest struct {
	Body   ChirpBody `json:"body"`
	UserID UserID    `json:"user_id"`
}

func handleCreateChirp(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
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

	cleanedBody := cfg.TextService.CleanProfanity(string(req.Body))

	userID, err := uuid.Parse(string(req.UserID))
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
		ID:        ChirpID(dbChirp.ID.String()),
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      ChirpBody(dbChirp.Body),
		UserID:    UserID(dbChirp.UserID.String()),
	}

	utils.RespondWithJSON(w, http.StatusCreated, responseChirp)
}

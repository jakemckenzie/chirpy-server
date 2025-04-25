package handlers

import (
	"time"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/config"
	"github.com/jakemckenzie/chirpy-server/internal/utils"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func CreateUserHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		type userRequest struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		var req userRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Email == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Email is required")
			return
		}

		dbUser, err := cfg.DBQueries.CreateUser(r.Context(), req.Email)
		if err != nil {
			log.Printf("Error creating user: %s", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		responseUser := User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		}

		utils.RespondWithJSON(w, http.StatusCreated, responseUser)
	}
}

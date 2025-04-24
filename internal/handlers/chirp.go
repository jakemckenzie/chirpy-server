package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jakemckenzie/chirpy-server/internal/services"
	"github.com/jakemckenzie/chirpy-server/internal/utils"
)

func ValidateChirpHandler(ts *services.TextService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		type validateChirpRequest struct {
			Body string `json:"body"`
		}
		type cleanedResponse struct {
			CleanedBody string `json:"cleaned_body"`
		}

		decoder := json.NewDecoder(r.Body)
		var req validateChirpRequest
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

		cleaned := ts.CleanProfanity(req.Body)
		utils.RespondWithJSON(w, http.StatusOK, cleanedResponse{CleanedBody: cleaned})
	}
}

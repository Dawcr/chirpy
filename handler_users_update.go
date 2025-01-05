package main

import (
	"encoding/json"
	"net/http"

	"github.com/dawcr/chirpy/internal/auth"
	"github.com/dawcr/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided in header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate JWT", err)
		return
	}

	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	param := request{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	dbResponse, err := cfg.db.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		ID:             userID,
		Email:          param.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update email/password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Email:     dbResponse.Email,
		},
	})
}

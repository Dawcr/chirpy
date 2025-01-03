package main

import (
	"encoding/json"
	"net/http"

	"github.com/dawcr/chirpy/internal/database/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	param := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	dbResponse, err := cfg.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve password", err)
		return
	}

	if auth.CheckPasswordHash(param.Password, dbResponse.HashedPassword) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbResponse.ID,
		CreatedAt: dbResponse.CreatedAt,
		UpdatedAt: dbResponse.UpdatedAt,
		Email:     dbResponse.Email,
	})
}

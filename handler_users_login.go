package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dawcr/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
	}

	param := loginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	dbResponse, err := cfg.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckPasswordHash(param.Password, dbResponse.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expirationTime time.Duration
	if param.ExpiresInSeconds == nil || *param.ExpiresInSeconds > 3600 || *param.ExpiresInSeconds <= 0 { // if not set, longer than 1 hour or negative
		expirationTime = time.Hour
	} else {
		expirationTime = time.Duration(*param.ExpiresInSeconds) * time.Second
	}

	userJWT, err := auth.MakeJWT(dbResponse.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create JWT for user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Email:     dbResponse.Email,
			Token:     userJWT,
		},
	})
}

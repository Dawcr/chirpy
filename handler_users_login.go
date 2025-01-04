package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dawcr/chirpy/internal/auth"
	"github.com/dawcr/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(dbResponse.ID, cfg.jwtSecret, accessTokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbResponse.ID,
		ExpiresAt: time.Now().UTC().Add(60 * time.Hour * 24),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Email:     dbResponse.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

package main

import (
	"net/http"
	"time"

	"github.com/dawcr/chirpy/internal/auth"
)

const (
	accessTokenDuration = time.Hour
)

func (cfg *apiConfig) handlerTokensRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to find refresh token", err)
		return
	}

	dbResponse, err := cfg.db.GetUserByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}
	if time.Now().UTC().Compare(dbResponse.ExpiresAt) > 0 {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbResponse.UserID, cfg.jwtSecret, accessTokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerTokensRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to find refresh token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

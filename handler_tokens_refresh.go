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
		respondWithError(w, http.StatusBadRequest, "No token provided in header", err)
		return
	}

	dbResponse, err := cfg.db.GetUserByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to fetch user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbResponse.ID, cfg.jwtSecret, accessTokenDuration)
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
		respondWithError(w, http.StatusBadRequest, "No token provided in header", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) webhookUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		User_id string `json:"user_id"`
	}

	type parameters struct {
		Event string      `json:"event"`
		Data  requestData `json:"data"`
	}

	param := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params", err)
	}

	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(param.Data.User_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse user ID", err)
		return
	}

	if _, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

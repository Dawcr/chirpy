package main

import (
	"net/http"
	"slices"

	"github.com/dawcr/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")

	var dbResponse []database.Chirp

	if s == "" {
		var err error
		dbResponse, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
			return
		}
	} else {
		uid, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Invalid author ID", err)
			return
		}
		dbResponse, err = cfg.db.GetChirpsFromUser(r.Context(), uid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
			return
		}
	}

	response := []Chirp{}
	for _, chirp := range dbResponse {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if order == "desc" {
		slices.Reverse(response)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerChirpsGetSingle(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse chirp ID", err)
		return
	}

	dbResponse, err := cfg.db.GetChirpsSingle(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbResponse.ID,
		CreatedAt: dbResponse.CreatedAt,
		UpdatedAt: dbResponse.UpdatedAt,
		Body:      dbResponse.Body,
		UserID:    dbResponse.UserID,
	})
}

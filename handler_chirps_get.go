package main

import "net/http"

func (cfg *apiConfig) handlerchirpsGet(w http.ResponseWriter, r *http.Request) {

	dbResponse, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
		return
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

	respondWithJSON(w, http.StatusOK, response)
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dawcr/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	cleaned_body, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbResponse, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body, // use the version without profanities
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Body:      dbResponse.Body,
			UserID:    dbResponse.UserID,
		},
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	banned_words := []string{"kerfuffle", "sharbert", "fornax"}

	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned_body := replaceProfane(body, mapBanned(banned_words...))

	return cleaned_body, nil
}

func replaceProfane(chirp string, banned map[string]struct{}) string {
	segments := strings.Split(chirp, " ")
	for i, segment := range segments {
		if _, ok := banned[strings.ToLower(segment)]; ok {
			segments[i] = "****"
		}
	}
	return strings.Join(segments, " ")
}

func mapBanned(words ...string) map[string]struct{} {
	dict := make(map[string]struct{}, len(words))
	for _, word := range words {
		dict[strings.ToLower(word)] = struct{}{}
	}
	return dict
}

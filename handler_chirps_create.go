package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dawcr/chirpy/internal/database"
	"github.com/google/uuid"
)

const (
	maxChirpLength = 140
)

type Chirp_message struct {
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
		Chirp_message
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	banned_words := mapBanned("kerfuffle", "sharbert", "fornax")
	cleaned_body := replaceProfane(params.Body, banned_words)

	dbResponse, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body, // use the version without profanities
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp_message: Chirp_message{
			ID:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Body:      dbResponse.Body,
			UserID:    dbResponse.UserID,
		},
	})
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

package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	maxChirpLength = 140
)

func handlerChirpsValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned_body,
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

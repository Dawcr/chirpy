package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	cfg.db.DeleteUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All user data deleted"))
}

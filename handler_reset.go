package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in dev environment.", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and user database cleared"))
}

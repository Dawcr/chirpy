package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/dawcr/chirpy/internal/database"
)

const (
	port                  = "8080"
	filepathRoot          = "."
	filepathSite          = "/app/"
	filepathHealthz       = "/api/healthz"
	filepathMetricz       = "/admin/metrics"
	filepathResetMetricz  = "/admin/reset"
	filepathValidateChirp = "/api/validate_chirp"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func startServer() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL missing")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to db: %s", err)
	}

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
	}

	mux.Handle(filepathSite, apiCfg.middlewareMetricsInc(http.StripPrefix(filepathSite, http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET "+filepathHealthz, handlerReadiness)
	mux.HandleFunc("GET "+filepathMetricz, apiCfg.handlerHitCount)
	mux.HandleFunc("POST "+filepathResetMetricz, apiCfg.handlerResetHitCount)
	mux.HandleFunc("POST "+filepathValidateChirp, handlerChirpsValidation)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

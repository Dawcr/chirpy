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
	filepathReset         = "/admin/reset"
	filepathValidateChirp = "/api/validate_chirp"
	filepathCreateUser    = "/api/users"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
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
		platform:       os.Getenv("PLATFORM"),
	}

	mux.Handle(filepathSite, apiCfg.middlewareMetricsInc(http.StripPrefix(filepathSite, http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET "+filepathHealthz, handlerReadiness)
	mux.HandleFunc("GET "+filepathMetricz, apiCfg.handlerHitCount)
	mux.HandleFunc("POST "+filepathReset, apiCfg.handlerReset)
	mux.HandleFunc("POST "+filepathValidateChirp, handlerChirpsValidation)
	mux.HandleFunc("POST "+filepathCreateUser, apiCfg.handlerCreateUser)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

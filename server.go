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
	port            = "8080"
	path_Root       = "."
	path_App        = "/app/"
	path_Healthz    = "/api/healthz"
	path_Metrics    = "/admin/metrics"
	path_Reset      = "/admin/reset"
	path_Chirps     = "/api/chirps"
	path_CreateUser = "/api/users"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func startServer() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to db: %s", err)
	}

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
	}

	mux.Handle(path_App, apiCfg.middlewareMetricsInc(http.StripPrefix(path_App, http.FileServer(http.Dir(path_Root)))))
	mux.HandleFunc("GET "+path_Healthz, handlerReadiness)
	mux.HandleFunc("GET "+path_Metrics, apiCfg.handlerHitCount)
	mux.HandleFunc("POST "+path_Reset, apiCfg.handlerReset)
	mux.HandleFunc("POST "+path_Chirps, apiCfg.handlerChirpsValidation)
	mux.HandleFunc("POST "+path_CreateUser, apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET "+path_Chirps, apiCfg.handlerChirpsGet)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

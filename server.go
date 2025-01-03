package main

import (
	"log"
	"net/http"
	"sync/atomic"
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
}

func startServer() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
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

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

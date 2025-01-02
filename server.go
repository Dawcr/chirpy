package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

const (
	port                 = "8080"
	filepathRoot         = "."
	filepathAPP          = "/app/"
	filepathHealthz      = "/healthz"
	filepathMetricz      = "/metrics"
	filepathResetMetricz = "/reset"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func startServer() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle(filepathAPP, apiCfg.middlewareMetricsInc(http.StripPrefix(filepathAPP, http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc(filepathHealthz, handlerReadiness)
	mux.HandleFunc(filepathMetricz, apiCfg.handlerHitCount)
	mux.HandleFunc(filepathResetMetricz, apiCfg.handlerResetHitCount)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

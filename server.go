package main

import "net/http"

const (
	port = "8080"
)

func newServer() *http.Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}
	return server
}

func startServer(server *http.Server) {
	server.ListenAndServe()
}

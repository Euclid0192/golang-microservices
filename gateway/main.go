package main

import (
	"log"
	"net/http"
)

const (
	httpAddr = ":8080"
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting http server on port: %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}

}

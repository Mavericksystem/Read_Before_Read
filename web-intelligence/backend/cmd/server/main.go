package main

import (
	"log"
	"net/http"
	"os"
	"web-intelligence/backend/internal/handler"
	"web-intelligence/backend/internal/nim"
)

func main() {
	nimClient, err := nim.NewClient()
	if err != nil {
		log.Fatalf("Failed to create NIM client: %v", err)
	}

	analyze := &handler.Analyzer{NimClient: nimClient}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/analyze", analyzer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

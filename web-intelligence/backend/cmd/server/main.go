package main

import (
	"log"
	"net/http"
	"os"
	"web-intelligence/backend/internal/handler"
	"web-intelligence/backend/internal/nim"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	nimClient, err := nim.NewClient()
	if err != nil {
		log.Fatalf("Failed to create NIM client: %v", err)
	}

	analyze := &handler.Analyzer{NimClient: nimClient}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/analyze", analyze)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

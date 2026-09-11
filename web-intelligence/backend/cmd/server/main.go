package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"web-intelligence/backend/internal/browser"
	"web-intelligence/backend/internal/extractor"
	"web-intelligence/backend/internal/handler"
	"web-intelligence/backend/internal/nim"
)

// Pool sizes are placeholders until Open Questions ("Rust worker pool
// size", "headless Chrome pool size", "NIM call throttling") are resolved
// via benchmarking — do not treat these as final, tune with real numbers.
const (
	extractorPoolSize = 8
	browserPoolSize   = 2
	nimMaxInFlight    = 2
	admissionCapacity = 50 // max requests in flight before 503
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
	throttledNim := nim.NewThrottled(nimClient, nimMaxInFlight)

	extractorPool, err := extractor.NewPooled(extractorPoolSize)
	if err != nil {
		log.Fatalf("Failed to start extractor pool: %v", err)
	}
	defer extractorPool.Shutdown()

	browserPool, err := browser.New(browserPoolSize)
	if err != nil {
		log.Fatalf("Failed to start browser pool: %v", err)
	}
	defer browserPool.Shutdown()

	analyze := &handler.Analyzer{
		NimClient: throttledNim,
		Extractor: extractorPool,
		Browser:   browserPool,
		Admission: make(chan struct{}, admissionCapacity),
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/analyze", analyze)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: withCORS(mux),
	}

	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for Ctrl+C / kill so the deferred pool Shutdown() calls above
	// actually run — without this, killing the process leaks Rust
	// subprocesses and the Chrome instance instead of cleaning them up.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("shutting down...")
}

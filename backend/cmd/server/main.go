package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/teamcastaldi/rc-car-cam/backend/internal/camera"
	"github.com/teamcastaldi/rc-car-cam/backend/internal/stream"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mockSource, err := camera.NewMockSource(100 * time.Millisecond)
	if err != nil {
		log.Fatalf("create mock camera source: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.Handle("GET /stream", stream.NewHandler(mockSource))

	addr := ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		// ReadHeaderTimeout mitigates slowloris-style connections. ReadTimeout and
		// WriteTimeout are intentionally left unset: once this server hosts the
		// video stream, either would cap the duration of a long-lived response.
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("rc-car-cam backend listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("healthz: write response: %v", err)
	}
}

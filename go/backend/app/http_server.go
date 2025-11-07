package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *App) startHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", sseHandler)
	mux.HandleFunc("/health", a.healthHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("HTTP server is running on :8080")

		a.httpReady.Store(true)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("app startHTTPServer srv.ListenAndServe error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("HTTP server shutting down")
		srv.Shutdown(ctx)
		return nil
	case err := <-errCh:
		return fmt.Errorf("app startHTTPServer error: %w", err)
	}
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	if a.httpReady.Load() && a.mqttReady.Load() {
		log.Println("Health check OK")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	log.Println("Health check NOT OK")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("Service Unavailable"))
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Client connected to SSE")

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for client disconnection
	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-clientGone:
			log.Println("Client disconnected")
			return
		case <-t.C:
			// Send an event to the client
			// Here we send only the "data" field, but there are few others
			_, err := fmt.Fprintf(w, "data: The time is %s\n\n", time.Now().Format(time.UnixDate))
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
	}
}

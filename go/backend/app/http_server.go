package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/max-pv/fourier/go-shared/models"
)

const (
	typeInitialData = "initial_data"
	typeUpdateData  = "update_data"
)

type SSEPayload struct {
	Kind    string              `json:"kind"`
	Payload []*models.DataPoint `json:"payload"`
}

func (a *App) startHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", a.sseHandler)
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
		// log.Println("Health check OK")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	log.Println("Health check NOT OK")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("Service Unavailable"))
}

func (a *App) sseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Client connected to SSE")

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for this client
	clientChan := make(chan *models.DataPoint)
	a.mu.Lock()
	a.clients[clientChan] = struct{}{}
	a.mu.Unlock()

	// Remove the client when they disconnect
	defer func() {
		a.mu.Lock()
		delete(a.clients, clientChan)
		a.mu.Unlock()
		close(clientChan)
		log.Println("Client disconnected")
	}()

	// Parse query parameters
	query := r.URL.Query()
	dataType := query.Get("type") // Get the "type" parameter
	if dataType == "" {
		dataType = typeEverything
	}

	log.Printf("Fetching initial data with type: %s, start: %s, end: %s", dataType, time.Now().Add(-1*time.Hour), time.Now())

	// Create a channel for client disconnection
	initialData, err := a.db.GetByTypeAndTimeRange(r.Context(), dataType, time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		log.Printf("sseHandler GetByTypeAndTimeRange error: %v", err)
	}

	log.Printf("len %d", len(initialData))

	rc := http.NewResponseController(w)

	if len(initialData) > 0 {
		payload := SSEPayload{
			Kind:    typeInitialData,
			Payload: initialData,
		}
		data, err := json.Marshal(payload)
		if err == nil {
			fmt.Fprintf(w, "data: %s\n\n", data)
			if err := rc.Flush(); err != nil {
				log.Printf("sseHandler rc.Flush error: %v", err)
				return
			}
		}
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-clientChan:
			payload := SSEPayload{
				Kind:    typeUpdateData,
				Payload: []*models.DataPoint{msg},
			}
			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("sseHandler json.Marshal error: %v", err)
				return
			}

			_, err = fmt.Fprintf(w, "%s\n\n", data)
			if err != nil {
				return
			}

			flusher, ok := w.(http.Flusher)
			if ok {
				flusher.Flush()
			}
		}
	}
}

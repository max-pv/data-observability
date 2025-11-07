package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/max-pv/fourier/go-shared/models"
)

const (
	topicTelemetry = "telemetry"
)

type Storage interface {
	InsertDataPoint(ctx context.Context, dp *models.DataPoint) error
	GetByTypeAndTimeRange(ctx context.Context, dataType string, start, end time.Time) ([]*models.DataPoint, error)
}

type App struct {
	httpReady atomic.Bool
	mqttReady atomic.Bool

	db         Storage
	sseClients map[chan *models.DataPoint]struct{}
	mu         sync.Mutex
}

func New() *App {
	return &App{
		sseClients: make(map[chan *models.DataPoint]struct{}),
	}
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 2)

	// this is blocking because we need the DB connection for MQTT hooks
	if err := a.connectToDatabase(ctx); err != nil {
		return fmt.Errorf("app Run connectToDatabase error: %w", err)
	}

	go func() {
		errChan <- a.startMQTTServer(ctx)
	}()

	go func() {
		errChan <- a.startHTTPServer(ctx)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("app Run error: %w", err)
		}
	case <-ctx.Done():
		log.Println("App shutting down")
		return ctx.Err()
	}

	return nil
}

func (a *App) Broadcast(dataPoint *models.DataPoint) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for clientChan := range a.sseClients {
		select {
		case clientChan <- dataPoint:
		default:
			// If the client is not ready to receive, skip it
			log.Println("Skipping client due to slow connection")
		}
	}
}

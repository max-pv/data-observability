package app

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/max-pv/fourier/go-shared/models"
)

const (
	topicTelemetry = "telemetry"
)

type Storage interface {
	InsertDataPoint(ctx context.Context, dp *models.DataPoint) error
}

type App struct {
	httpReady atomic.Bool
	mqttReady atomic.Bool

	db Storage
}

func New() *App {
	return &App{}
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

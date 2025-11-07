package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/max-pv/fourier/producer/app"
)

func main() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigs
		cancel()
	}()

	app, err := app.New()
	if err != nil {
		log.Fatalf("app.New returned error: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		if err != context.Canceled {
			log.Printf("app.Run returned error: %v", err)
		}
	}

	<-ctx.Done()
	log.Println("producer stopped")
}

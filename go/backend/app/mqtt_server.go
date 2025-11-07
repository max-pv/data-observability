package app

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/max-pv/data-observability/go-shared/models"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	mqtt "github.com/mochi-mqtt/server/v2"
	hooks "github.com/mochi-mqtt/server/v2/hooks/auth"
)

type listenerHook struct {
	mqtt.HookBase

	app *App
}

func (a *App) startMQTTServer(ctx context.Context) error {
	s := mqtt.New(nil)

	s.AddHook(&hooks.AllowHook{}, nil)
	s.AddHook(&listenerHook{app: a}, nil)

	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})

	err := s.AddListener(tcp)
	if err != nil {
		return fmt.Errorf("app startMQTTServer s.AddListener error: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("MQTT server starting on :1883")

		a.mqttReady.Store(true)
		if err := s.Serve(); err != nil {
			errCh <- fmt.Errorf("app startMQTTServer s.Serve error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("MQTT server shutting down")
		s.Close()
		return nil
	case err := <-errCh:
		return err
	}
}

// Provides indicates which hook methods this hook provides.
func (l *listenerHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnPublish,
	}, []byte{b})
}

func (l *listenerHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	if pk.TopicName != topicTelemetry {
		return pk, nil
	}

	dp, err := models.DataPointFromJSON(string(pk.Payload))
	if err != nil {
		log.Printf("app OnPublish DataPointFromJSON error: %v", err)
		return pk, nil
	}

	// Broadcast to connected SSE clients
	l.app.Broadcast(dp)

	// Run DB insertion in a separate goroutine to avoid blocking the MQTT server
	// not going to handle insertion errors here, just log them - in the real world dropping events should not trigger system failures
	go func() {
		err = l.app.db.InsertDataPoint(context.Background(), dp)
		if err != nil {
			log.Printf("app OnPublish InsertDataPoint error: %v", err)
		}
		log.Printf(`Inserted "%s" datapoint`, dp.Type)
	}()

	return pk, nil
}

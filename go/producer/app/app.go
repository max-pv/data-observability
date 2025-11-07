package app

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/max-pv/fourier/go-shared/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	PowerInputType             = "PowerInput"
	WaterFlowRateType          = "WaterFlowRate"
	TemperatureType            = "Temperature"
	HydrogenProductionRateType = "HydrogenProductionRate"
	EfficiencyType             = "Efficiency"
)

const (
	topicTelemetry = "telemetry"
)

type App struct {
	mqttClient mqtt.Client
}

func New() (*App, error) {
	mqttClient, err := createMQTTClient()
	if err != nil {
		return nil, fmt.Errorf("app New createMQTTClient error: %w", err)
	}

	return &App{
		mqttClient: *mqttClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	dataChan := make(chan *models.DataPoint, 100)

	go producePowerInput(ctx, dataChan)
	go produceWaterFlowRate(ctx, dataChan)
	go produceTemperature(ctx, dataChan)
	go produceHydrogenProductionRate(ctx, dataChan)
	go produceEfficiency(ctx, dataChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(dataChan)
				return
			case data := <-dataChan:
				json := data.ToJSON()
				token := a.mqttClient.Publish(topicTelemetry, 0, true, json)
				token.Wait()
				if token.Error() != nil {
					log.Printf("app Run publish error: %v", token.Error())
				} else {
					// log.Printf("Published: %s", json)
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func createMQTTClient() (*mqtt.Client, error) {
	brokerURL := os.Getenv("MQTT_BROKER")
	if brokerURL == "" {
		brokerURL = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL) // Replace with your MQTT broker address
	opts.SetClientID("telemetry-producer")

	// Create and connect the MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	log.Println("Connected to MQTT broker")
	return &client, nil
}

// Random interval between 1.5s and 4s
func newRandomTicker() *time.Ticker {
	return time.NewTicker(time.Millisecond * time.Duration(1500+rand.Float64()*2000))
}

func producePowerInput(ctx context.Context, ch chan<- *models.DataPoint) {
	ticker := newRandomTicker()
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v := 1000 + 500*math.Sin(float64(i)*0.1)
			ch <- models.NewDataPoint(v, PowerInputType)
		}
	}
}

func produceWaterFlowRate(ctx context.Context, ch chan<- *models.DataPoint) {
	ticker := newRandomTicker()
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v := 50 + rand.Float64()*10
			ch <- models.NewDataPoint(v, WaterFlowRateType)
		}
	}
}

func produceTemperature(ctx context.Context, ch chan<- *models.DataPoint) {
	ticker := newRandomTicker()
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v := 25 + 5*math.Sin(float64(i)*0.2) + rand.Float64()
			ch <- models.NewDataPoint(v, TemperatureType)
		}
	}
}

func produceHydrogenProductionRate(ctx context.Context, ch chan<- *models.DataPoint) {
	ticker := newRandomTicker()
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v := 10 + float64(i)*0.1 + rand.Float64()
			ch <- models.NewDataPoint(v, HydrogenProductionRateType)
		}
	}
}

func produceEfficiency(ctx context.Context, ch chan<- *models.DataPoint) {
	ticker := newRandomTicker()
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v := 90 + rand.Float64()*5
			ch <- models.NewDataPoint(v, EfficiencyType)
		}
	}
}

package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/max-pv/data-observability/go-shared/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "data-observability"
	collectionName = "telemetry"
	typeEverything = "*"
)

type Database struct {
	client *mongo.Client
}

func NewDatabase(client *mongo.Client) *Database {
	return &Database{client: client}
}

func (db *Database) InsertDataPoint(ctx context.Context, dp *models.DataPoint) error {
	collection := db.client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(ctx, dp)
	if err != nil {
		return fmt.Errorf("Database InsertDataPoint InsertOne error: %w", err)
	}
	return nil
}

func (db *Database) GetByTypeAndTimeRange(ctx context.Context, dataType string, start, end time.Time) ([]*models.DataPoint, error) {
	collection := db.client.Database(dbName).Collection(collectionName)

	filter := map[string]interface{}{
		"timestamp": map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}

	if dataType != typeEverything && dataType != "" {
		filter["type"] = dataType
	}

	// opts := options.Find().SetLimit(60)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Database GetByTypeAndTimeRange Find error: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*models.DataPoint
	for cursor.Next(ctx) {
		var dp models.DataPoint
		if err := cursor.Decode(&dp); err != nil {
			return nil, fmt.Errorf("Database GetByTypeAndTimeRange Decode error: %w", err)
		}
		results = append(results, &dp)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("Database GetByTypeAndTimeRange cursor error: %w", err)
	}

	return results, nil
}

func (a *App) connectToDatabase(ctx context.Context) error {
	uri := "mongodb://root:example@mongo:27017"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("app connectToDatabase mongo.Connect error: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = client.Ping(ctxPing, nil); err != nil {
		return fmt.Errorf("app connectToDatabase client.Ping error: %w", err)
	}

	// a.db = client.Database(dbName)
	a.db = NewDatabase(client)
	log.Println("Connected to MongoDB")

	go func() {
		<-ctx.Done()
		log.Println("Disconnecting from MongoDB")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	return nil
}

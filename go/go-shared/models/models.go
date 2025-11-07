package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type DataPoint struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Value     float64   `json:"value" bson:"value"`
	Type      string    `json:"type" bson:"type"`
}

func NewDataPoint(value float64, dataType string) *DataPoint {
	return &DataPoint{
		Timestamp: time.Now().UTC(),
		Value:     value,
		Type:      dataType,
	}
}

func (dp *DataPoint) ToJSON() string {
	return fmt.Sprintf(
		`{"timestamp":"%s","value":%.2f,"type":"%s"}`,
		dp.Timestamp.Format(time.RFC3339),
		dp.Value,
		dp.Type,
	)
}

func DataPointFromJSON(jsonStr string) (*DataPoint, error) {
	dp := &DataPoint{}

	err := json.Unmarshal([]byte(jsonStr), dp)
	if err != nil {
		return nil, fmt.Errorf("models DataPointFromJSON error: %w", err)
	}

	return dp, nil
}

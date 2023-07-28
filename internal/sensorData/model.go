package sensordata

import (
	"sensors-generator/internal/spiece"
	"time"
)

type SensorData struct {
	ID              int
	SensorID        int
	Temperature     float32
	Transparency    uint8
	DetectedSpieces []spiece.Spiece
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateSensorDataDTO struct {
	SensorID     int     `json:"sensor_id"`
	Temperature  float32 `json:"temperature"`
	Transparency uint8   `json:"transparency"`
}

type SensorDataFilters struct {
}

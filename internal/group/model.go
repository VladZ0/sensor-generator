package group

import (
	"sensors-generator/internal/sensor"
	"time"
)

type SensorGroup struct {
	ID        int
	Name      string
	Sensors   []sensor.Sensor
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateSensorGroupDTO struct {
	Name string
}

type SensorGroupFilters struct {
	TopLimit int
	FromDate time.Time
	TillDate time.Time
}

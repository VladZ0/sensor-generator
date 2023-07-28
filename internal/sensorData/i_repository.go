package sensordata

import (
	"context"
	"sensors-generator/internal/spiece"
)

type ISensorDataRepository interface {
	FindAll(ctx context.Context, filters SensorDataFilters) ([]SensorData, error)
	FindOneByID(ctx context.Context, id int, filters SensorDataFilters) (*SensorData, error)
	Create(ctx context.Context, sensorData CreateSensorDataDTO) (int, error)
	AddDetectedSpiece(ctx context.Context, sensorDataID int, spiece spiece.Spiece) error
}

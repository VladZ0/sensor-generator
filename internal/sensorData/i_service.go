package sensordata

import (
	"context"
	"sensors-generator/internal/spiece"
)

type ISensorDataService interface {
	GetOneByID(ctx context.Context, id int, filters SensorDataFilters) (*SensorData, error)
	Create(ctx context.Context, sensorData ...CreateSensorDataDTO) ([]int, error)
	AddDetectedSpieces(ctx context.Context, sensorDataID int, spieces ...spiece.Spiece) error
}

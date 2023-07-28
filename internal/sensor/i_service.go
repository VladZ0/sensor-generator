package sensor

import "context"

type ISensorService interface {
	GetAll(ctx context.Context, filters SensorFilters) ([]Sensor, error)
	Create(ctx context.Context, sensors ...CreateSensorDTO) error
	AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error
	GetExtremumTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates, min bool) (float32, error)
	GetAvgTemperatureForSensor(ctx context.Context, filters SensorFilters) (float32, error)
}

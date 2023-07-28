package sensor

import "context"

type ISensorRepository interface {
	FindAll(ctx context.Context, filters SensorFilters) ([]Sensor, error)
	FindOneByID(ctx context.Context, id int, filters SensorFilters) (*Sensor, error)
	Create(ctx context.Context, sensor CreateSensorDTO) error
	AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error
	FindMaxTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates) (float32, error)
	FindMinTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates) (float32, error)
	FindAvgTemperatureForSensor(ctx context.Context, filters SensorFilters) (float32, error)
}

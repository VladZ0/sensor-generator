package sensor

import (
	"context"
	"sensors-generator/internal/sensor"

	"github.com/stretchr/testify/mock"
)

type MockSensorService struct {
	mock.Mock
}

func (m *MockSensorService) GetAll(ctx context.Context, filters sensor.SensorFilters) ([]sensor.Sensor, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]sensor.Sensor), args.Error(1)
}

func (m *MockSensorService) Create(ctx context.Context, sensors ...sensor.CreateSensorDTO) error {
	args := m.Called(ctx, sensors)
	return args.Error(0)
}

func (m *MockSensorService) AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error {
	args := m.Called(ctx, sensorID, groupID)
	return args.Error(0)
}

func (m *MockSensorService) GetExtremumTemperatureForRegion(ctx context.Context, minCoords, maxCoords sensor.Coordinates, min bool) (float32, error) {
	args := m.Called(ctx, minCoords, maxCoords, min)
	return float32(args.Get(0).(float32)), args.Error(1)
}

func (m *MockSensorService) GetAvgTemperatureForSensor(ctx context.Context, filters sensor.SensorFilters) (float32, error) {
	args := m.Called(ctx, filters)
	return float32(args.Get(0).(float32)), args.Error(1)
}

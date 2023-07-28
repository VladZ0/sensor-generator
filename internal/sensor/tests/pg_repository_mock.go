package sensor

import (
	"context"
	"sensors-generator/internal/sensor"

	"github.com/stretchr/testify/mock"
)

type MockSensorRepository struct {
	mock.Mock
}

func (m *MockSensorRepository) FindAll(ctx context.Context, filters sensor.SensorFilters) ([]sensor.Sensor, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]sensor.Sensor), args.Error(1)
}

func (m *MockSensorRepository) FindOneByID(ctx context.Context, id int, filters sensor.SensorFilters) (*sensor.Sensor, error) {
	args := m.Called(ctx, id, filters)
	if obj := args.Get(0); obj != nil {
		return obj.(*sensor.Sensor), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSensorRepository) Create(ctx context.Context, sensor sensor.CreateSensorDTO) error {
	args := m.Called(ctx, sensor)
	return args.Error(0)
}

func (m *MockSensorRepository) AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error {
	args := m.Called(ctx, sensorID, groupID)
	return args.Error(0)
}

func (m *MockSensorRepository) FindMaxTemperatureForRegion(ctx context.Context, minCoords, maxCoords sensor.Coordinates) (float32, error) {
	args := m.Called(ctx, minCoords, maxCoords)
	return float32(args.Get(0).(float32)), args.Error(1)
}

func (m *MockSensorRepository) FindMinTemperatureForRegion(ctx context.Context, minCoords, maxCoords sensor.Coordinates) (float32, error) {
	args := m.Called(ctx, minCoords, maxCoords)
	return float32(args.Get(0).(float32)), args.Error(1)
}

func (m *MockSensorRepository) FindAvgTemperatureForSensor(ctx context.Context, filters sensor.SensorFilters) (float32, error) {
	args := m.Called(ctx, filters)
	return float32(args.Get(0).(float32)), args.Error(1)
}

package sensordata

import (
	"context"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"

	"github.com/stretchr/testify/mock"
)

type MockSensorDataRepository struct {
	mock.Mock
}

func (m *MockSensorDataRepository) FindAll(ctx context.Context, filters sensordata.SensorDataFilters) ([]sensordata.SensorData, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]sensordata.SensorData), args.Error(1)
}

func (m *MockSensorDataRepository) FindOneByID(ctx context.Context, id int, filters sensordata.SensorDataFilters) (*sensordata.SensorData, error) {
	args := m.Called(ctx, id, filters)
	if obj := args.Get(0); obj != nil {
		return obj.(*sensordata.SensorData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSensorDataRepository) Create(ctx context.Context, sensorData sensordata.CreateSensorDataDTO) (int, error) {
	args := m.Called(ctx, sensorData)
	return args.Int(0), args.Error(1)
}

func (m *MockSensorDataRepository) AddDetectedSpiece(ctx context.Context, sensorDataID int, spiece spiece.Spiece) error {
	args := m.Called(ctx, sensorDataID, spiece)
	return args.Error(0)
}

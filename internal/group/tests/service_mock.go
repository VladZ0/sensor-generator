package group

import (
	"context"
	"sensors-generator/internal/group"
	"sensors-generator/internal/spiece"

	"github.com/stretchr/testify/mock"
)

type MockSensorGroupService struct {
	mock.Mock
}

func (m *MockSensorGroupService) GetSpiecesInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (map[*spiece.Spiece]int, error) {
	args := m.Called(ctx, groupName, filters)
	actualSpieces := args.Get(0).(map[spiece.Spiece]int)
	expectedSpieces := make(map[*spiece.Spiece]int)

	for key, value := range actualSpieces {
		tempKey := key
		expectedSpieces[&tempKey] = value
	}

	return expectedSpieces, args.Error(1)
}

func (m *MockSensorGroupService) GetAvgTrasparencyInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (uint8, error) {
	args := m.Called(ctx, groupName, filters)
	return args.Get(0).(uint8), args.Error(1)
}

func (m *MockSensorGroupService) GetAvgTemperatureInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (float32, error) {
	args := m.Called(ctx, groupName, filters)
	return args.Get(0).(float32), args.Error(1)
}

func (m *MockSensorGroupService) Create(ctx context.Context, groups ...group.CreateSensorGroupDTO) error {
	args := m.Called(ctx, groups)
	return args.Error(0)
}

func (m *MockSensorGroupService) GetAll(ctx context.Context, filters group.SensorGroupFilters) ([]group.SensorGroup, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]group.SensorGroup), args.Error(1)
}

package group

import (
	"context"
	"sensors-generator/internal/group"
	"sensors-generator/internal/spiece"

	"github.com/stretchr/testify/mock"
)

type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) FindSpiecesInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (map[*spiece.Spiece]int, error) {
	args := m.Called(ctx, groupName, filters)
	return args.Get(0).(map[*spiece.Spiece]int), args.Error(1)
}

func (m *MockGroupRepository) FindOneByID(ctx context.Context, id int, filters group.SensorGroupFilters) (*group.SensorGroup, error) {
	args := m.Called(ctx, id, filters)
	return args.Get(0).(*group.SensorGroup), args.Error(1)
}

func (m *MockGroupRepository) FindAll(ctx context.Context, filters group.SensorGroupFilters) ([]group.SensorGroup, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]group.SensorGroup), args.Error(1)
}

func (m *MockGroupRepository) FindAvgTransparencyInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (uint8, error) {
	args := m.Called(ctx, groupName, filters)
	return args.Get(0).(uint8), args.Error(1)
}

func (m *MockGroupRepository) FindAvgTemperatureInGroup(ctx context.Context, groupName string, filters group.SensorGroupFilters) (float32, error) {
	args := m.Called(ctx, groupName, filters)
	return args.Get(0).(float32), args.Error(1)
}

func (m *MockGroupRepository) Create(ctx context.Context, grp group.CreateSensorGroupDTO) error {
	args := m.Called(ctx, grp)
	return args.Error(0)
}

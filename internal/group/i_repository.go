package group

import (
	"context"
	"sensors-generator/internal/spiece"
)

type ISensorGroupRepository interface {
	FindAll(ctx context.Context, filters SensorGroupFilters) ([]SensorGroup, error)
	FindOneByID(ctx context.Context, id int, filters SensorGroupFilters) (*SensorGroup, error)
	FindSpiecesInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (map[*spiece.Spiece]int, error)
	FindAvgTransparencyInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (uint8, error)
	FindAvgTemperatureInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (float32, error)
	Create(ctx context.Context, grp CreateSensorGroupDTO) error
}

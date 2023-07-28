package group

import (
	"context"
	"sensors-generator/internal/spiece"
)

type ISensorGroupService interface {
	GetAll(ctx context.Context, filters SensorGroupFilters) ([]SensorGroup, error)
	GetSpiecesInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (map[*spiece.Spiece]int, error)
	GetAvgTrasparencyInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (uint8, error)
	GetAvgTemperatureInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (float32, error)
	Create(ctx context.Context, groups ...CreateSensorGroupDTO) error
}

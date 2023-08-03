package generator

import (
	"context"
	"sensors-generator/internal/group"
	"sensors-generator/internal/sensor"
	"sensors-generator/internal/spiece"
)

type MainEntities struct {
	Groups  []group.CreateSensorGroupDTO
	Sensors []sensor.CreateSensorDTO
	Spieces []spiece.CreateSpieceDTO
}

type MainEntitiesGenerator struct {
	mainEntities MainEntities
	services     Services
}

func NewMainEntitiesGenerator(mainEntities MainEntities, services Services) *MainEntitiesGenerator {
	return &MainEntitiesGenerator{
		mainEntities: mainEntities,
		services:     services,
	}
}

func (gssg *MainEntitiesGenerator) Generate() error {
	if err := gssg.services.SensorGroupService.
		Create(context.Background(), gssg.mainEntities.Groups...); err != nil {
		return err
	}

	if err := gssg.services.SensorService.
		Create(context.Background(), gssg.mainEntities.Sensors...); err != nil {
		return err
	}

	if err := gssg.services.SpieceService.
		Create(context.Background(), gssg.mainEntities.Spieces...); err != nil {
		return err
	}

	return nil
}

func (gssg *MainEntitiesGenerator) IsGenerated() bool {
	if sGroups, _ := gssg.services.SensorGroupService.GetAll(context.Background(), group.SensorGroupFilters{}); sGroups == nil || len(sGroups) <= 0 {
		return false
	}

	return true
}

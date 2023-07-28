package sensor

import (
	"context"
	"sensors-generator/config"
	"sensors-generator/pkg/logging"
)

type service struct {
	sensorRepo ISensorRepository
	logger     *logging.Logger
	cfg        *config.Config
}

func NewService(sensorRepo ISensorRepository,
	logger *logging.Logger, cfg *config.Config) *service {
	return &service{
		sensorRepo: sensorRepo,
		logger:     logger,
		cfg:        cfg,
	}
}

func (s *service) GetAll(ctx context.Context, filters SensorFilters) ([]Sensor, error) {
	s.logger.Info("GET ALL SENSORS.")
	return s.sensorRepo.FindAll(ctx, filters)
}

func (s *service) Create(ctx context.Context, sensors ...CreateSensorDTO) error {
	s.logger.Info("CREATE SENSORS.")
	for _, sensor := range sensors {
		if err := s.sensorRepo.Create(ctx, sensor); err != nil {
			return err
		}
	}

	s.logger.Info("Sensors was created successfully.")
	return nil
}

func (s *service) AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error {
	s.logger.Info("ADD SENSOR TO GROUP")
	return s.sensorRepo.AddSensorToGroup(ctx, sensorID, groupID)
}

func (s *service) GetExtremumTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates, min bool) (float32, error) {
	if min {
		s.logger.Info("GET MIN TEMPERATURE FOR REGION.")
		return s.sensorRepo.FindMinTemperatureForRegion(ctx, minCoords, maxCoords)
	}

	s.logger.Info("GET MAX TEMPERATURE FOR REGION.")
	return s.sensorRepo.FindMaxTemperatureForRegion(ctx, minCoords, maxCoords)
}

func (s *service) GetAvgTemperatureForSensor(ctx context.Context, filters SensorFilters) (float32, error) {
	s.logger.Info("GET AVERAGE TEMPERATURE FOR REGION.")
	return s.sensorRepo.FindAvgTemperatureForSensor(ctx, filters)
}

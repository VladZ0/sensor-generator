package sensordata

import (
	"context"
	"sensors-generator/config"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
)

type service struct {
	sensorDataRepo ISensorDataRepository
	logger         *logging.Logger
	cfg            *config.Config
}

func NewService(sensorDataRepo ISensorDataRepository,
	logger *logging.Logger, cfg *config.Config) *service {
	return &service{
		sensorDataRepo: sensorDataRepo,
		logger:         logger,
		cfg:            cfg,
	}
}

func (s *service) Create(ctx context.Context, sensorData ...CreateSensorDataDTO) ([]int, error) {
	s.logger.Info("CREATE SENSOR DATA.")
	ids := make([]int, 0)
	for _, sd := range sensorData {
		id, err := s.sensorDataRepo.Create(ctx, sd)
		if err != nil {
			ids = append(ids, 0)
			return ids, err
		}
		ids = append(ids, id)
	}

	s.logger.Info("Sensor data created successfully.")
	return ids, nil
}

func (s *service) AddDetectedSpieces(ctx context.Context, sensorDataID int, spieces ...spiece.Spiece) error {
	s.logger.Info("ADD DETECTED SPIECES.")
	for _, spiece := range spieces {
		if err := s.sensorDataRepo.AddDetectedSpiece(ctx, sensorDataID, spiece); err != nil {
			return err
		}
	}

	s.logger.Info("Detected spieces was added successfully.")
	return nil
}

func (s *service) GetOneByID(ctx context.Context, id int, filters SensorDataFilters) (*SensorData, error) {
	s.logger.Info("GET SENSOR DATA.")
	return s.sensorDataRepo.FindOneByID(ctx, id, filters)
}

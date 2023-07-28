package group

import (
	"context"
	"sensors-generator/config"
	"sensors-generator/internal/spiece"
	clients "sensors-generator/pkg/client/interfaces"
	"sensors-generator/pkg/client/redis"
	"sensors-generator/pkg/logging"
	"strconv"
)

type service struct {
	sensorGroupRepo ISensorGroupRepository
	cache           clients.Cache
	logger          *logging.Logger
	cfg             *config.Config
}

func NewService(sensorGroupRepo ISensorGroupRepository, cache clients.Cache,
	logger *logging.Logger, cfg *config.Config) *service {
	return &service{
		sensorGroupRepo: sensorGroupRepo,
		cache:           cache,
		logger:          logger,
		cfg:             cfg,
	}
}

func (s *service) GetSpiecesInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (map[*spiece.Spiece]int, error) {
	s.logger.Info("GET SPIECES IN GROUP.")
	return s.sensorGroupRepo.FindSpiecesInGroup(ctx, groupName, filters)
}

func (s *service) GetAvgTrasparencyInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (uint8, error) {
	s.logger.Info("GET AVERAGE TRANSPARENCY IN GROUP.")
	transparencyKey := groupName + "AvgTransparency"
	var transparency uint8

	transparencyValue, err := s.cache.Get(ctx, transparencyKey)
	if err != nil {
		s.logger.Warn("Transparency not found in cache.")
	}

	if transparencyValue == "" {
		transparency, err = s.sensorGroupRepo.FindAvgTransparencyInGroup(ctx, groupName, filters)
		if err != nil {
			return 0, err
		}

		if err := s.cache.Set(ctx, transparencyKey, transparency, redis.TTL); err != nil {
			s.logger.Errorf("Cannot set cache value: %d.", transparency)
			return 0, err
		}
	} else {
		transparencyInt, err := strconv.Atoi(transparencyValue)
		if err != nil {
			s.logger.Errorf("String \"%s\" cannot be converted to int.", transparencyValue)
			return 0, err
		}
		transparency = uint8(transparencyInt)
	}

	return transparency, nil
}

func (s *service) GetAvgTemperatureInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (float32, error) {
	s.logger.Info("GET AVERAGE TEMPERATURE IN GROUP.")
	temperatureKey := groupName + "AvgTemperature"
	var temperature float32

	temperatureValue, err := s.cache.Get(ctx, temperatureKey)
	if err != nil {
		s.logger.Warn("Temperature not found in cache.")
	}

	if temperatureValue == "" {
		temperature, err = s.sensorGroupRepo.FindAvgTemperatureInGroup(ctx, groupName, filters)
		if err != nil {
			return 0, err
		}

		if err := s.cache.Set(ctx, temperatureKey, temperature, redis.TTL); err != nil {
			s.logger.Errorf("Cannot set cache value: %f.", temperature)
			return 0, err
		}
	} else {
		temperatureF64, err := strconv.ParseFloat(temperatureValue, 32)
		if err != nil {
			s.logger.Errorf("String \"%s\" cannot be converted to float.", temperatureValue)
			return 0, err
		}
		temperature = float32(temperatureF64)
	}

	return temperature, nil
}

func (s *service) Create(ctx context.Context, groups ...CreateSensorGroupDTO) error {
	s.logger.Info("CREATE SENSOR GROUPS.")
	for _, grp := range groups {
		if err := s.sensorGroupRepo.Create(ctx, grp); err != nil {
			return err
		}
	}
	s.logger.Info("Sensor groups was created successfully.")
	return nil
}

func (s *service) GetAll(ctx context.Context, filters SensorGroupFilters) ([]SensorGroup, error) {
	s.logger.Info("GET SENSOR GROUPS.")
	return s.sensorGroupRepo.FindAll(ctx, filters)
}

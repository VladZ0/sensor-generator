package spiece

import (
	"context"
	"sensors-generator/config"
	"sensors-generator/pkg/logging"
)

type service struct {
	spieceRepo ISpiecesRepository
	logger     *logging.Logger
	cfg        *config.Config
}

func NewService(spieceRepo ISpiecesRepository,
	logger *logging.Logger, cfg *config.Config) *service {
	return &service{
		spieceRepo: spieceRepo,
		logger:     logger,
		cfg:        cfg,
	}
}

func (s *service) GetAll(ctx context.Context, filters SpieceFilters) ([]Spiece, error) {
	s.logger.Info("GET SPIECES.")
	return s.spieceRepo.FindAll(ctx, filters)
}

func (s *service) Create(ctx context.Context, spieces ...CreateSpieceDTO) error {
	s.logger.Info("CREATE SPIECES.")
	for _, spiece := range spieces {
		if err := s.spieceRepo.Create(ctx, spiece); err != nil {
			return err
		}
	}

	s.logger.Info("Spieces was created successfully.")
	return nil
}

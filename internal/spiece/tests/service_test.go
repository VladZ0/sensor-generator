package spiece

import (
	"context"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SpieceService_GetAll(t *testing.T) {
	logging.Init("trace", true)
	repo := &MockSpieceRepository{}

	service := spiece.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	filters := spiece.SpieceFilters{}

	expectedSpieces := []spiece.Spiece{
		{
			ID:        1,
			Name:      "Spiece1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		{
			ID:        2,
			Name:      "Spiece2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		{
			ID:        3,
			Name:      "Spiece3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	repo.On("FindAll", ctx, filters).Return(expectedSpieces, nil)

	spieces, err := service.GetAll(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, spieces, expectedSpieces)

	repo.AssertExpectations(t)
}

func Test_SpieceService_Create(t *testing.T) {
	repo := &MockSpieceRepository{}

	service := spiece.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()

	spieces := []spiece.CreateSpieceDTO{
		{
			Name: "Spiece1",
		},

		{
			Name: "Spiece2",
		},

		{
			Name: "Spiece3",
		},
	}

	repo.On("Create", ctx, spieces[0]).Return(nil)
	repo.On("Create", ctx, spieces[1]).Return(nil)
	repo.On("Create", ctx, spieces[2]).Return(nil)

	err := service.Create(ctx, spieces...)

	assert.NoError(t, err)

	repo.AssertCalled(t, "Create", ctx, spieces[0])
	repo.AssertCalled(t, "Create", ctx, spieces[1])
	repo.AssertCalled(t, "Create", ctx, spieces[2])
}

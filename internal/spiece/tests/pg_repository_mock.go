package spiece

import (
	"context"
	"sensors-generator/internal/spiece"

	"github.com/stretchr/testify/mock"
)

type MockSpieceRepository struct {
	mock.Mock
}

func (m *MockSpieceRepository) FindAll(ctx context.Context, filters spiece.SpieceFilters) ([]spiece.Spiece, error) {
	args := m.Called(ctx, filters)

	return args.Get(0).([]spiece.Spiece), args.Error(1)
}

func (m *MockSpieceRepository) FindOneByID(ctx context.Context, id int, filters spiece.SpieceFilters) (*spiece.Spiece, error) {
	args := m.Called(ctx, id, filters)

	return args.Get(0).(*spiece.Spiece), args.Error(1)
}

func (m *MockSpieceRepository) Create(ctx context.Context, spiece spiece.CreateSpieceDTO) error {
	args := m.Called(ctx, spiece)

	return args.Error(0)
}

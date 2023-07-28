package spiece

import "context"

type ISpiecesService interface {
	GetAll(ctx context.Context, filters SpieceFilters) ([]Spiece, error)
	Create(ctx context.Context, spieces ...CreateSpieceDTO) error
}

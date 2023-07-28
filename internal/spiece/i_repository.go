package spiece

import "context"

type ISpiecesRepository interface {
	FindAll(ctx context.Context, filters SpieceFilters) ([]Spiece, error)
	FindOneByID(ctx context.Context, id int, filters SpieceFilters) (*Spiece, error)
	Create(ctx context.Context, spiece CreateSpieceDTO) error
}

package spiece

import "time"

type Spiece struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSpieceDTO struct {
	Name string `json:"name"`
}

type SpieceFilters struct {
	N         int
	GroupName string
}

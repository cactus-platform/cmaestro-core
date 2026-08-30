package models

import (
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Revision string    `json:"revision"`
	Hash     string    `json:"hash"`
	Size     int64     `json:"size"`
	Format   string    `json:"format"`
	Status   string    `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

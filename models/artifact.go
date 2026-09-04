package models

import (
	"time"

	"github.com/google/uuid"
)

type Artifact struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	RepositoryID uuid.UUID `json:"repository_id" gorm:"index"`

	Name string `json:"name"`
	Path string `json:"path" gorm:"unique"`

	Revision string `json:"revision"`

	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
	Format string `json:"format"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

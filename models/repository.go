package models

import (
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID        uuid.UUID   `json:"id" gorm:"primaryKey"`
	Name      string      `json:"name" gorm:"unique"`
	Status    string      `json:"status"`
	Artifacts []*Artifact `json:"artifacts" gorm:"foreignKey:RepositoryID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

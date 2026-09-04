package models

import (
	"time"

	"github.com/google/uuid"
)

type IngestStatus string

const (
	IngestStatusPending IngestStatus = "pending"
)

type Ingest struct {
	RepositoryID uuid.UUID    `json:"repository_id"`
	Revision     uuid.UUID    `json:"revision"`
	Status       IngestStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

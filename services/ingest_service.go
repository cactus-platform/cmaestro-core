package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cactus-platform/cmaestro-core/models"

	"github.com/cactus-platform/cmaestro-core/storage/keyval"
)

type IngestService interface {
	Ingest(ctx context.Context, repository *models.Repository) error
}

type IngestServiceImpl struct {
	keyVal *keyval.Client
}

func NewIngestService(keyVal *keyval.Client) IngestService {
	return &IngestServiceImpl{keyVal: keyVal}
}

func (s *IngestServiceImpl) Ingest(ctx context.Context, repository *models.Repository) error {
	if s.keyVal == nil {
		return errors.New("key-val client cannot be nil")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if len(repository.Artifacts) == 0 {
		return errors.New("no artifacts defined")
	}

	now := time.Now().UTC()
	value, err := json.Marshal(models.Ingest{
		RepositoryID: repository.ID,
		Revision:     repository.Artifacts[0].ID,
		Status:       models.IngestStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return err
	}

	return s.keyVal.Set(
		"ingest:"+repository.ID.String(),
		string(value),
		-1,
	)
}

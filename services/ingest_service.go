package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cactus-platform/cmaestro-core/models"
)

type IngestService interface {
	Ingest(ctx context.Context, repository *models.Repository) error
	Update(ctx context.Context, repository *models.Repository, status models.IngestStatus) error
	Delete(ctx context.Context, repository *models.Repository) error
	IngestStatusReader
}

type IngestStatusReader interface {
	Get(ctx context.Context, repository *models.Repository) (*models.Ingest, error)
}

type IngestStore interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type IngestServiceImpl struct {
	keyVal IngestStore
}

func NewIngestService(keyVal IngestStore) IngestService {
	return &IngestServiceImpl{keyVal: keyVal}
}

func (s *IngestServiceImpl) Ingest(ctx context.Context, repository *models.Repository) error {
	if s.keyVal == nil {
		return errors.New("key-val client cannot be nil")
	}
	if repository == nil {
		return errors.New("repository cannot be nil")
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

func (s *IngestServiceImpl) Update(ctx context.Context, repository *models.Repository, status models.IngestStatus) error {
	if s.keyVal == nil {
		return errors.New("key-val client cannot be nil")
	}
	if repository == nil {
		return errors.New("repository cannot be nil")
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
		Status:       status,
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

func (s *IngestServiceImpl) Delete(ctx context.Context, repository *models.Repository) error {
	if s.keyVal == nil {
		return errors.New("key-val client cannot be nil")
	}
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return s.keyVal.Delete("ingest:" + repository.ID.String())
}

func (s *IngestServiceImpl) Get(ctx context.Context, repository *models.Repository) (*models.Ingest, error) {
	if s.keyVal == nil {
		return nil, errors.New("key-val client cannot be nil")
	}
	if repository == nil {
		return nil, errors.New("repository cannot be nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	value, err := s.keyVal.Get("ingest:" + repository.ID.String())
	if err != nil {
		return nil, err
	}

	var ingest models.Ingest
	if err := json.Unmarshal([]byte(value), &ingest); err != nil {
		return nil, err
	}

	return &ingest, nil
}

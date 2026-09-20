package services

import (
	"context"
	"errors"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/repositories"
	"github.com/google/uuid"
)

type RepositoryService interface {
	Create(ctx context.Context, repository *models.Repository) error
	CreateRevision(ctx context.Context, repository *models.Repository) error
	Get(ctx context.Context, id uuid.UUID) (*models.Repository, error)
	Update(ctx context.Context, repository *models.Repository) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	CreateOrUpdate(ctx context.Context, repository *models.Repository) error
}

type RepositoryServiceImpl struct {
	repository repositories.RepositoryRepository
	ingest     IngestStatusReader
}

func NewRepositoryService(
	repository repositories.RepositoryRepository,
	ingest IngestStatusReader,
) RepositoryService {
	return &RepositoryServiceImpl{
		repository: repository,
		ingest:     ingest,
	}
}

func (s *RepositoryServiceImpl) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.Repository, error) {
	repository, err := s.repository.Get(ctx, id)
	if err != nil {
		return repository, err
	}

	ingest, err := s.ingest.Get(ctx, repository)
	if err != nil {
		return nil, err
	}

	repository.Status = string(ingest.Status)
	for _, artifact := range repository.Artifacts {
		if artifact != nil {
			artifact.Status = string(ingest.Status)
		}
	}

	return repository, nil
}

func (s *RepositoryServiceImpl) Create(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	return s.repository.Create(ctx, repository)
}

func (s *RepositoryServiceImpl) CreateRevision(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	return s.repository.CreateRevision(ctx, repository)
}

func (s *RepositoryServiceImpl) Update(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	return s.repository.Update(ctx, repository)
}

func (s *RepositoryServiceImpl) Exists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	return s.repository.Exists(ctx, id)
}

func (s *RepositoryServiceImpl) CreateOrUpdate(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	exists, err := s.repository.Exists(ctx, repository.ID)
	if err != nil {
		return err
	}

	if exists {
		return s.repository.CreateRevision(ctx, repository)
	}

	return s.repository.Create(ctx, repository)
}

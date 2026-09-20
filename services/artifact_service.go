package services

import (
	"context"
	"errors"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/repositories"

	"github.com/google/uuid"
)

type ArtifactService interface {
	CreateOrUpdateArtifact(ctx context.Context, artifact *models.Artifact) error
	GetArtifact(ctx context.Context, artifactID uuid.UUID) (*models.Artifact, error)
}

type ArtifactServiceImpl struct {
	repository repositories.ArtifactRepository
	ingest     IngestStatusReader
}

func NewArtifactService(
	repository repositories.ArtifactRepository,
	ingest IngestStatusReader,
) ArtifactService {
	return &ArtifactServiceImpl{
		repository: repository,
		ingest:     ingest,
	}
}

func (s *ArtifactServiceImpl) CreateOrUpdateArtifact(
	ctx context.Context,
	artifact *models.Artifact,
) error {
	if artifact == nil {
		return errors.New("artifact cannot be nil")
	}

	exists, err := s.repository.ArtifactExists(ctx, artifact.ID)
	if err != nil {
		return err
	}

	if exists {
		return s.repository.UpdateArtifact(ctx, artifact)
	}

	return s.repository.CreateArtifact(ctx, artifact)
}

func (s *ArtifactServiceImpl) GetArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
) (*models.Artifact, error) {
	artifact, err := s.repository.GetArtifact(ctx, artifactID)
	if err != nil {
		return artifact, err
	}

	ingest, err := s.ingest.Get(ctx, &models.Repository{ID: artifact.RepositoryID})
	if err != nil {
		return nil, err
	}

	artifact.Status = string(ingest.Status)
	return artifact, nil
}

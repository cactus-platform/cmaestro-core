package repositories

import (
	"context"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/google/uuid"
)

type RepositoryRepository interface {
	Create(ctx context.Context, repository *models.Artifact) error
	Get(ctx context.Context, id uuid.UUID) (*models.Artifact, error)
	Update(ctx context.Context, repository *models.Artifact) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type RepositoryRepositoryImpl struct {
	artifactRepository ArtifactRepository
}

func NewRepositoryRepository(
	artifactRepository ArtifactRepository,
) RepositoryRepository {
	return &RepositoryRepositoryImpl{
		artifactRepository: artifactRepository,
	}
}

func (r *RepositoryRepositoryImpl) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.Artifact, error) {
	return r.artifactRepository.GetArtifact(ctx, id)
}

func (r *RepositoryRepositoryImpl) Create(
	ctx context.Context,
	repository *models.Artifact,
) error {
	return r.artifactRepository.CreateArtifact(ctx, repository)
}

func (r *RepositoryRepositoryImpl) Update(
	ctx context.Context,
	repository *models.Artifact,
) error {
	return r.artifactRepository.UpdateArtifact(ctx, repository)
}

func (r *RepositoryRepositoryImpl) Exists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	return r.artifactRepository.ArtifactExists(ctx, id)
}

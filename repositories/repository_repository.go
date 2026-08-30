package repositories

import (
	"context"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/google/uuid"
)

type RepositoryRepository interface {
	Create(ctx context.Context, repository *models.Repository) error
	Get(ctx context.Context, id uuid.UUID) (*models.Repository, error)
	Update(ctx context.Context, repository *models.Repository) error
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
) (*models.Repository, error) {
	artifact, err := r.artifactRepository.GetArtifact(ctx, id)
	if err != nil {
		return nil, err
	}

	return repositoryFromArtifact(artifact), nil
}

func (r *RepositoryRepositoryImpl) Create(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return nil
	}

	return r.artifactRepository.CreateArtifact(ctx, artifactFromRepository(repository))
}

func (r *RepositoryRepositoryImpl) Update(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return nil
	}

	return r.artifactRepository.UpdateArtifact(ctx, artifactFromRepository(repository))
}

func (r *RepositoryRepositoryImpl) Exists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	return r.artifactRepository.ArtifactExists(ctx, id)
}

func artifactFromRepository(repository *models.Repository) *models.Artifact {
	return &models.Artifact{
		Id:        repository.ID,
		Name:      repository.Name,
		Path:      repository.Path,
		Revision:  repository.Revision,
		Hash:      repository.Hash,
		Size:      repository.Size,
		Format:    repository.Format,
		Status:    repository.Status,
		CreatedAt: repository.CreatedAt,
		UpdatedAt: repository.UpdatedAt,
	}
}

func repositoryFromArtifact(artifact *models.Artifact) *models.Repository {
	return &models.Repository{
		ID:        artifact.Id,
		Name:      artifact.Name,
		Path:      artifact.Path,
		Revision:  artifact.Revision,
		Hash:      artifact.Hash,
		Size:      artifact.Size,
		Format:    artifact.Format,
		Status:    artifact.Status,
		CreatedAt: artifact.CreatedAt,
		UpdatedAt: artifact.UpdatedAt,
	}
}

package repositories

import (
	"context"
	"errors"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/storage/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrRepositoryNotFound = errors.New("repository not found")

type RepositoryRepository interface {
	Create(ctx context.Context, repository *models.Repository) error
	CreateRevision(ctx context.Context, repository *models.Repository) error
	Get(ctx context.Context, id uuid.UUID) (*models.Repository, error)
	Update(ctx context.Context, repository *models.Repository) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type RepositoryRepositoryImpl struct {
	db *sql.Client
}

func NewRepositoryRepository(
	db *sql.Client,
) RepositoryRepository {
	return &RepositoryRepositoryImpl{
		db: db,
	}
}

func (r *RepositoryRepositoryImpl) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.Repository, error) {
	var repository models.Repository
	err := r.db.WithContext(ctx).
		Preload("Artifacts").
		Where("id = ?", id).
		First(&repository).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRepositoryNotFound
	}
	return &repository, err
}

func (r *RepositoryRepositoryImpl) Create(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Artifacts").Create(repository).Error; err != nil {
			return err
		}
		for _, artifact := range repository.Artifacts {
			if artifact == nil {
				continue
			}
			artifact.RepositoryID = repository.ID
			if err := tx.Create(artifact).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RepositoryRepositoryImpl) CreateRevision(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Repository{}).
			Where("id = ?", repository.ID).
			Updates(map[string]any{
				"status":     repository.Status,
				"updated_at": repository.UpdatedAt,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrRepositoryNotFound
		}

		for _, artifact := range repository.Artifacts {
			if artifact == nil {
				continue
			}
			artifact.RepositoryID = repository.ID
			if err := tx.Create(artifact).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RepositoryRepositoryImpl) Update(
	ctx context.Context,
	repository *models.Repository,
) error {
	if repository == nil {
		return errors.New("repository cannot be nil")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Repository{}).
		Where("id = ?", repository.ID).
		Updates(repository)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRepositoryNotFound
	}
	return nil
}

func (r *RepositoryRepositoryImpl) Exists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Repository{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

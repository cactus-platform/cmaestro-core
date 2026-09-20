package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/repositories"
	"github.com/google/uuid"
)

type fakeRepositoryRepository struct {
	repository *models.Repository
	exists     bool
	err        error
	created    *models.Repository
	revised    *models.Repository
	updated    *models.Repository
}

func (f *fakeRepositoryRepository) Create(_ context.Context, repository *models.Repository) error {
	f.created = repository
	return f.err
}

func (f *fakeRepositoryRepository) CreateRevision(_ context.Context, repository *models.Repository) error {
	f.revised = repository
	return f.err
}

func (f *fakeRepositoryRepository) Get(_ context.Context, _ uuid.UUID) (*models.Repository, error) {
	return f.repository, f.err
}

func (f *fakeRepositoryRepository) Update(_ context.Context, repository *models.Repository) error {
	f.updated = repository
	return f.err
}

func (f *fakeRepositoryRepository) Exists(_ context.Context, _ uuid.UUID) (bool, error) {
	return f.exists, f.err
}

type fakeRepositoryIngestReader struct {
	value *models.Ingest
	err   error
	seen  *models.Repository
}

func (f *fakeRepositoryIngestReader) Get(_ context.Context, repository *models.Repository) (*models.Ingest, error) {
	f.seen = repository
	return f.value, f.err
}

func TestRepositoryServiceGetAddsIngestStatus(t *testing.T) {
	repositoryID := uuid.New()
	repository := &models.Repository{
		ID: repositoryID,
		Artifacts: []*models.Artifact{
			{ID: uuid.New()},
			nil,
		},
	}
	status := &fakeRepositoryIngestReader{value: &models.Ingest{Status: "processing"}}
	service := NewRepositoryService(&fakeRepositoryRepository{repository: repository}, status)

	got, err := service.Get(context.Background(), repositoryID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Status != "processing" || got.Artifacts[0].Status != "processing" {
		t.Fatalf("status was not propagated: %+v", got)
	}
	if status.seen != repository {
		t.Fatal("expected ingest status lookup to receive the repository")
	}
}

func TestRepositoryServiceCreateOrUpdate(t *testing.T) {
	repository := &models.Repository{ID: uuid.New()}
	fake := &fakeRepositoryRepository{}
	service := NewRepositoryService(fake, &fakeRepositoryIngestReader{})

	if err := service.CreateOrUpdate(context.Background(), repository); err != nil {
		t.Fatalf("CreateOrUpdate() create error = %v", err)
	}
	if fake.created != repository {
		t.Fatal("expected repository to be created")
	}

	fake.exists = true
	if err := service.CreateOrUpdate(context.Background(), repository); err != nil {
		t.Fatalf("CreateOrUpdate() revision error = %v", err)
	}
	if fake.revised != repository {
		t.Fatal("expected repository revision to be created")
	}
}

func TestRepositoryServicePropagatesErrors(t *testing.T) {
	wantErr := errors.New("repository dependency failed")
	repositoryService := NewRepositoryService(&fakeRepositoryRepository{err: wantErr}, &fakeRepositoryIngestReader{})

	if _, err := repositoryService.Get(context.Background(), uuid.New()); !errors.Is(err, wantErr) {
		t.Fatalf("Get() error = %v, want %v", err, wantErr)
	}
	if err := repositoryService.Create(context.Background(), nil); err == nil {
		t.Fatal("expected nil repository error")
	}
	if err := repositoryService.CreateOrUpdate(context.Background(), &models.Repository{ID: uuid.New()}); !errors.Is(err, wantErr) {
		t.Fatalf("CreateOrUpdate() error = %v, want %v", err, wantErr)
	}
}

var _ repositories.RepositoryRepository = (*fakeRepositoryRepository)(nil)

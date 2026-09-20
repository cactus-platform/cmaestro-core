package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/repositories"
	"github.com/google/uuid"
)

type fakeArtifactRepository struct {
	artifact       *models.Artifact
	exists         bool
	err            error
	created        *models.Artifact
	updated        *models.Artifact
	existsRequests []uuid.UUID
}

func (f *fakeArtifactRepository) CreateArtifact(_ context.Context, artifact *models.Artifact) error {
	f.created = artifact
	return f.err
}

func (f *fakeArtifactRepository) GetArtifact(_ context.Context, id uuid.UUID) (*models.Artifact, error) {
	return f.artifact, f.err
}

func (f *fakeArtifactRepository) UpdateArtifact(_ context.Context, artifact *models.Artifact) error {
	f.updated = artifact
	return f.err
}

func (f *fakeArtifactRepository) ArtifactExists(_ context.Context, id uuid.UUID) (bool, error) {
	f.existsRequests = append(f.existsRequests, id)
	return f.exists, f.err
}

type fakeArtifactIngestReader struct {
	value *models.Ingest
	err   error
	seen  *models.Repository
}

func (f *fakeArtifactIngestReader) Get(_ context.Context, repository *models.Repository) (*models.Ingest, error) {
	f.seen = repository
	return f.value, f.err
}

func TestArtifactServiceGetAddsIngestStatus(t *testing.T) {
	artifactID := uuid.New()
	repositoryID := uuid.New()
	artifact := &models.Artifact{ID: artifactID, RepositoryID: repositoryID}
	status := &fakeArtifactIngestReader{value: &models.Ingest{Status: "ready"}}
	service := NewArtifactService(&fakeArtifactRepository{artifact: artifact}, status)

	got, err := service.GetArtifact(context.Background(), artifactID)
	if err != nil {
		t.Fatalf("GetArtifact() error = %v", err)
	}
	if got.Status != "ready" || status.seen.ID != repositoryID {
		t.Fatalf("artifact status = %q, repository ID = %v", got.Status, status.seen.ID)
	}
}

func TestArtifactServiceCreateOrUpdate(t *testing.T) {
	repository := &fakeArtifactRepository{}
	service := NewArtifactService(repository, &fakeArtifactIngestReader{})
	artifact := &models.Artifact{ID: uuid.New()}

	if err := service.CreateOrUpdateArtifact(context.Background(), artifact); err != nil {
		t.Fatalf("create error = %v", err)
	}
	if repository.created != artifact {
		t.Fatal("expected artifact to be created")
	}

	repository.exists = true
	if err := service.CreateOrUpdateArtifact(context.Background(), artifact); err != nil {
		t.Fatalf("update error = %v", err)
	}
	if repository.updated != artifact {
		t.Fatal("expected artifact to be updated")
	}
}

func TestArtifactServicePropagatesErrors(t *testing.T) {
	wantErr := errors.New("artifact dependency failed")
	service := NewArtifactService(&fakeArtifactRepository{err: wantErr}, &fakeArtifactIngestReader{})

	if _, err := service.GetArtifact(context.Background(), uuid.New()); !errors.Is(err, wantErr) {
		t.Fatalf("GetArtifact() error = %v, want %v", err, wantErr)
	}
	if err := service.CreateOrUpdateArtifact(context.Background(), &models.Artifact{ID: uuid.New()}); !errors.Is(err, wantErr) {
		t.Fatalf("CreateOrUpdateArtifact() error = %v, want %v", err, wantErr)
	}
	if err := service.CreateOrUpdateArtifact(context.Background(), nil); err == nil {
		t.Fatal("expected nil artifact error")
	}
}

var _ repositories.ArtifactRepository = (*fakeArtifactRepository)(nil)

package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/google/uuid"
)

type fakeIngestStore struct {
	values  map[string]string
	setKey  string
	setData string
	delKey  string
	err     error
}

func (f *fakeIngestStore) Set(key, value string, _ time.Duration) error {
	f.setKey, f.setData = key, value
	if f.err != nil {
		return f.err
	}
	if f.values == nil {
		f.values = make(map[string]string)
	}
	f.values[key] = value
	return nil
}

func (f *fakeIngestStore) Get(key string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.values[key], nil
}

func (f *fakeIngestStore) Delete(key string) error {
	f.delKey = key
	if f.err != nil {
		return f.err
	}
	delete(f.values, key)
	return nil
}

func TestIngestServiceLifecycle(t *testing.T) {
	repositoryID := uuid.New()
	artifactID := uuid.New()
	repository := &models.Repository{
		ID: repositoryID,
		Artifacts: []*models.Artifact{{
			ID: artifactID,
		}},
	}
	store := &fakeIngestStore{}
	service := &IngestServiceImpl{keyVal: store}

	if err := service.Ingest(context.Background(), repository); err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}
	value, err := service.Get(context.Background(), repository)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if value.RepositoryID != repositoryID || value.Revision != artifactID || value.Status != models.IngestStatusPending {
		t.Fatalf("unexpected ingest value: %+v", value)
	}

	if err := service.Update(context.Background(), repository, "complete"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	value, err = service.Get(context.Background(), repository)
	if err != nil {
		t.Fatalf("Get() after Update error = %v", err)
	}
	if value.Status != "complete" {
		t.Fatalf("status = %q, want complete", value.Status)
	}

	if err := service.Delete(context.Background(), repository); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if store.delKey != "ingest:"+repositoryID.String() {
		t.Fatalf("Delete() key = %q", store.delKey)
	}
}

func TestIngestServiceRejectsInvalidInput(t *testing.T) {
	service := &IngestServiceImpl{keyVal: &fakeIngestStore{}}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		repo *models.Repository
		want string
	}{
		{"nil repository", context.Background(), nil, "repository cannot be nil"},
		{"no artifacts", context.Background(), &models.Repository{}, "no artifacts defined"},
		{"cancelled context", cancelled, &models.Repository{Artifacts: []*models.Artifact{{ID: uuid.New()}}}, "context canceled"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.Ingest(test.ctx, test.repo)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestIngestServicePropagatesStoreErrors(t *testing.T) {
	wantErr := errors.New("store failed")
	repository := &models.Repository{Artifacts: []*models.Artifact{{ID: uuid.New()}}}
	service := &IngestServiceImpl{keyVal: &fakeIngestStore{err: wantErr}}

	if err := service.Ingest(context.Background(), repository); !errors.Is(err, wantErr) {
		t.Fatalf("Ingest() error = %v, want %v", err, wantErr)
	}
	if _, err := service.Get(context.Background(), repository); !errors.Is(err, wantErr) {
		t.Fatalf("Get() error = %v, want %v", err, wantErr)
	}
	if err := service.Delete(context.Background(), repository); !errors.Is(err, wantErr) {
		t.Fatalf("Delete() error = %v, want %v", err, wantErr)
	}
}

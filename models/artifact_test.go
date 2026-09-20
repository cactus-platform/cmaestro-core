package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestArtifactJSONRoundTrip(t *testing.T) {
	artifact := Artifact{
		ID:           uuid.New(),
		RepositoryID: uuid.New(),
		Name:         "archive.zip",
		Path:         "archives/archive.zip",
		Revision:     "revision-1",
		Hash:         "hash",
		Size:         42,
		Format:       "zip",
		Status:       "ready",
	}

	data, err := json.Marshal(artifact)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var got Artifact
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got != artifact {
		t.Fatalf("round trip = %+v, want %+v", got, artifact)
	}
}

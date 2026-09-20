package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestRepositoryJSONRoundTrip(t *testing.T) {
	repository := Repository{
		ID:     uuid.New(),
		Name:   "example",
		Status: "ready",
		Artifacts: []*Artifact{{
			ID:   uuid.New(),
			Name: "archive.zip",
		}},
	}

	data, err := json.Marshal(repository)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var got Repository
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.ID != repository.ID || got.Name != repository.Name || got.Status != repository.Status || len(got.Artifacts) != 1 || got.Artifacts[0].Name != "archive.zip" {
		t.Fatalf("round trip = %+v, want %+v", got, repository)
	}
}

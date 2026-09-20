package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestIngestJSONRoundTrip(t *testing.T) {
	ingest := Ingest{
		RepositoryID: uuid.New(),
		Revision:     uuid.New(),
		Status:       IngestStatusPending,
	}

	data, err := json.Marshal(ingest)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var got Ingest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.RepositoryID != ingest.RepositoryID || got.Revision != ingest.Revision || got.Status != ingest.Status {
		t.Fatalf("round trip = %+v, want %+v", got, ingest)
	}
}

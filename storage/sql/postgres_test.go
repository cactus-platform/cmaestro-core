package sql

import "testing"

func TestNewReturnsConnectionErrorWhenPostgresIsUnavailable(t *testing.T) {
	_, err := New(Config{
		Endpoint: "127.0.0.1",
		Username: "user",
		Password: "pass",
		Database: "database",
		Port:     1,
	})
	if err == nil {
		t.Fatal("New() expected a connection error")
	}
}

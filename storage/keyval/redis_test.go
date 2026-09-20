package keyval

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func newTestClient(t *testing.T) (*Client, *miniredis.Miniredis) {
	t.Helper()
	server := miniredis.RunT(t)
	client, err := New(Config{Addr: server.Addr()})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client, server
}

func TestClientStringLifecycle(t *testing.T) {
	client, server := newTestClient(t)

	if err := client.Set("key", "value", 0); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	value, err := client.Get("key")
	if err != nil || value != "value" {
		t.Fatalf("Get() = %q, %v; want value, nil", value, err)
	}
	exists, err := client.Exists("key")
	if err != nil || !exists {
		t.Fatalf("Exists() = %v, %v; want true, nil", exists, err)
	}
	if err := client.Expire("key", time.Minute); err != nil {
		t.Fatalf("Expire() error = %v", err)
	}
	if err := client.Delete("key"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if server.Exists("key") {
		t.Fatal("Delete() did not remove key")
	}
}

func TestClientJSONAndCounters(t *testing.T) {
	client, _ := newTestClient(t)

	want := map[string]string{"status": "ready"}
	if err := client.SetJSON("json", want, 0); err != nil {
		t.Fatalf("SetJSON() error = %v", err)
	}
	var got map[string]string
	if err := client.GetJSON("json", &got); err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if got["status"] != want["status"] {
		t.Fatalf("decoded status = %q, want %q", got["status"], want["status"])
	}
	count, err := client.Increment("count")
	if err != nil || count != 1 {
		t.Fatalf("Increment() = %d, %v; want 1, nil", count, err)
	}
	if err := client.FlushDB(); err != nil {
		t.Fatalf("FlushDB() error = %v", err)
	}
}

package bucket

import (
	"context"
	"strings"
	"testing"
)

func TestNewValidatesConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{"missing endpoint", Config{Bucket: "bucket", AccessKey: "key", SecretKey: "secret"}, "endpoint is required"},
		{"invalid endpoint", Config{Endpoint: "://bad", Bucket: "bucket", AccessKey: "key", SecretKey: "secret"}, "invalid SeaweedFS endpoint"},
		{"missing bucket", Config{Endpoint: "http://localhost:8333", AccessKey: "key", SecretKey: "secret"}, "bucket is required"},
		{"missing credentials", Config{Endpoint: "http://localhost:8333", Bucket: "bucket"}, "access key and secret key are required"},
		{"small part size", Config{Endpoint: "http://localhost:8333", Bucket: "bucket", AccessKey: "key", SecretKey: "secret", PartSize: 1}, "part size must be at least"},
		{"invalid concurrency", Config{Endpoint: "http://localhost:8333", Bucket: "bucket", AccessKey: "key", SecretKey: "secret", Concurrency: -1}, "concurrency must be positive"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(context.Background(), test.cfg)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("New() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestObjectKeyAndPrefixValidation(t *testing.T) {
	client := &Client{root: "root/", bucket: "bucket"}

	key, err := client.objectKey("archives\\file.ZIP")
	if err != nil || key != "root/archives/file.ZIP" {
		t.Fatalf("objectKey() = %q, %v", key, err)
	}
	for _, name := range []string{"", "folder/", "../file.zip", "file.tar"} {
		if _, err := client.objectKey(name); err == nil {
			t.Fatalf("objectKey(%q) expected error", name)
		}
	}
	prefix, err := normalizePrefix("/application/data/")
	if err != nil || prefix != "application/data/" {
		t.Fatalf("normalizePrefix() = %q, %v", prefix, err)
	}
	if _, err := normalizePrefix("../data"); err == nil {
		t.Fatal("normalizePrefix() expected traversal error")
	}
	if got := zipContentType("archive.zip"); got != "application/zip" {
		t.Fatalf("zipContentType() = %q", got)
	}
}

package dbutil

import (
	"bytes"
	"io"
	"testing"
)

func TestHasherReaderHashesAllReadContent(t *testing.T) {
	reader := NewHasherReader(bytes.NewBufferString("hello world"))

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("read data = %q", data)
	}
	if got, want := reader.Hash(), "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"; got != want {
		t.Fatalf("Hash() = %q, want %q", got, want)
	}
}

func TestHasherReaderHashIsIncompleteUntilRead(t *testing.T) {
	reader := NewHasherReader(bytes.NewBufferString("content"))
	if got := reader.Hash(); got == "9a0364b9e99bb480dd25e1f0284c8555e5e4f5f6" {
		t.Fatal("Hash() unexpectedly included unread content")
	}
}

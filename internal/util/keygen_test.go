package util

import (
	"testing"
)

func TestGenerateReadableKey(t *testing.T) {
	key, err := GenerateReadableKey(16, 4)
	if err != nil {
		t.Fatalf("Error generating readable key: %v", err)
	}
	if len(key) != 19 { // 16 chars + 3 hyphens
		t.Errorf("Expected key length 19, got %d", len(key))
	}
	t.Logf("Generated Readable Key: %s", key)

	key, err = GenerateReadableKey(32, 0)
	if err != nil {
		t.Fatalf("Error generating readable key: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}
	t.Logf("Generated Readable Key: %s", key)
}

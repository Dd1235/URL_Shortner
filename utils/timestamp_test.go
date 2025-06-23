package utils

import (
	"testing"
)

func TestTimestampGenerator_Generate(t *testing.T) {
	gen := &TimestampGenerator{}

	code := gen.Generate("ignored")

	// Check that output is exactly 8 characters long
	if len(code) != 8 {
		t.Errorf("Expected length 8, got %d, value: %s", len(code), code)
	}

	// Check that it's URL-safe base64 (only contains allowed characters)
	for _, c := range code {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			t.Errorf("Invalid character in base64 URL-encoded string: %c", c)
		}
	}
}

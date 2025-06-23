// utils/random_test

package utils

import (
	"testing"
)

func TestRandomGenerator_Generate(t *testing.T) {
	l := 8
	rg := &RandomGenerator{8}
	got := rg.Generate("don't care (heavy accent, iykyk)")
	if len(got) != l {
		t.Errorf("Expected length %d Got %d", l, len(got))
	}
	for _, c := range got {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			t.Errorf("Invalid character in base64 URL-encoded string: %c", c)
		}
	}
}

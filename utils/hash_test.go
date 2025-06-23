// utils/hash_test.go
package utils

import "testing"

func TestHashGenerator_Generate(t *testing.T) {
	hg := &HashGenerator{}
	var got string = hg.Generate("https://example.com") // auto dereferenced
	if len(got) != 8 {
		t.Errorf("Expected 8-char hash, got: %s", got)
	}
}

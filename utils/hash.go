package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

type HashGenerator struct{}

func (g *HashGenerator) Generate(input string) string {
	hash := sha256.Sum256([]byte(input)) // you get a 32-byte array, 256/8

	return base64.RawStdEncoding.EncodeToString(hash[:])[:8] // first 8 characters

	// this is deterministic, 6 bits per character (encode to base 64), 8 characters, so 2^48 possibilities.

	// works good maybe for under a million urls
}

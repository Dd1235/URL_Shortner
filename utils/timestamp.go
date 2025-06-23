package utils

import (
	"encoding/base64"
	"encoding/binary"
	"time"
)

type TimestampGenerator struct{}

func (g *TimestampGenerator) Generate(_ string) string {
	ts := time.Now().UnixNano()
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(ts)) // encode a 64 bit unsigned integer into 8-byte slice
	// eg [12 34 56 78 90 AB CD EF]  most significant byte firts, like 123 1 is the msdigit

	key := base64.RawURLEncoding.EncodeToString(b)

	if len(key) > 8 {
		return key[len(key)-8:]
	}
	return key

}

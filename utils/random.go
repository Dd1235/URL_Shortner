// utils/random.go

package utils

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandomGenerator struct {
	Length int
}

func (g *RandomGenerator) Generate(_ string) string {

	b := make([]byte, g.Length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// say length = 6, thn 62^6 = 56 billion different combinations.
// consider 1 million urls in a day, in 10 years that is 3.65billion
// but for about 50% chance of collision we only need about 240K urls. (its like the birthday paradox)
// so should check for uniqueness
// but increasing hte length to 7 ~ 3.5 trillion

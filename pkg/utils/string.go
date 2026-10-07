package utils

import (
	"bytes"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type KeyGenerator interface {
	GenerateKey(length int) string
}

type randomStringGenerator struct {
	rng *rand.Rand
}

func NewKeyGenerator() KeyGenerator {
	return &randomStringGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
func (r *randomStringGenerator) GenerateKey(length int) string {
	return randomString(r.rng, length)
}

func GenRandomString(length int) string {
	// Seed the random number generator with the current time

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randomString(rng, length)
}

func randomString(rng *rand.Rand, length int) string {
	var strBuilder bytes.Buffer
	for i := 0; i < length; i++ {
		strBuilder.WriteByte(charset[rng.Intn(len(charset))])
	}
	return strBuilder.String()
}

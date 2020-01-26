// Package cmd from cobra
package cmd

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a human readable string of length X
// TODO: Move these to a core lib
func RandomString(length int) string {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for idx := range b {
		b[idx] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// RandomInteger return random int no larger than max
func RandomInteger(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomBool return random boolen
func RandomBool() bool {
	return rand.Intn(2) == 0
}

// RandomFloat returns a random float
func RandomFloat() float64 {
	return rand.Float64()
}

// RandomUUID generate a new UUID
func RandomUUID() uuid.UUID {
	return uuid.New()
}

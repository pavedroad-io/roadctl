/*
 */

//Note: above and below blank lines required for golint.
//Related to required documentation format for packages.

package cmd

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString Returns a random string of length specified
//
func RandomString(length int) string {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for idx := range b {
		b[idx] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// RandomInteger Returns a random int between speified min and max int
//
func RandomInteger(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomBool Returns true or false
//
func RandomBool() bool {
	return rand.Intn(2) == 0
}

// RandomFloat Returns a random float64
//
func RandomFloat() float64 {
	return rand.Float64()
}

// RandomUUID Returns a random UUID
//
func RandomUUID() uuid.UUID {
	return uuid.New()
}

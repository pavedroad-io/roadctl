package cmd

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for idx := range b {
		b[idx] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func RandomInteger(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomBool() bool {
	return rand.Intn(2) == 0
}

func RandomFloat() float64 {
	return rand.Float64()
}

func RandomUUID() uuid.UUID {
	return uuid.New()
}

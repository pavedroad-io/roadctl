// Package cmd from cobra
package cmd

/*
Copyright © 2019 PavedRoad <info@pavedroad.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

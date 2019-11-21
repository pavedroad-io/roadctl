/*
 */

//Note: above and below blank lines required for golint.
//Related to required documentation format for packages.

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

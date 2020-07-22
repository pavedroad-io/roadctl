// Package cmd from gobra
//   Valid JSON and Go type supported in definitions files
//   With associated methods

package cmd

import (
	"fmt"
	"time"
)

// mappedTypes defines valid Go types
type mappedTypes struct {
	validDefinitionType string // acceptable types in a definition file
	name                string // go built in type or defined type
	random              interface{}
	quoted              bool
}

// builtInTypes
//
var builtInTypes = []mappedTypes{
	{"string", "string", RandomString(15), true},
	{"boolean", "bool", RandomBool(), false},
	{"number", "int", RandomInteger(0, 254), false},
	{"int", "int", RandomInteger(0, 254), false},
	{"uint", "uint", RandomInteger(0, 254), false},
	{"byte", "byte", RandomString(1), true},
	{"rune", "byte", RandomString(1), true},
	{"float", "float64", RandomFloat(), false},
	{"float64", "float64", RandomFloat(), false},
	{"float32", "float32", RandomFloat(), false},
	{"uuid", "string", RandomUUID(), true},
	{"time", "time.Time", time.Now().Format(time.RFC3339), true},
}

func (t *mappedTypes) validInputType(key string) bool {
	for _, input := range builtInTypes {
		if input.validDefinitionType == key {
			return true
		}
	}
	fmt.Println("Failure: ", key, " not found")
	return false
}

func (t *mappedTypes) inputToGoType(key string) string {
	for _, input := range builtInTypes {
		if input.validDefinitionType == key {
			return input.name
		}
	}
	return ""
}

func (t *mappedTypes) randomJSONData(key string) interface{} {
	for _, input := range builtInTypes {
		if input.validDefinitionType == key {
			if input.quoted {
				v := fmt.Sprintf("\"%v\"", input.random)
				return v
			} else {
				return input.random
			}
		}
	}
	return ""
}

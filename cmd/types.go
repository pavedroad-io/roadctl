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
}

// builtInTypes
//
var builtInTypes = []mappedTypes{
	{"string", "string", RandomString(15)},
	{"boolean", "bool", RandomBool()},
	{"number", "int", RandomInteger(0, 254)},
	{"int", "int", RandomInteger(0, 254)},
	{"uint", "uint", RandomInteger(0, 254)},
	{"byte", "byte", RandomString(1)},
	{"rune", "byte", RandomString(1)},
	{"float", "float", RandomFloat()},
	{"uuid", "uuid", RandomString(44)},
	{"time", "Time.time", time.Now().Format(time.RFC3339)},
}

func (t *mappedTypes) validInputType(key string) bool {
	for _, input := range builtInTypes {
		if input.validDefinitionType == key {
			fmt.Println("Success: ", input.validDefinitionType, key)
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

func (t *mappedTypes) randomGoData(key string) interface{} {
	for _, input := range builtInTypes {
		if input.validDefinitionType == key {
			return input.random
		}
	}
	return ""
}

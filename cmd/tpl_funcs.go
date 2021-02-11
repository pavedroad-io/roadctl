// cmd - Template function maps
//
// Blocks enable quickly adding core functionality to a microservice, function, or CRD
//

package cmd

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Template function maps
//
func stringFunctionMap() template.FuncMap {
	stringFuncMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"ToCamel": strcase.ToCamel,
		"ToSnake": strcase.ToSnake,
	}
	return stringFuncMap
}

func lookupFunctionMap(name string) template.FuncMap {
	switch name {
	case "stringFunctionMap()":
		return stringFunctionMap()
	}
	return nil
}

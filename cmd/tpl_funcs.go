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

func sonarCloudFunctoinMap(defs bpDef) template.FuncMap {

	//Sonarcloud
	si := defs.findIntegration("sonarcloud")

	if si.Name != "" {
		sonarFuncMap := template.FuncMap{
			"SonarKey": func() string {
				return si.SonarCloudConfig.Key
			},
			"SonarLogin": func() string {
				return si.SonarCloudConfig.Login
			},
			"SonarPrefix": func() string {
				return SONARPREFIX
			},
			"SonarCloudEnabled": func() bool {
				return si.Enabled
			},
		}
		return sonarFuncMap
	}
	return nil
}

func lookupFunctionMap(name string, defs bpDef) template.FuncMap {
	switch name {
	case "stringFunctionMap()":
		return stringFunctionMap()
	case "sonarCloudFunctoinMap()":
		return sonarCloudFunctoinMap(defs)
	}
	return nil
}

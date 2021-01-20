// cmd - code fragements
//
// Blocks enable quickly adding core functionality to a microservice, function, or CRD
//

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//
// HTTP blocks configuration objects
//

//
// endpointConfig template variables for generating routes and
// associates handlers
type endpointConfig struct {
	MicroServiceName string   `json:"microServiceName"`
	APIVersion       string   `json:"apiVersion"`
	Namespace        string   `json:"namespace"`
	ResourceType     string   `json:"resourceType"`
	EndPointName     string   `json:"endPointName"`
	Methods          []string `json:"methods"` // HTTP methods
}

// GenerateRoutes given a configFragmenet, load its templates
// executing each one for HTTP methods and/or events defined
// return the combined code framgent as []byte
//
// Executing a route template requires all data passed in as
// a single go object, that is the function of the
// tplRouteObject
//
// If a route requires template function maps if required
// TODO move this into, move these dependencies into
// CodeFragment so support programmatic invocation

func (hc *endpointConfig) GenerateRoutes() (configFragment []byte, err error) {
	var configFragments strings.Builder

	for _, m := range hc.Methods {

		if found, fragement := findHTTPTemplate(m, GorillaCodeFragments); found {
			fmt.Printf("Processing : %s:%s\n", m, fragement.FileName)
			tplName := GorillaCodeFragments.BaseDirectory + fragement.FileName
			if fragement.TemplatePtr == nil {
				if tpl, err := loadTemplate(tplName, GorillaCodeFragments.Family); err != nil {
					msg := fmt.Errorf("File not found: [%s][%v'\n", tplName, err)
					return nil, msg
				} else {
					fragement.TemplatePtr = tpl
				}
			}
			var b strings.Builder
			var tplData = tplRouteObject{
				APIVersion:   hc.APIVersion,
				EndPointName: hc.EndPointName,
				Namespace:    hc.Namespace,
				Method:       m,
			}
			e := fragement.TemplatePtr.ExecuteTemplate(&b, fragement.FileName, &tplData)
			if e != nil {
				msg := fmt.Errorf("Template execution failed : [%s][%v]\n", tplName, e)
				return nil, msg
			}
			configFragments.WriteString(b.String())
		}
	}
	return []byte(configFragments.String()), nil
}

// Composite objects for templale execution

// tplRouteObject for building HTTP gorilla routes
type tplRouteObject struct {
	APIVersion   string `json:"apiVersion"`
	Namespace    string `json:"namespace"`
	EndPointName string `json:"endPointName"`
	Method       string `json:"method"`
}

// Template based classes

// CodeFragment defines a family of teplates and the directory location
// Then a list HTTP methods or Kafka events that trigger generating the code
// using the specified templates
type CodeFragment struct {
	Family        string                   `json:"family"`        // Family these templates belong too
	BaseDirectory string                   `json:"baseDirectory"` // Directory relative to TLD of blueprints
	HTTPMappings  []HTTPMethodTemplateMap  `json:"httpMappings"`  // Mapping of methods to templates
	EventMappings []EventMethodTemplateMap `json:"eventMappings"` // Mapping of methods to templates
}

type HTTPMethodTemplateMap struct {
	HTTPMethods []string           `json:"http_methods"` // HTTP methods using this template
	FileName    string             `json:"file_name"`    // Name of the template file
	TemplatePtr *template.Template `json:"templatePtr"`  // Pointer to a compiled template or nil
}

type EventMethodTemplateMap struct {
	Events      []string           `json:"events"`      // Events using this template
	FileName    string             `json:"file_name"`   // Name of the template file
	TemplatePtr *template.Template `json:"templatePtr"` // Pointer to a compiled template or nil
}

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

// findHTTPTemplate give an CodeFragment object and an HTTP method
//   return an HTTPMethodTemplateMap containing template to
//   use for code generation
//
//
func findHTTPTemplate(method string, cf CodeFragment) (bool, *HTTPMethodTemplateMap) {
	for _, m := range cf.HTTPMappings {
		// HTTPMethods holds a list methods supported by a template
		if t, _ := containsString(method, m.HTTPMethods); t == true {
			return true, &m
		}
	}
	return false, nil
}

// loadTemplate reads a template from the provided location,
// parses it and returns a pointer to the parsed template list.
// When ParseFiles is used, the first template in the list is ""
// The templates names are based the {{define "name"}} in the
// template file
//
// To execute a single template, use
//   template.ExecuteTemplate(writer, filename, data)
//
// Location is relative to the currently configured blueprints
// directory, i.e. the default $HOME/.pavedroad/blueprints
//
func loadTemplate(location, name string) (tpl *template.Template, e error) {
	// read blueprint cache
	tc, te := NewBlueprintCache()
	if te.errno != tcSuccess {
		log.Fatalf("Failed to read blueprint cache, Got (%v)\n", te)
	}

	// Test the directory location
	fileList := []string{}
	readPath := tc.location.Location() + location
	if _, err := os.Stat(readPath); os.IsNotExist(err) {
		msg := fmt.Errorf("File not found: [%s]\n", readPath)
		return nil, msg
	}

	fileList = append(fileList, readPath)
	// err := template.Must(template.New("").ParseFiles(fileList...))
	fm := stringFunctionMap()
	tpl, err := template.New(name).Funcs(fm).ParseFiles(fileList...)
	if err != nil {
		msg := fmt.Errorf("Parsing tpl [%s] failed with error [%s]\n", readPath, err)
		return nil, msg
	}
	return tpl, nil
}

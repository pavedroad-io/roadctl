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
	EndPointName     string   `json:"endPointName"`
	Methods          []string `json:"methods"` // HTTP methods
}

// loadFromDefinitions given a definitions file, create a list of one
// or more endpointConfigs.  Normally, we only expect one endpoint for
// a microservice, but allow for more if needed
func (hc *endpointConfig) loadFromDefinitions(defs bpDef) (endPoints []endpointConfig, err error) {

	if len(defs.Project.Endpoints) == 0 {
		// Older blueprints don't have dynamically discovered
		// blueprints so this is not a terminal error
		msg := fmt.Errorf("No end points found")
		return nil, msg
	}

	for _, ep := range defs.Project.Endpoints {
		methodList := []string{}
		for _, m := range ep.Methods {
			// TODO: modify to also extract query parameters
			methodList = append(methodList, m.Method)
		}
		newEP := endpointConfig{
			MicroServiceName: defs.Info.Name,
			APIVersion:       defs.Info.APIVersion,
			Namespace:        defs.Project.Kubernetes.Namespace,
			EndPointName:     ep.Name,
			Methods:          methodList,
		}
		endPoints = append(endPoints, newEP)
	}
	return endPoints, nil
}

// GenerateRoutes given a configFragmenet, load its templates
// executing each one for HTTP methods and/or events defined
// return the combined code fragment as []byte
//
// Executing a route template requires all data passed in as
// a single go object, that is the function of the
// tplRouteObject
//
// If a route requires template function maps if required
// TODO move this into, move these dependencies into
// Block so support programmatic invocation

func (hc *endpointConfig) GenerateBlock(block Block) (configFragment []byte, err error) {
	var configFragments strings.Builder

	for _, m := range hc.Methods {

		if found, fragement := findHTTPTemplate(m, block); found {
			//			fmt.Printf("Loading block: %s:%s\n", m, fragement.FileName)
			tplName := block.BaseDirectory + fragement.Template.FileName
			if fragement.Template.TemplatePtr == nil {
				if tpl, err := loadTemplate(tplName, block.Family); err != nil {
					msg := fmt.Errorf("File not found: [%s][%v'\n", tplName, err)
					return nil, msg
				} else {
					fragement.Template.TemplatePtr = tpl
				}
			}
			var b strings.Builder
			var tplData = tplRouteObject{
				APIVersion:   hc.APIVersion,
				EndPointName: hc.EndPointName,
				Namespace:    hc.Namespace,
				Method:       m,
			}
			e := fragement.Template.TemplatePtr.ExecuteTemplate(&b, fragement.Template.FileName, &tplData)
			if e != nil {
				msg := fmt.Errorf("template execution failed : [%s][%v]", tplName, e)
				return nil, msg
			}
			configFragments.WriteString(b.String())
		}
	}
	return []byte(configFragments.String()), nil
}

// Composite objects for template execution

// tplRouteObject for building HTTP gorilla routes
type tplRouteObject struct {
	APIVersion   string `json:"apiVersion"`
	Namespace    string `json:"namespace"`
	EndPointName string `json:"endPointName"`
	Method       string `json:"method"`
}

// Block defines a family of templates and the directory location
// Then a list HTTP methods or Kafka events that trigger generating the code
// using the specified templates
type Block struct {

	// Inverted namespace ID unique to these templates
	ID string `json:"id"` // io.pavedroard.core.loggers.http_access

	// BlockType A type that determines how this block is processed
	BlockType string `json:"block_type"` // type of block; i.e. template, function

	// Description a human readable description
	Description string `json:"description"` // Friendly Description of this template

	// Family friendly name for this grouping of templates or functions
	Family string `json:"family"` // Family these templates belong too

	// Imports required modules for these templates
	Imports []string `json:"imports"` // Required package imports

	// Language the computer programming language
	Language string `json:"language"` // Programming language

	// BaseDirectory in blueprints repository
	BaseDirectory string `json:"base_drectory"` // Directory relative to TLD of blueprints

	// Mapping methods for functions and templates
	//

	// TemplateMap a simple map
	TemplateMap []TemplateItem `json:"template_map"` // Directory relative to TLD of blueprintsm
	// HTTPMappings templates mapped by HTTP methods
	HTTPMappings []HTTPMethodTemplateMap `json:"http_mappings"` // HTTP to template mappings
	// EventMappings templates mapped by events
	EventMappings []EventMethodTemplateMap `json:"event_appings"` // Event to template mapping

	// TemplateExports
	TemplateExports []ExportedItem `json:"exported_template_variables"` // Directory relative to TLD of blueprintsm
}

// TemplateItem
// Templates that are not tied to events or methods
type TemplateItem struct {

	// FileName the file name of this template in the directory
	FileName string `json:"file_name"` // Name of the template file

	// TemplateFunction the name of the function map required for this template
	TemplateFunction string `json:"template_function"` // Name of the template file

	// TemplatePtr a pointer if the template if already initialized
	TemplatePtr *template.Template `json:"templatePtr"` // Pointer to a compiled template or nil
}

// ExporteddItem
// Template variables exported by this template
type ExportedItem struct {

	// TemplateVar the names or functions available for inclusion in this templates
	TemplateVar string `json:"templateVar"`

	// SourceInDefinitions where in the definitions file this value is populated from
	SourceInDefinitions string `json:"source_in_definitions"`
}

// HTTPMethodTemplateMap a list of methods, the assoicated tpl file,
// and a ptr to a compiled instance of it
type HTTPMethodTemplateMap struct {
	HTTPMethods []string     `json:"http_methods"` // HTTP methods using this template
	Template    TemplateItem `json:"template"`
}

// EventMethodTemplateMap
type EventMethodTemplateMap struct {
	Events   []string     `json:"events"` // Events using this template
	Template TemplateItem `json:"template"`
}

// findHTTPTemplate give an Block object and an HTTP method
//   return an HTTPMethodTemplateMap containing template to
//   use for code generation
//
//
func findHTTPTemplate(method string, cf Block) (bool, *HTTPMethodTemplateMap) {
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

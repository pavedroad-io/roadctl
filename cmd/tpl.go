// cmd template functions
package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"text/template"
)

// TemplateItem information about a Go tpl and its
// associated requirements.  If the template has been parsed
// it also contains a pointer too the template for execution
//
type TemplateItem struct {

	// FileName the file name of this template in the directory
	FileName string `json:"file_name"` // Name of the template file

	// TemplateFunction the name of the function map required for this template
	TemplateFunction interface{} `json:"template_function"` // Name of the template file

	// TemplatePtr a pointer if the template if already initialized
	TemplatePtr *template.Template `json:"templatePtr"` // Pointer to a compiled template or nil
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

// ExportedItem
// Template variables exported by this template
type ExportedItem struct {

	// TemplateVar the names or functions available for inclusion in this templates
	TemplateVar string `json:"templateVar"`

	// SourceInDefinitions where in the definitions file this value is populated from
	SourceInDefinitions string `json:"source_in_definitions"`

	// Description of this blocks capabilities
	Description string `json:"description"`

	// Required is this item required in a blueprint using this block
	Required bool `json:"required"`
}

//
// Composite objects for template execution
// As only one object can be passed to a template
// these objects represent a combination of other objects
//

// tplRouteObject for building HTTP gorilla routes
type tplRouteObject struct {
	// APIVersion i,e. 1, 2, 3
	APIVersion string `json:"apiVersion"`

	// Namespace in Kubernetes cluster
	Namespace string `json:"namespace"`

	// EndPointName name used in URL as a end point
	EndPointName string `json:"endPointName"`

	// Method the HTTP method to create a route for
	Method string `json:"method"`
}

//
// Template functions
//

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
// TODO: rewrite to take Block and input and load the proper
// template funcs
func loadTemplate(location, name string, tplFunc interface{}) (tpl *template.Template, e error) {
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

	var err error

	// TypeOf will return nil is assertion fails
	t := reflect.TypeOf(tplFunc.(template.FuncMap))
	if t == nil {
		tpl, err = template.New(name).ParseFiles(fileList...)
	} else {
		tpl, err = template.New(name).Funcs(tplFunc.(template.FuncMap)).ParseFiles(fileList...)
	}
	if err != nil {
		msg := fmt.Errorf("Parsing tpl [%s] failed with error [%s]\n", readPath, err)
		return nil, msg
	}
	return tpl, nil
}

// findHTTPTemplate give a Block object and an HTTP method
//   return an HTTPMethodTemplateMap containing template to
//   use for code generation
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

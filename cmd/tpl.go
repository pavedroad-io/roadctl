// cmd template functions
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	TPLOutputFile = iota
	TPLOutputByteSlice
)

var tplTypes = map[string]int{
	"OutputFile":  TPLOutputFile,
	"OutputBytes": TPLOutputByteSlice,
}

// TemplateItem information about a Go tpl and its
// associated requirements.  If the template has been parsed
// it also contains a pointer too the template for execution
//
type TemplateItem struct {

	// FileName the file name of this template in
	// the directory and in the template define
	FileName string `yaml:"fileName"`

	// OutputFileName for templates that create files
	OutputFileName string `yaml:"outputFileName"`

	// OutputType controls execution of block processing
	//   - TPLOutputFile writes this template to the file name
	//     specified by OutputFileName
	//   - TPLOutputByteSlice aggregates one or more templates
	//     for use in a TPLOutputFile
	OutputType string `yaml:"outputType"`

	// FilePermissions to set when creating a file
	ExecutePermissions bool `yaml:"executePermissions"`

	// TemplateFunction the name of the function map
	// required for this template
	TemplateFunction interface{} `yaml:"templateFunction"`

	// TemplatePtr a pointer if the template if already
	// initialized or nil
	TemplatePtr *template.Template `yaml:"templatePtr"`

	// Description user friendly description
	Description string `yaml:"description"`
}

// HTTPMethodTemplateMap a list of methods, the assoicated tpl file,
// and a ptr to a compiled instance of it
type HTTPMethodTemplateMap struct {
	// HTTP methods using this template
	HTTPMethods []string `yaml:"httpMethods"`

	// Template a TemplateItem
	Template TemplateItem `yaml:"template"`
}

// EventMethodTemplateMap
type EventMethodTemplateMap struct {
	// Events using this template
	Events []string `yaml:"events"`

	// Template a TemplateItem
	Template TemplateItem `yaml:"template"`
}

// ExportedItem
// Template variables exported by this template
type ExportedItem struct {

	// TemplateVar the names or functions available for inclusion in this templates
	TemplateVar string `yaml:"templateVar"`

	// SourceInDefinitions where in the definitions file this value is populated from
	SourceInDefinitions string `yaml:"sourceInDefinitions"`

	// Description of this blocks capabilities
	Description string `yaml:"description"`

	// Required is this item required in a blueprint using this block
	Required bool `yaml:"required"`
}

//
// Composite objects for template execution
// As only one object can be passed to a template
// these objects represent a combination of other objects
//

// tplRouteObject for building HTTP gorilla routes
type tplRouteObject struct {
	// APIVersion i,e. 1, 2, 3
	APIVersion string `yaml:"apiVersion"`

	// Namespace in Kubernetes cluster
	Namespace string `yaml:"namespace"`

	// EndPointName name used in URL as a end point
	EndPointName string `yaml:"endPointName"`

	// Method the HTTP method to create a route for
	Method string `yaml:"method"`
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
func loadTemplate(location, name string, tplFunc interface{}) (tpl *template.Template, e error) {
	// read blueprint cache
	tc, te := NewBlueprintCache()
	if te.errno != tcSuccess {
		log.Fatalf("Failed to read blueprint cache, Got (%v)\n", te)
	}

	// Test the directory location
	fileList := []string{}
	readPath := filepath.Join(tc.location.Location(), location)
	if _, err := os.Stat(readPath); os.IsNotExist(err) {
		msg := fmt.Errorf("File not found: [%s]\n", readPath)
		return nil, msg
	}

	fileList = append(fileList, readPath)

	// err := template.Must(template.New("").ParseFiles(fileList...))

	var err error

	// TypeOf will return nil if assertion fails
	f, ok := tplFunc.(template.FuncMap)

	// Dynamicly loaded blocks need the function map
	// name mapped to an actual function
	if !ok {
		f = lookupFunctionMap(tplFunc.(string))
		if f != nil {
			ok = true
		}
	}
	if !ok {
		tpl, err = template.New(name).ParseFiles(fileList...)
	} else {
		tpl, err = template.New(name).Funcs(f).ParseFiles(fileList...)
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

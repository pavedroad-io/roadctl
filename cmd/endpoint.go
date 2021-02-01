// cmd - endpoints

package cmd

import (
	"fmt"
	"strings"
)

// endpointConfig Configuration definition for generating
// routes and their associated handlers
type endpointConfig struct {
	// MicroServiceName The microserivce name and is allowed to be different
	// than the endpoint name
	MicroServiceName string `json:"microServiceName"`

	// APIVersion A Kubernetes conforming API version
	APIVersion string `json:"apiVersion"`

	// Namespace A Kubernetes namespace
	Namespace string `json:"namespace"`

	// EndPointName The name of this endpoint
	EndPointName string `json:"endPointName"`

	// Method this endpoint supports
	Methods []string `json:"methods"` // HTTP methods
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
func (hc *endpointConfig) GenerateBlock(block Block) (configFragment []byte, err error) {
	var configFragments strings.Builder

	for _, m := range hc.Methods {

		if found, fragement := findHTTPTemplate(m, block); found {
			//			fmt.Printf("Loading block: %s:%s\n", m, fragement.FileName)
			tplName := block.BaseDirectory + fragement.Template.FileName
			if fragement.Template.TemplatePtr == nil {
				if tpl, err := loadTemplate(tplName, block.Family, fragement.Template.TemplateFunction); err != nil {
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
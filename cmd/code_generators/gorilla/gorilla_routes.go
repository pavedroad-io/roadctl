package cmd

import "fmt"

// endpointConfig template variables for generating routes and assoicated handlers
type endpointConfig struct {
	microServiceName string
	APIVersion       string
	namespace        string
	resourceType     string
	endPointName     string
	methods          []string // HTTP methods
}

// routeTemplate assoicates and HTTP method with a template for generating its route
type routeTemplate struct {
	method   string
	template string // HTTP methods
}

// routeTemplates list mapping methods to templates
var routeTemplates = []routeTemplate{
	routeTemplate{method: "DELETE", template: "keyedRoute.tpl"},
	routeTemplate{method: "GET", template: "keyedRoute.tpl"},
	routeTemplate{method: "HEAD", template: "keyedRoute.tpl"},
	routeTemplate{method: "LIST", template: "listRoute.tpl"},
	routeTemplate{method: "PATCH", template: "keyedRoute.tpl"},
	routeTemplate{method: "POST", template: "keyedRoute.tpl"},
	routeTemplate{method: "PUT", template: "keyedRoute.tpl"},
	routeTemplate{method: "OPTIONS", template: "optionsRoute.tpl"},
	routeTemplate{method: "TRACE", template: "traceRoute.tpl"},
}

func (hc *endpointConfig) GenerateRoutes() (configFragment []byte, err error) {
	configFragment = []byte("hi")

	fmt.Println(hc)
	for _, m := range hc.methods {
		fmt.Println(m)

	}

	return configFragment, nil

}

func (hc *endpointConfig) GenerateRoute() (configFragment []byte, err error) {
	return nil, nil

}

/*
func main() {

	//m := [10]string{"GET", "HEAD"}

	myEndpoint := endpointConfig{microServiceName: "puzzle",
		APIVersion:   "1",
		namespace:    "pr",
		resourceType: "puzzle",
		endPointName: "Puzzle",
		methods:      []string{"GET", "HEAD"},
	}

	routes, _ := myEndpoint.GenerateRoutes()
	fmt.Println(string(routes))

}
*/

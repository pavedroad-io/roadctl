package cmd

// TODO: Updated with actual mapped values
// TODO: Consider a mapping function as part of the template
// GorillaRouteBlocks
var GorillaRouteBlocks Block = Block{
	ID:          "io.pavedroard.http.routers.gorilla",
	Family:      "gorilla/mux",
	BlockType:   "template",
	Description: "Gorilla route generations",
	Language:    "go",
	Imports: []string{
		"github.com/gorilla/mux",
		"_github.com/lib/pq"},
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			Template: TemplateItem{
				FileName:         "keyed_route.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"LIST"},
			Template: TemplateItem{
				FileName:         "list_route.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"POST", "OPTIONS", "TRACE"},
			Template: TemplateItem{
				FileName:         "non_keyed_route.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.AccessLoggerInit}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.AccessLoggerShutdown}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.DebugLoggerName}}",
			SourceInDefinitions: "project.logging.access.name",
		},
	},
}

// GorillaMethodBlocks
var GorillaMethodBlocks Block = Block{
	ID:            "io.pavedroard.http.methods.gorilla",
	BlockType:     "template",
	Description:   "Gorilla method generations",
	Family:        "gorilla/mux",
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			Template: TemplateItem{
				FileName:         "keyed_method.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"LIST"},
			Template: TemplateItem{
				FileName:         "list_method.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"OPTIONS", "TRACE"},
			Template: TemplateItem{
				FileName:         "non_keyed_method.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"POST"},
			Template: TemplateItem{
				FileName:         "post_method.tpl",
				TemplateFunction: "",
				TemplatePtr:      nil},
		},
	},
}

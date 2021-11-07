package cmd

// GorillaRouteBlocks
var GorillaRouteBlocks Block = Block{
	APIVersion: "v1beta",
	Kind:       "EndpointsBlock",
	ID:         "io.pavedroad.http.routers.gorilla",
	Family:     "gorilla/mux",
	Metadata: Metadata{
		Labels: []string{"gorilla", "router", "http"},
		Tags:   []string{"http", "request router"},
		Information: BlockInformation{
			Description: "Gorilla route generator",
			Title:       "Gorilla route generator",
			Contact: Contact{
				Author:       "John Scharber",
				Organization: "PavedRoad",
				Email:        "support@pavedroad.io",
				Website:      "www.pavedroad.io",
				Support:      "pavedroad-io.slack.com",
			},
		},
	},
	UsageRights: UsageRights{
		TermsOfService: "As is",
		Licenses:       "Apache 2",
		AccessToken:    "",
	},
	Language: "go",
	Imports: []string{
		`"github.com/gorilla/mux"`,
		`"github.com/gorilla/handlers"`,
		`_ "github.com/lib/pq"`},
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			Template: TemplateItem{
				FileName:         "keyed_route.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"LIST"},
			Template: TemplateItem{
				FileName:         "list_route.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"POST", "TRACE"},
			Template: TemplateItem{
				FileName:         "non_keyed_route.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"OPTIONS"},
			Template: TemplateItem{
				FileName:         "options_route.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Method}}",
			SourceInDefinitions: "defs.Project.Endpoints.Methods",
		},
		{
			TemplateVar:         "{{.EndPointName}}",
			SourceInDefinitions: "defs.Project.Endpoitns.Name",
		},
		{
			TemplateVar:         "{{.Namespace}}",
			SourceInDefinitions: "defs.Project.Kubernetes.Namespace",
		},
		{
			TemplateVar:         "{{.APIVersion}}",
			SourceInDefinitions: "defs.Project.APIVersion",
		},
		{
			TemplateVar:         "{{.ToCamel}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToLower}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToUpper}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToSnake}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
	},
}

// GorillaMethodBlocks
var GorillaMethodBlocks Block = Block{
	APIVersion: "v1beta",
	Kind:       "EndpointsBlocks",
	ID:         "io.pavedroad.http.methods.gorilla",
	Family:     "gorilla/mux",
	Metadata: Metadata{
		Labels: []string{"gorilla", "methods", "http"},
		Tags:   []string{"http", "methods"},
		Information: BlockInformation{
			Description: "Gorilla method generator",
			Title:       "Gorilla method generator",
			Contact: Contact{
				Author:       "John Scharber",
				Organization: "PavedRoad",
				Email:        "support@pavedroad.io",
				Website:      "www.pavedroad.io",
				Support:      "pavedroad-io.slack.com",
			},
		},
	},
	Language: "go",
	Imports: []string{
		`"github.com/gorilla/mux"`,
		`"github.com/gorilla/handlers"`,
		`_ "github.com/lib/pq"`},
	UsageRights: UsageRights{
		TermsOfService: "As is",
		Licenses:       "Apache 2",
		AccessToken:    "",
	},
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			Template: TemplateItem{
				FileName:         "keyed_method.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"LIST"},
			Template: TemplateItem{
				FileName:         "list_method.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"OPTIONS", "TRACE"},
			Template: TemplateItem{
				FileName:         "keyed_method.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		}, {
			HTTPMethods: []string{"POST"},
			Template: TemplateItem{
				FileName:         "post_method.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Method}}",
			SourceInDefinitions: "defs.Project.Endpoints.Methods",
		},
		{
			TemplateVar:         "{{.EndPointName}}",
			SourceInDefinitions: "defs.Project.Endpoitns.Name",
		},
		{
			TemplateVar:         "{{.Namespace}}",
			SourceInDefinitions: "defs.Project.Kubernetes.Namespace",
		},
		{
			TemplateVar:         "{{.APIVersion}}",
			SourceInDefinitions: "defs.Project.APIVersion",
		},
		{
			TemplateVar:         "{{.ToCamel}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToLower}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToUpper}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToSnake}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
	},
}

// GorillaMethodHookBlocks
var GorillaMethodHookBlocks Block = Block{
	APIVersion: "v1beta",
	Kind:       "EndpointsBlocks",
	ID:         "io.pavedroad.http.methods.gorilla",
	Family:     "gorilla/mux",
	Metadata: Metadata{
		Labels: []string{"gorilla", "methods", "hooks", "http"},
		Tags:   []string{"http", "methods", "hooks"},
		Information: BlockInformation{
			Description: "Gorilla method pre-post processing hooks generator",
			Title:       "Gorilla method pre-post processing hooks generator",
			Contact: Contact{
				Author:       "John Scharber",
				Organization: "PavedRoad",
				Email:        "support@pavedroad.io",
				Website:      "www.pavedroad.io",
				Support:      "pavedroad-io.slack.com",
			},
		},
	},
	Language: "go",
	Imports:  []string{},
	UsageRights: UsageRights{
		TermsOfService: "As is",
		Licenses:       "Apache 2",
		AccessToken:    "",
	},
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PUT", "PATCH", "OPTIONS"},
			Template: TemplateItem{
				FileName:         "method-keyed-hooks.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		},
		{
			HTTPMethods: []string{"POST", "TRACE", "CONNECT"},
			Template: TemplateItem{
				FileName:         "method-hooks.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		},
		{
			HTTPMethods: []string{"LIST"},
			Template: TemplateItem{
				FileName:         "method-list-hooks.tpl",
				TemplateFunction: stringFunctionMap(),
				TemplatePtr:      nil},
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Method}}",
			SourceInDefinitions: "defs.Project.Endpoints.Methods",
		},
		{
			TemplateVar:         "{{.EndPointName}}",
			SourceInDefinitions: "defs.Project.Endpoitns.Name",
		},
		{
			TemplateVar:         "{{.Namespace}}",
			SourceInDefinitions: "defs.Project.Kubernetes.Namespace",
		},
		{
			TemplateVar:         "{{.NameExported}}",
			SourceInDefinitions: "defs.Info.Name",
		},
		{
			TemplateVar:         "{{.APIVersion}}",
			SourceInDefinitions: "defs.Project.APIVersion",
		},
		{
			TemplateVar:         "{{.ToCamel}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToLower}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToUpper}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
		{
			TemplateVar:         "{{.ToSnake}}",
			SourceInDefinitions: "stringFunctionMap()",
		},
	},
}

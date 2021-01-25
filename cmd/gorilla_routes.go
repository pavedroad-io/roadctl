package cmd

// GorillaRouteBlocks
var GorillaRouteBlocks BlockFragment = BlockFragment{
	Family:        "gorilla/mux",
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			FileName:    "keyed_route.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"LIST"},
			FileName:    "list_route.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"POST", "OPTIONS", "TRACE"},
			FileName:    "non_keyed_route.tpl",
			TemplatePtr: nil,
		},
	},
}

// GorillaMethodBlocks
var GorillaMethodBlocks BlockFragment = BlockFragment{
	Family:        "gorilla/mux",
	BaseDirectory: "/blocks/go/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "PUT"},
			FileName:    "keyed_method.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"LIST"},
			FileName:    "list_method.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"OPTIONS", "TRACE"},
			FileName:    "non_keyed_method.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"POST"},
			FileName:    "post_method.tpl",
			TemplatePtr: nil,
		},
	},
}

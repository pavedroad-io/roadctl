package cmd

var GorillaCodeFragments CodeFragment = CodeFragment{
	Family:        "gorilla/mux",
	BaseDirectory: "/common/go/code_fragments/gorilla/",
	HTTPMappings: []HTTPMethodTemplateMap{
		{
			HTTPMethods: []string{"GET", "HEAD", "DELETE", "PATCH", "POST", "PUT"},
			FileName:    "keyed_route.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"LIST"},
			FileName:    "list_route.tpl",
			TemplatePtr: nil,
		}, {
			HTTPMethods: []string{"OPTIONS", "TRACE"},
			FileName:    "non_keyed_route.tpl",
			TemplatePtr: nil,
		},
	},
}

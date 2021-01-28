// cmd
package cmd

// Skaffold
var SkaffoldConfigBlock = Block{
	APIVersion: "v1beta",
	Kind:       "FileTemplate",
	ID:         "io.pavedroard.core.skaffold.config",
	Family:     "pavedroad/core/skaffold",
	Metadata: Metadata{
		Labels: []string{"pavedroad", "skaffold", "skaffold",
			"W3C"},
		Tags: []string{"pavedroad", "skaffold", "skaffold"},
		Information: BlockInformation{
			Description: "Generate skaffold file",
			Title:       "Generate skaffold file",
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
	Language:      "go",
	Imports:       []string{`log "github.com/pavedroad-io/skaffold/skaffold"`},
	BaseDirectory: "/blocks/go/pavedroad/core/skaffold/",
	HomeDirectory: "/manifests/",
	HomeFilename:  "skaffold.yaml",
	Environment:   "dev",

	TemplateMap: []TemplateItem{
		{
			FileName:         "skaffold_config.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Generate skaffold config for this micro-service",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Skaffold.Build}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Skaffold.Deploy}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Skaffold.Profile}}",
			SourceInDefinitions: "",
		},
	},
}

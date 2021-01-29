// cmd
package cmd

// Skaffold
var SkaffoldDBConfigBlock = Block{
	APIVersion: "v1beta",
	Kind:       "FileTemplate",
	ID:         "io.pavedroard.core.skaffold.config",
	Family:     "pavedroad/core/skaffold",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"skaffold",
			"ci",
			"kubernetes",
			"cockroach",
			"debug",
			"microservice",
			"insecure registry",
			"sha256 tagging"},
		Tags: []string{"pavedroad",
			"skaffold",
			"ci",
			"kubernetes",
			"cockroach",
			"debug",
			"microservice",
			"insecure registry",
			"sha256 tagging"},
		Information: BlockInformation{
			Description: "Skaffold file supporting Cockroachdb and debugging",
			Title:       "Skaffold file with cockroachdb",
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
	BaseDirectory: "/blocks/go/pavedroad/core/skaffold/",
	HomeDirectory: "/manifests/",
	HomeFilename:  "skaffold.yaml",
	Environment:   "dev",

	TemplateMap: []TemplateItem{
		{
			FileName:         "skaffold_cockroach.tpl",
			TemplateFunction: stringFunctionMap(),
			TemplatePtr:      nil,
			Description:      "Generate skaffold config for this micro-service",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.SkaffoldExports.Build}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.SkaffoldExports.Deploy}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.SkaffoldExports.Profile}}",
			SourceInDefinitions: "",
		},
	},
}

//
// Input definitions
//

type SkaffoldInputs struct {
	Organization     string
	MicroserviceName string
}

//
// Exports
//

type SkaffoldExports struct {
	Build   string
	Deploy  string
	Profile string
}

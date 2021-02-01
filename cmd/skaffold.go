// cmd
package cmd

// Skaffold
var SkaffoldDBConfigBlock = Block{
	APIVersion: "v1beta",
	Kind:       "SkaffoldBlock",
	ID:         "io.pavedroard.core.skaffold.config",
	Family:     "pavedroad/core/skaffold",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"skaffold",
			"helm3",
			"ci",
			"kubernetes",
			"cockroach",
			"debug",
			"microservice",
			"insecure registry",
			"sha256 tagging"},
		Tags: []string{"pavedroad",
			"go",
			"skaffold",
			"management API",
			"metrics API",
			"RESTAPI",
			"events",
			"ci",
			"kubernetes",
			"database",
			"postgresql",
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
			Description:      "Generate skaffold config for a database micro-service",
		},
	},
	ImportedBlocks: []Block{
		{ID: "io.pavedroard.core.skaffold.build"},
		{ID: "io.pavedroard.core.skaffold.deploy"},
		{ID: "io.pavedroard.core.skaffold.profile",
			Metadata: Metadata{
				Labels: []string{
					"skaffold",
					"profile",
					"go",
					"debug",
				},
			},
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

/*
This is a possible second way of handling these files
Load the major sections and combine them to create the template

Then execute the template using the inputs


type SkaffoldExports struct {
	Build   string
	Deploy  string
	Profile string
}

func (sc *SkaffoldInputs) readAndSave(url url.URL, block Block, in SkaffoldInputs) (err erorr) {

}
*/

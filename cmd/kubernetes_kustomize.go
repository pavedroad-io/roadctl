// cmd
package cmd

//  PRKubernetesKustomize
var PRKubernetesKustomize Block = Block{
	APIVersion: "v1beta",
	Kind:       "CompositeBlock",
	ID:         "io.pavedroard.core.manifests.kubernetes.kustomize",
	Family:     "pavedroad/core/manifests/kubernetes/kustomize",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"kubernetes",
			"kafka",
			"zookepper",
			"kustomize",
			"envconfiguration",
			"dev",
		},
		Tags: []string{
			"pavedroad",
			"kubernetes",
			"kafka",
			"zookepper",
			"kustomize",
			"envconfiguration",
			"dev",
		},
		Information: BlockInformation{
			Description: "Generate kubernetes deployment using kustomize",
			Title:       "Generate kubernetes deployment using kustomize",
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
	ImportedBlocks: []Block{
		{
			ID: "io.pavedroard.core.manifests.kubernetes.kustomize.kafka",
			Metadata: Metadata{
				Labels: []string{
					"kubernetes",
					"kafka",
					"zookepper",
					"kustomize",
					"dev",
				},
			},
		},
	},
	BaseDirectory: "/blocks/go/pavedroad/manifests/kubernetes/kustomize/",
	HomeDirectory: "manifests/kubernetes/dev",

	TemplateMap: []TemplateItem{
		{
			FileName:         "kustomize.tpl",
			TemplateFunction: stringFunctionMap(),
			TemplatePtr:      nil,
			Description:      "Generate top level kustomize configuration",
		},
		{
			FileName:         "namespace.tpl",
			TemplateFunction: stringFunctionMap(),
			TemplatePtr:      nil,
			Description:      "Generate top level kustomize configuration",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Kustomize.Version}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Namespace}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Name}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Kustomize.Resources}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Kustomize.Bases}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Kustomize.Annotations}}",
			SourceInDefinitions: "",
		},
	},
}

// cmd
package cmd

// DockerCompose for main docker-file
var DockerCompose Block = Block{
	APIVersion: "v1beta",
	Kind:       "FileTemplate",
	ID:         "io.pavedroard.core.docker.compose",
	Family:     "pavedroad/core/docker",
	Metadata: Metadata{
		Labels: []string{"pavedroad", "docker", "docker-compose",
			"W3C"},
		Tags: []string{"pavedroad", "docker-compoase", "docker"},
		Information: BlockInformation{
			Description: "Generate docker-compose file",
			Title:       "Generate docker-compose file",
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
	BaseDirectory: "/blocks/go/pavedroad/core/docker/",
	HomeDirectory: "/manifests/",
	HomeFilename:  "docker-compose.yaml",
	Environment:   "dev",

	TemplateMap: []TemplateItem{
		{
			FileName:         "docker_compose.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Generate docker config for this micro-service",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Docker.Service.Environment}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Service.DependsOn}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Kafka}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Zookepper}}",
			SourceInDefinitions: "",
		},
	},
}

// DockerCompose for database only
var DBOnlyDockerCompose Block = Block{
	APIVersion: "v1beta",
	Kind:       "FileTemplate",
	ID:         "io.pavedroard.core.docker.compose",
	Family:     "pavedroad/core/docker",
	Metadata: Metadata{
		Labels: []string{"pavedroad", "docker", "docker-compose",
			"W3C"},
		Tags: []string{"pavedroad", "docker-compoase", "docker"},
		Information: BlockInformation{
			Description: "Generate docker-compose file",
			Title:       "Generate docker-compose file",
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
	BaseDirectory: "/blocks/go/pavedroad/core/docker/",
	HomeDirectory: "/manifests/",
	HomeFilename:  "docker-db-only.yaml",
	Environment:   "dev",

	TemplateMap: []TemplateItem{
		{
			FileName:         "docker_compose.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Generate docker config for this micro-service",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.Docker.Service.Environment}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Service.DependsOn}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Kafka}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.Docker.Zookepper}}",
			SourceInDefinitions: "",
		},
	},
}

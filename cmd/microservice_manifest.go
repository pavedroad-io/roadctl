// cmd
package cmd

// PRMicroServiceManifests
var PRMicroServiceManifests Block = Block{
	APIVersion: "v1beta",
	Kind:       "CompositeBlock",
	ID:         "io.pavedroard.blueprints.microservice.manifests",
	Family:     "pavedroad/blueprints/manifests",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"microservice",
			"manifests",
			"docker",
			"docker-file",
			"skaffold",
			"kubernetes"},
		Tags: []string{
			"pavedroad",
			"microservice",
			"manifests",
			"docker",
			"docker-file",
			"skaffold",
			"kubernetes"},
		Information: BlockInformation{
			Description: "Composite for building manifests for a microservice",
			Title:       "Composite for building manifests for a microservice",
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
			ID: "io.pavedroard.core.manifests.kubernetes.kustomize",
			Metadata: Metadata{
				Labels: []string{
					"kubernetes",
					"kafka",
					"zookepper",
					"kustomize",
					"envconfiguration",
					"microservice",
					"dev",
				},
			},
		},
		{
			ID: "io.pavedroard.core.manifests.docker.docker-compoase",
			Metadata: Metadata{
				Labels: []string{
					"docker",
					"docker-compose",
					"kafka",
					"zookepper",
					"envconfiguration",
					"microservice",
					"dev",
				},
			},
		},
		{
			ID: "io.pavedroard.core.manifests.docker.dockerfile",
			Metadata: Metadata{
				Labels: []string{
					"docker",
					"dockerfile",
					"microservice",
					"dev",
				},
			},
		},
		{
			ID: "io.pavedroard.core.manifests.skaffold.config",
			Metadata: Metadata{
				Labels: []string{
					"skaffold",
					"kustomize",
					"microservice",
					"dev",
					"dev-debug",
				},
			},
		},
	},
}

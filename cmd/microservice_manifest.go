// cmd
package cmd

// PRMicroServiceManifests
var PRMicroServiceManifests Block = Block{
	APIVersion: "v1beta",
	Kind:       "CompositeBlock",
	ID:         "cache://io.pavedroard.blocks/microservice/manifests",
	Family:     "pavedroad/microservice/blocks/manifests",
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
			ID: "cache://io.pavedroard.blocks/microservice/manifests/kubernetes/kustomize",
			Metadata: Metadata{
				Labels: []string{
					"pavedroad",
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
			ID: "cache://io.pavedroard.blocks/microservice/manifests/docker/docker-compose.yaml",
			Metadata: Metadata{
				Labels: []string{
					"pavedroad",
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
			ID: "cache://io.pavedroard.blocks/microservice/manifests/docker/dockerfile.yaml",
			Metadata: Metadata{
				Labels: []string{
					"pavedroad",
					"docker",
					"dockerfile",
					"microservice",
					"dev",
				},
			},
		},
		{
			ID: "cache://io.pavedroard.blocks/microservice/manifests/skaffold",
			Metadata: Metadata{
				Labels: []string{
					"pavedroad",
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

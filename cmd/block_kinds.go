package cmd

const (
	// FileBlock a block generated from a file
	// containing a Go template
	FileBlock = "FileBlock"

	// CoreBlock a block that integrates core libraries
	CoreBlock = "CoreBlock"

	// CompositeBlock a block that integrates core libraries
	CompositeBlock = "CompositeBlock"

	// WebFileBlock a block generated from a file
	// containing a Go template using the HTML safe
	// typing processor
	HTMLFileBlock = "WebFileBlock"

	// DockerBlock generates a dockerfile
	DockerfileBlock = "DockerfileBlock"

	// DockerComposeBlock generates a docker-compose block
	DockerComposeBlock = "DockerComposeBlock"

	// EndpointsBlock generates HTTP endpoint routers and handlers
	EndpointsBlock = "EndpointsBlock"

	// FunctionBlock call a function that generates
	// non template based blocks
	FunctionBlock = "FunctionBlock"

	// KustomizeBlock creates a kustomize configuration
	// or fragments
	KustomizeBlock = "KustomizeBlock"

	// TemplateBlock creates a tpl fragments
	TemplateBlock = "TemplateBlock"

	// SkaffoldBlock creates a skaffold configuration
	// or fragments
	SkaffoldBlock = "SkaffoldBlock"
)

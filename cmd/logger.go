// cmd
package cmd

// PRApplicationLogger
var PRApplicationLogger Block = Block{
	APIVersion: "v1beta",
	Kind:       "PavedRoad.template",
	ID:         "io.pavedroard.core.loggers.application",
	Family:     "pavedroad/core/logger",
	Metadata: Metadata{
		Labels: []string{"pavedroad", "logger", "http access",
			"W3C"},
		Tags: []string{"pavedroad", "HTTP access logger"},
		Information: BlockInformation{
			Description: "Application debug log",
			Title:       "Application debug log",
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
	Imports:       []string{`log "github.com/pavedroad-io/go-core/logger"`},
	BaseDirectory: "/blocks/go/pavedroad/core/logger/",
	TemplateMap: []TemplateItem{
		{
			FileName:         "docker_compose.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Generate docker config for this micro-service",
		},
		{
			FileName:         "kubernetes_kafka_service.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Create a k8s service for kafka",
		},
		{
			FileName:         "kubernetes_kafka_deployment.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Create a k8s deployment for kafka",
		},
		{
			FileName:         "kubernetes_zookeeper_service.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Create a k8s service for zookeeper",
		},
		{
			FileName:         "kubernetes_zookeeper_deployment.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Create a k8s deployment for zookeeper",
		},
		{
			FileName:         "kustomize.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "Kustomize configuration for logging",
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

// Data to be read from definitions and passed to
// the templates or functions
// A list of loggers is allowed in the definitions
// file
type Logger struct {
	// ID BlockID
	ID string

	// APIVersion API version
	APIVersion string

	// Name of the logger
	Name string

	// EnableDocker support
	EnableDocker bool

	// EnableKubernetes support
	EnableKubernetes bool

	// Outputs to enable
	Outputs LogOutput
}

// FileOutput configuration
type FileOutput struct {
	// Enable or disable
	Enable bool

	// Directory to place log in
	Directory string

	// FileName of log
	FileName string
}

// LogOutput configuration
type LogOutput struct {
	// Console logging
	Console bool

	// Disk logging
	Disk FileOutput

	// EventStream logging
	EventStream EventOutput
}

// EventOutput configuration
type EventOutput struct {
	// Enable or disable
	Enable bool

	// OutputTopics a list of topics to publish
	// to
	OutputTopics []Topic

	// EnableCloudEvents use cloud event formatting
	EnableCloudEvents bool
}

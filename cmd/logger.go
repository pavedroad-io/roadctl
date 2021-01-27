package cmd

var PRHTTPAccessLogger Block = Block{
	APIVersion: "v1beta",
	Kind:       "PavedRoad.template",
	ID:         "io.pavedroard.core.loggers.http_access",
	Family:     "pavedroad/core/logger",
	Metadata: Metadata{
		Labels: []string{"pavedroad", "logger", "http access",
			"W3C"},
		Tags: []string{"pavedroad", "HTTP access logger"},
		Information: BlockInformation{
			Description: "HTTP W3C access logger",
			Title:       "HTTP W3C access logger",
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
			FileName:         "debug_logger.tpl",
			TemplateFunction: nil,
			TemplatePtr:      nil,
			Description:      "HTTP W3C access logger",
		},
	},
	TemplateExports: []ExportedItem{
		{
			TemplateVar:         "{{.DockerEnvForService}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.DockerKafka}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.DockerZookepper}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.AccessLoggerInit}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.AccessLoggerShutdown}}",
			SourceInDefinitions: "",
		},
		{
			TemplateVar:         "{{.DebugLoggerName}}",
			SourceInDefinitions: "project.logging.access.name",
		},
	},
}

package cmd

var PRHTTPAccessLogger Block = Block{
	ID:            "io.pavedroard.core.loggers.http_access",
	Family:        "pavedroad/core/logger",
	BlockType:     "template",
	Description:   "HTTP access logger",
	Language:      "go",
	Imports:       []string{"import line 1", "import line 2"},
	BaseDirectory: "/blocks/go/pavedroad/core/logger/",
	TemplateMap: []TemplateItem{
		{
			FileName:         "debug_logger.tpl",
			TemplateFunction: "",
			TemplatePtr:      nil,
		},
	},
	TemplateExports: []ExportedItem{
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

// cmd
package cmd

import "fmt"

// PRApplicationLogger
var PRApplicationLogger Block = Block{
	APIVersion: "v1beta",
	Kind:       "CompositeBlock",
	ID:         "io.pavedroard.core.loggers.application",
	Family:     "pavedroad/core/logger",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"logger",
			"http access",
			"W3C"},
		Tags: []string{
			"pavedroad",
			"HTTP access logger"},
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
	Language: "go",
	Imports:  []string{`log "github.com/pavedroad-io/go-core/logger"`},
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
					"dev",
				},
			},
		},
	},
}

// Data to be read from definitions and passed to
// the templates or functions
// A list of loggers is allowed in the definitions
// file
type Logger struct {
	// ID BlockID
	ID string `json:"id"`

	// APIVersion API version
	APIVersion string `json:"apiVersion"`

	// Name of the logger
	Name string `json:"name"`

	Metadata Metadata `json:"metadata"`

	// AutoInit true or false
	AutoInit bool `json:"autoInit"`

	// ConfigType env
	ConfigType string `json:"configType"`

	// EnableDocker support
	EnableDocker bool `json:"enableDocker"`

	// EnableKubernetes support
	EnableKubernetes bool `json:"enableKubernetes"`

	// Outputs to enable
	Outputs LogOutput `json:"outputs"`
}

// LogOutput configuration
type LogOutput struct {
	// Console logging
	Console bool `json:"console"`

	// Disk logging
	Disk FileOutput `json:"disk"`

	// EventStream logging
	EventStream EventOutput `json:"eventStream"`
}

// FileOutput configuration
type FileOutput struct {
	// Enable or disable
	Enable bool `json:"enable"`

	// FileFormant text, JSON, etc
	FileFormant string `json:"fileFormant"`

	// Directory to place log in
	Directory string `json:"directory"`

	// FileName of log
	FileName string `json:"fileName"`
}

// EventOutput configuration
type EventOutput struct {
	// Enable or disable
	Enable bool `json:"enable"`

	// OutputTopics a list of topics to publish
	// to
	OutputTopics []Topic `json:"outputTopics"`

	// KafkaBrokers "kafka:9092", ....
	KafkaBrokers string `json:"kafkaBrokers"`

	// EnableCloudEvents use cloud event formatting
	EnableCloudEvents bool `json:"enableCloudEvents"`
}

func (l *Logger) getLoggerImports(defs bpDef) (imports []string, err error) {
	b := Block{}
	var uniqueImports []string

	if len(defs.Project.Loggers) == 0 {
		msg := fmt.Errorf("No loggers found")
		return uniqueImports, msg
	}

	for _, l := range defs.Project.Loggers {

		if l.ID == "" {
			msg := fmt.Errorf("Warning no ID specified in call to loadBlock: [%v]'\n", l.ID)
			fmt.Println(msg)

		}

		nb := &Block{}
		if nb, err = b.loadBlock(l.ID, l.Metadata.Labels); err != nil {
			msg := fmt.Errorf("Loading block failed for ID[%s]: [%v]'\n", l.ID)
			fmt.Println(msg)
		}

		for _, i := range nb.Imports {
			if test, _ := containsString(i, uniqueImports); test == true {
				continue
			}
			uniqueImports = append(uniqueImports, i)
		}
	}

	return uniqueImports, nil
}

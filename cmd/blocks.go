// cmd - blocks
//
// Blocks enable quickly adding core functionality to a microservice, function, or CRD
//

package cmd

//
// HTTP blocks configuration objects
//

// Block defines a family of templates and the directory location
// Then a list HTTP methods or Kafka events that trigger generating the code
// using the specified templates
type Block struct {

	// APIVersion version of this block
	// Added to the path for URL and file system
	// access to blocks
	// i.e. templatedir/v1/block
	// TODO: implement this
	APIVersion string `json:"api_version`

	// Kind a type that determines how this block
	// is processed
	// - template
	// - function
	// - fileTemplate a template that also creates thefile
	Kind string `json:"block_type"` // i.e. template, function

	// Metadata for this block
	Metadata Metadata `json:"metadata"`

	// Inverted namespace ID unique to these templates
	// For example, io.pavedroard.core.loggers.http_access
	ID string `json:"id"`

	// Family friendly name for this grouping of templates
	// or functions, for example, gorilla/mux
	Family string `json:"family"`

	// UsageRights for using the block
	UsageRights UsageRights

	// Imports required modules for these templates
	// Required package imports
	Imports []string `json:"imports"`

	// Language the computer programming language
	Language string `json:"language"`

	// BaseDirectory in blueprints repository
	BaseDirectory string `json:"base_drectory"`

	// HomeDirectory to place a fileTemplate in
	HomeDirectory string `json:"home_directory"`

	// HomeFileName to create
	HomeFilename string `json:"home_filename"`

	// Environment this template applies to
	Environment string `json:"environment"`

	//
	// Mapping methods for functions and templates
	//

	// TemplateMap a simple map
	TemplateMap []TemplateItem `json:"template_map"`

	// HTTPMappings templates mapped by HTTP methods
	HTTPMappings []HTTPMethodTemplateMap `json:"http_mappings"`
	// EventMappings templates mapped by events
	EventMappings []EventMethodTemplateMap `json:"event_appings"`

	// TemplateExports variables for templates and the
	// data source that provides them
	TemplateExports []ExportedItem `json:"exported_template_variables"`
}

// UsageRights terms of service, licensing, and access tokens
type UsageRights struct {
	// TermsOfService for example, as is
	TermsOfService string

	// Licenses cost for example,  annual, per use, perpetual
	Licenses string

	// ContributeLink for donation to the developer
	ContributeLink string

	// AccessToken for downloading this block
	AccessToken string
}

// Metrics that support data driven development
// and operations
type Metrics struct {

	// DORA metrics
	DORAStatistics DORA

	// GitHub metrics
	GitHub GitStatistics

	// Operations metrics developed by PavedRoad
	Operations OperationalStatictics
}

// GitStatistics tracked from GitHub repositories holding blocks
type GitStatistics struct {
	Stars     int
	Forks     int
	Clones    int
	Watchers  int
	Downloads int
}

// OperationalStatictics created automatically when deploying on
// the PR SaaS service
type OperationalStatictics struct {
	// NumberOfTimesDeployed
	NumberOfTimesDeployed int

	// ActiveDeployments
	ActiveDeployments int

	// Failures int
	Failures int

	// Response times per HTTP method or pub/sub event
	Performance map[string]int
}

// DORA metrics are a result of six years worth of surveys
// conducted by the DevOps Research and Assessments (DORA) team
// These metrics guide determine how successful a company is
// at DevOps - ranging from elite performer
type DORA struct {
	// DF deployment Frequency
	DF float64

	// MLT mean Lead Time for changes
	MLT float64

	// MTTR Mean Time To Recover
	MTTR float64

	// CFR Change Failure Rate
	CFR float64
}

type Metadata struct {
	// Label's allow blueprints to be associated
	Labels []string `json:"labels"`

	// Tags catagorize blocks for search
	Tags []string `json:"tags"`

	// Information about block author and support
	Information BlockInformation `json:"information"`
}

type BlockInformation struct {
	// Description a paragraph or  two about this block
	Description string `json:"description"`

	// Title a single line description
	Title string `json:"title"`

	// Contact information for suport
	Contact Contact `json:"contact"`
}

// Contact information for this block
type Contact struct {
	// Author of block
	Author string `json:"author"`

	// Organization who built this block
	Organization string `json:"organization"`

	// Email address for support
	Email string `json:"email"`

	// Website for more information
	Website string `json:"website"`

	// Support channel like slack URL
	Support string `json:"support"`
}

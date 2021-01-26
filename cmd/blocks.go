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

	// Inverted namespace ID unique to these templates
	ID string `json:"id"` // io.pavedroard.core.loggers.http_access

	// BlockType A type that determines how this block is processed
	BlockType string `json:"block_type"` // type of block; i.e. template, function

	// Description a human readable description
	Description string `json:"description"` // Friendly Description of this template

	// Family friendly name for this grouping of templates or functions
	Family string `json:"family"` // Family these templates belong too

	// Imports required modules for these templates
	Imports []string `json:"imports"` // Required package imports

	// Language the computer programming language
	Language string `json:"language"` // Programming language

	// BaseDirectory in blueprints repository
	BaseDirectory string `json:"base_drectory"` // Directory relative to TLD of blueprints

	// Mapping methods for functions and templates
	//

	// TemplateMap a simple map
	TemplateMap []TemplateItem `json:"template_map"` // Directory relative to TLD of blueprintsm
	// HTTPMappings templates mapped by HTTP methods
	HTTPMappings []HTTPMethodTemplateMap `json:"http_mappings"` // HTTP to template mappings
	// EventMappings templates mapped by events
	EventMappings []EventMethodTemplateMap `json:"event_appings"` // Event to template mapping

	// TemplateExports
	TemplateExports []ExportedItem `json:"exported_template_variables"` // Directory relative to TLD of blueprintsm
}

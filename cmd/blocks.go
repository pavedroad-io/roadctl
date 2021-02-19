// cmd - blocks
//
// Blocks enable quickly adding core functionality to a microservice, function, or CRD
//

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	//log "github.com/pavedroad-io/go-core/logger"
	"gopkg.in/yaml.v2"
)

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
	APIVersion string `yaml:"apiVersion"`

	// Kind a type that determines how this block
	// is processed
	// - template
	// - function
	// - fileTemplate a template that also creates thefile
	Kind string `yaml:"kind"` // i.e. template, function

	// Metadata for this block
	Metadata Metadata `yaml:"metadata"`

	// Inverted namespace ID unique to these templates
	// For example, io.pavedroard.core.loggers.http_access
	ID string `yaml:"id"`

	// Family friendly name for this grouping of templates
	// or functions, for example, gorilla/mux
	Family string `yaml:"family"`

	// UsageRights for using the block
	UsageRights UsageRights `yaml:"usageRights"`

	// Imports required modules for these templates
	// Required package imports
	Imports []string `yaml:"imports"`

	// ImportedBlocks additional blocks imported by this block
	ImportedBlocks []Block `yaml:"importedBlocks"`

	// Language the computer programming language
	Language string `yaml:"language"`

	// BaseDirectory in blueprints repository
	BaseDirectory string `yaml:"baseDirectory"`

	// HomeDirectory to place a fileTemplate in
	HomeDirectory string `yaml:"homeDirectory"`

	// HomeFileName to create
	HomeFilename string `yaml:"homeFilename"`

	// Environment this template applies to
	Environment string `yaml:"environment"`

	//
	// Mapping methods for functions and templates
	//

	// TemplateMap a simple map
	TemplateMap []TemplateItem `yaml:"templateMap"`

	// HTTPMappings templates mapped by HTTP methods
	HTTPMappings []HTTPMethodTemplateMap `yaml:"httpMappings"`
	// EventMappings templates mapped by events
	EventMappings []EventMethodTemplateMap `yaml:"eventMappings"`

	// TemplateExports variables for templates and the
	// data source that provides them
	TemplateExports []ExportedItem `yaml:"templateExports"`
}

// UsageRights terms of service, licensing, and access tokens
type UsageRights struct {
	// TermsOfService for example, as is
	TermsOfService string `yaml:"termsOfService"`

	// Licenses cost for example,  annual, per use, perpetual
	Licenses string `yaml:"licenses"`

	// ContributeLink for donation to the developer
	ContributeLink string `yaml:"contributeLink"`

	// AccessToken for downloading this block
	AccessToken string `yaml:"accessToken"`
}

// Metrics that support data driven development
// and operations
type Metrics struct {

	// DORA metrics
	DORAStatistics DORA `yaml:"doraStatistics"`

	// GitHub metrics
	GitHub GitStatistics `yaml:"gitHub"`

	// Operations metrics developed by PavedRoad
	Operations OperationalStatictics `yaml:"operations"`
}

// GitStatistics tracked from GitHub repositories holding blocks
type GitStatistics struct {
	Stars     int `yaml:"stars"`
	Forks     int `yaml:"forks"`
	Clones    int `yaml:"clones"`
	Watchers  int `yaml:"watchers"`
	Downloads int `yaml:"downloads"`
}

// OperationalStatictics created automatically when deploying on
// the PR SaaS service
type OperationalStatictics struct {
	// NumberOfTimesDeployed
	NumberOfTimesDeployed int `yaml:"numberOfTimesDeployed"`

	// ActiveDeployments
	ActiveDeployments int `yaml:"activeDeployments"`

	// Failures int
	Failures int `yaml:"failures"`

	// Response times per HTTP method or pub/sub event
	Performance map[string]int `yaml:"performance"`
}

// DORA metrics are a result of six years worth of surveys
// conducted by the DevOps Research and Assessments (DORA) team
// These metrics guide determine how successful a company is
// at DevOps - ranging from elite performer
type DORA struct {
	// DF deployment Frequency
	DF float64 `yaml:"df"`

	// MLT mean Lead Time for changes
	MLT float64 `yaml:"mlt"`

	// MTTR Mean Time To Recover
	MTTR float64 `yaml:"mttr"`

	// CFR Change Failure Rate
	CFR float64 `yaml:"cfr"`
}

// TODO: break into its own go file for use in
// other types
type Metadata struct {
	// Label's allow blueprints to be associated
	Labels []string `yaml:"labels"`

	// Tags catagorize blocks for search
	Tags []string `yaml:"tags"`

	// Information about block author and support
	Information BlockInformation `yaml:"information"`
}

type BlockInformation struct {
	// Description a paragraph or  two about this block
	Description string `yaml:"description"`

	// Title a single line description
	Title string `yaml:"title"`

	// Contact information for suport
	Contact Contact `yaml:"contact"`
}

// Contact information for this block
type Contact struct {
	// Author of block
	Author string `yaml:"author"`

	// Organization who built this block
	Organization string `yaml:"organization"`

	// Email address for support
	Email string `yaml:"email"`

	// Website for more information
	Website string `yaml:"website"`

	// Support channel like slack URL
	Support string `yaml:"support"`
}

// loadBlock populate a block given its ID
// and a set of lables
func (b *Block) loadBlock(ID string, labels []string) (block *Block, err error) {
	var u *url.URL

	if u, err = url.Parse(ID); err != nil {
		log.Fatalf("Failed to parse ID (%s) error (%v)\n", ID, err)
	}

	switch u.Scheme {
	case "cache":
		return b.loadBlockFromCache(u, labels)
	case "http", "https":
		return b.loadBlockFromNetwork(ID, labels)
	}

	// TODO: fix this hack
	switch ID {
	case "io.pavedroard.core.loggers.application":
		b = &PRApplicationLogger
		break
	}

	return b, nil
}

func (b *Block) loadBlockFromCache(u *url.URL, labels []string) (block *Block, err error) {
	tc, te := NewBlueprintCache()
	if te.errno != tcSuccess {
		log.Fatalf("Failed to read blueprint cache, Got (%v)\n", te)
	}

	fn := tc.location.Location() + "/" + u.Host + u.Path

	// If it is a directory, use default.yaml as the file
	fileInfo, err := os.Stat(fn)
	if err != nil {
		fmt.Println("Error:", err)
		return b, err
	}
	if fileInfo.IsDir() {
		fn = checkDefault(fn)
	}

	df, err := os.Open(fn)
	if err != nil {
		fmt.Println("failed to open:", fn, ", error:", err)
	}
	defer df.Close()

	byteValue, e := ioutil.ReadAll(df)
	if e != nil {
		fmt.Println("read failed for ", df)
		os.Exit(-1)
	}

	err = yaml.Unmarshal([]byte(byteValue), b)
	if err != nil {
		fmt.Printf("Unmarshal faild for %v with %v", u.String(), err)
		return b, err
	}

	for i, sb := range b.ImportedBlocks {
		su, _ := url.Parse(sb.ID)
		nb, _ := sb.loadBlockFromCache(su, sb.Metadata.Labels)
		b.ImportedBlocks[i] = *nb

	}

	return b, nil
}

// checkDefault if you yaml file is specified look for
// default.yaml in the directory given
func checkDefault(s string) string {
	if s[len(s)-4:len(s)] == "yaml" && s[len(s)-3:len(s)] == "yml" {
		return s
	}
	return s + "/default.yaml"
}

// loadBlockFromNetwork
func (b *Block) loadBlockFromNetwork(ID string, labels []string) (block *Block, err error) {
	return b, nil
}

// getImports
func (b *Block) getImports() []string {
	return b.Imports
}

// GenerateBlock
func (b *Block) GenerateBlock(def bpDef) (output string, err error) {

	switch b.Kind {
	case SkaffoldBlock, DockerfileBlock, KustomizeBlock, TemplateBlock:
		for _, sb := range b.TemplateMap {
			fn := filepath.Join(b.BaseDirectory, sb.FileName)

			t, e := loadTemplate(fn, sb.FileName, def, sb.TemplateFunction)
			if e != nil {
				nw := bpError{Type: ErrorGeneric, Err: e}
				return "", nw.WrappedError()
			}

			rclog.Println("template ", t)
			var tplResult strings.Builder
			if e := t.ExecuteTemplate(&tplResult, sb.FileName, def); e != nil {
				nw := bpError{Type: ErrorGeneric, Err: e}
				return "", nw.WrappedError()
			}

			rclog.Println("result ", tplResult.String())
			if e := b.saveResults([]byte(tplResult.String()), sb, def); e != nil {
				return "", e
			}
		}
	}
	// Process all imported blocks
	for _, sb := range b.ImportedBlocks {
		if _, err := sb.GenerateBlock(def); err != nil {
			nw := bpError{Type: ErrorGeneric, Err: err}
			return "", nw.WrappedError()
		}
	}
	return "", nil
}

// saveResults writes generated output to the directory
// and file specified in the blocks creating directories as
// needed
func (b *Block) saveResults(buf []byte, ti TemplateItem, def bpDef) (err error) {
	if b.HomeDirectory == "" {
		e := errors.New("Home directory is required")
		nw := bpError{Type: ErrorGeneric, Err: e}
		return nw.WrappedError()
	}
	var ms macroSubstitutions
	macroDirName := ms.replaceAll(b.HomeDirectory, def)

	//	macroDirName := macroSubstition(b.HomeDirectory, , def.Info.Name)
	if _, err := os.Stat(macroDirName); os.IsNotExist(err) {
		err := os.MkdirAll(macroDirName, 0750)
		if err != nil {
			nw := bpError{Type: ErrorGeneric,
				Err: fmt.Errorf("Failed to create directory: %v", b.HomeDirectory)}
			return nw.WrappedError()
		}
	}

	mode := DefaultFileMode

	if ti.ExecutePermissions {
		mode = DefaultExecutable
	}

	//macroFileName := macroSubstition(ti.OutputFileName, substrings, def.Info.Name)
	macroFileName := ms.replaceAll(ti.OutputFileName, def)
	file, err := os.OpenFile(
		filepath.Join(macroDirName, macroFileName),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		log.Fatal(err,
			filepath.Join(macroDirName, macroFileName))
	}

	bw := bufio.NewWriter(file)

	if _, err := bw.Write(buf); err != nil {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Write failed: %v", macroFileName)}
		return ne.WrappedError()
	}

	if err := bw.Flush(); err != nil {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Flush failed: %v", macroFileName)}
		return ne.WrappedError()
	}

	if err := file.Close(); err != nil {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Close failed: %v", macroFileName)}
		return ne.WrappedError()
	}

	return nil
}

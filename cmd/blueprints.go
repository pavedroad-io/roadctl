// Package cmd from cobra
package cmd

/*
Copyright © 2019 PavedRoad <info@pavedroad.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//TODO: create standard error messages as const

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"

	fossa "github.com/pavedroad-io/integrations/fossa/cmd"
	gorpt "github.com/pavedroad-io/integrations/go/cmd"
)

// GitHub repository information for bpPull
var defaultOrg = "pavedroad-io"
var defaultRepo = "blueprints"
var defaultPath string

// Default blueprint directory on local machine
var defaultBlueprintDir = "blueprints"

var repoType = "GitHub"

// From CLI
var bpFile string    // Blueprint name to use
var bpDefFile string // Definition file to use form -f option

// The directory we found this blueprint in
var bpDirSelected string

// BLUEPRINT needs documentation.
const (
	bpResourceName = "blueprints"
	bpDefinition   = "definition.yaml"
	// PREFIX is the prefix to be replaced in front of the file name
	PREFIX = "template"
	// ORGANIZATION is the prefix to be replaced in front of the file name
	ORGANIZATION = "organization"

	// Name to add to sonarcloud projects to create unique namespace
	SONARPREFIX = "PavedRoad_"

	HOOK = "Hook"

	prCopyright = `
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//`
)

const strcutComment = `
//
//
//`

// JSON formatters
const (
	jsonObjectStart = "{\n"
	jsonObjectEnd   = "}"
	jsonListStart   = "["
	jsonListEnd     = "]"
	jsonSeperator   = ",\n"
	jsonField       = "\"%v\": "
	jsonValue       = "\"%v\"" // If new object, or last field strip the comma
	jsonNumber      = "%v"     // If new object, or last field strip the comma
)
const defMicroserviceName = "yourMicroserviceName"
const pavedroadSonarTestOrg = "acme-demo"
const swaggerRoute = "// swagger:response %s\n"
const structOpen = "type %s struct {\n"

// structUUID
// name of table, type of data json|yaml
const structUUID = "\t%sUUID string `%s:\"%suuid\"`\n"
const structClose = "}\n\n"

// structField
//  name, type, encoding, json|yaml, encoding options
const structField = "\t%s %s\t`%s:\"%s\"`\n"

// structSubStruct
// Same as structField except:
//  type will be the sub-table
//  No options
const structSubstruct = "\t%s %s\t`%s:\"%s\"`\n"
const structSubstructList = "\t%s []%s\t`%s:\"%s\"`\n"

// Makefile constants
const (
	allWithFossa    string = "all: mod-setup $(PREFLIGHT) $(FOSSATEST) compile check"
	allWithoutFossa string = "all: mod-setup $(PREFLIGHT) compile check"

	checkWithSonar    string = "check: lint docker-build sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)"
	checkWithoutSonar string = "check: lint docker-build $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)"

	// Fossa has a build section and a lint section
	fossaSection string = `
$(FOSSATEST):
	fossa init
`
	fossaLint string = `
	@echo "  >  running FOSSA license scan."
	$(shell (export GOPATH=$(GOPATH); @FOSSA_API_KEY=$(FOSSA_API_KEY) fossa analyze))
`
)

var blueprints *template.Template

type bpData struct {
	// Information about company and project
	Organization        string // Name of Organization
	OrgSQLSafe          string // Mapped to a safe name for using in SQL
	OrganazationInfo    string // Name of Organization
	OrganizationLicense string // Org license/copyright
	ProjectInfo         string // Project/service description

	//   scheduler to create
	SchedulerName string // For worker polls specifies the type of

	MaintainerName  string
	MaintainerEmail string
	MaintainerSlack string
	MaintainerWeb   string

	// Version Information
	Version    string // Version of this application
	APIVersion string // Version of this API
	TLD        string // Top level domain

	// Kubernetes endpoints
	// Namespace to deploy and use in URLs
	Namespace string

	// Liveness endpoint name
	Liveness string

	// Readiness endpoint name
	Readiness string

	// Metrics endpoint name
	Metrics string

	// Metrics endpoint name
	Management string

	// Returns API documentation in Swagger syntax
	Explain string

	// Integration's
	Badges            string // badges to include docs
	SonarKey          string
	SonarLogin        string
	SonarPrefix       string
	SonarCloudEnabled bool
	FOSSAEnabled      bool

	// Service and bp-names
	Name         string //service name
	NameExported string //camel case with first letter cap
	TplName      string //blueprint name
	DefFile      string //definition file used

	//PR license/copyright should be a function
	PavedroadInfo string //PR license/copyright

	//Swagger headers probably turn these into functions
	// TODO: remove these old swagger doc statements from blueprints
	//       they are now generated with blocks
	AllRoutesSwaggerDoc     string
	GetAllSwaggerDoc        string // swagger for list method
	GetSwaggerDoc           string // swagger for get method
	PutSwaggerDoc           string // swagger for put method
	PostSwaggerDoc          string // swagger for post method
	DeleteSwaggerDoc        string // swagger for delete method
	ExplainSwaggerDoc       string // swagger for explain endpoint
	SwaggerGeneratedStructs string // swagger doc and go struct
	DumpStructs             string // Generic dump of given object type

	// Endpoints blueprint specific
	EndpointRoutes   string // Holds gorilla routes to function initialization
	EndpointHandlers string // Holds methods for each route
	EndpointHooks    string // Generated pre/post hook functions for selected methods

	PrimaryTableName string // Used as the structure name for
	// Storing user data

	// Language specific inputs
	GoImports string // Imports added by digital blocks

	//JSON data
	PostJSON string // Sample data for a post
	PutJSON  string // Sample data for a put

	// Makefile options
	CheckBuildTarget  string //build line for check section
	AllBuildTarget    string //build line for check section
	FossaBuildSection string //build target for Fossa
	FossaLintSection  string //lint section for Fossa
}

//  bpDataMapper
//    Map data from definitions file to bpData structure
//    return error if required mappings are missing
//    TODO: jms
func bpDataMapper(defs bpDef, output *bpData) error {
	// Docker images names don't allow uppercase letters
	output.Name = strings.ToLower(defs.Info.Name)
	output.NameExported = strcase.ToCamel(defs.Info.Name)
	output.TplName = defs.Info.ID
	output.DefFile = bpDefFile
	output.OrganizationLicense = defs.Project.License
	output.Organization = defs.Info.Organization
	if len(defs.TableList) > 0 {
		output.PrimaryTableName = defs.TableList[0].TableName
	}

	// TODO: Write an SQL safe naming function
	output.OrgSQLSafe = strcase.ToCamel(defs.Info.Organization)
	output.ProjectInfo = defs.Project.Description
	output.MaintainerName = defs.Project.Maintainer.Name
	output.MaintainerEmail = defs.Project.Maintainer.Email
	output.MaintainerWeb = defs.Project.Maintainer.Web
	output.MaintainerSlack = defs.Project.Maintainer.Slack
	output.TLD = defs.Project.TLD
	output.SchedulerName = defs.Project.SchedulerName
	output.PavedroadInfo = prCopyright

	// Version info
	output.Version = defs.Info.Version
	output.APIVersion = defs.Info.APIVersion

	// Endpoint mappings
	output.Namespace = defs.Project.Kubernetes.Namespace
	output.Liveness = defs.Project.Kubernetes.Liveness
	output.Readiness = defs.Project.Kubernetes.Readiness
	output.Metrics = defs.Project.Kubernetes.Metrics
	output.Management = defs.Project.Kubernetes.Management

	// CI integrations

	//Sonarcloud
	si := defs.findIntegration("sonarcloud")

	if si.Name != "" {
		output.SonarKey = si.SonarCloudConfig.Key
		output.SonarLogin = si.SonarCloudConfig.Login
		output.SonarPrefix = SONARPREFIX
		output.SonarCloudEnabled = si.Enabled
	}

	if output.SonarCloudEnabled {
		// If sonarcloud is configured validate token and project
		err := validateIntegrations(output)

		if err != nil {
			fmt.Println("Validating integrations failed: ", err)
			os.Exit(-1)
		}
		output.CheckBuildTarget = checkWithSonar
	} else {
		output.CheckBuildTarget = checkWithoutSonar
	}

	bl, err := scBadges(output, si.SonarCloudConfig.Options.Badges)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, v := range bl {
		output.Badges += v
	}

	si = defs.findIntegration("fossa")
	if si.Name != "" {
		output.FOSSAEnabled = si.Enabled
	}

	if output.FOSSAEnabled {
		output.FossaBuildSection = fossaSection
		output.FossaLintSection = fossaLint
		output.AllBuildTarget = allWithFossa
		b := fossa.GetBadge(fossa.HTML, fossa.Shield, output.Name)
		if b != "" {
			output.Badges += b
		}
	} else {
		output.FossaBuildSection = ""
		output.FossaLintSection = ""
		output.AllBuildTarget = allWithoutFossa
	}

	// This assumes the repository name is the same as the microserice
	// fine for Go but not other langauges
	si = defs.findIntegration("go")
	if si.Name != "" && si.Enabled == true {
		b := gorpt.GetGoBadge(gorpt.GoHTMLink, output.Organization, output.Name)
		if b != "" {
			output.Badges += b
		}
	}
	return nil
}

//  bpJSONData
//    Use the schema definition found in bpDefs to create
//    Sample JSON data files
//
func bpJSONData(defs bpDef, output *bpData) error {
	var jsonString string
	order := defs.devineOrder()
	bpAddJSON(order, defs, &jsonString)

	//Make it pretty
	var pj bytes.Buffer
	err := json.Indent(&pj, []byte(jsonString), "", "\t")
	if err != nil {
		log.Fatal("Failed to generate json data with ", jsonString)
	}
	output.PostJSON = string(pj.String())
	output.PutJSON = string(pj.String())

	return nil
}

// bpAddJSON
//
//   Creates JSON sample data
//
func bpAddJSON(item bpTableItem, defs bpDef, jsonString *string) {
	table, _ := defs.tableByName(item.Name)

	// Start this table
	if item.Root {
		*jsonString = fmt.Sprintf(jsonObjectStart)
	} else {
		*jsonString += fmt.Sprintf(jsonField, strings.ToLower(item.Name))
		if item.IsList {
			*jsonString += fmt.Sprintf(jsonListStart)
		}
		*jsonString += fmt.Sprintf(jsonObjectStart)
	}

	// Only add the UUID if this is the parent table
	if item.Root {
		*jsonString += fmt.Sprintf(jsonField, strings.ToLower(item.Name+"UUID"))
		*jsonString += fmt.Sprintf(jsonValue, RandomUUID())
		*jsonString += fmt.Sprintf(jsonSeperator)
	}

	// Add this tables attributes
	numCol := len(table.Columns)
	m := mappedTypes{}
	for idx, col := range table.Columns {

		*jsonString += fmt.Sprintf(jsonField, strings.ToLower(col.Name))
		*jsonString += fmt.Sprintf("%v", m.randomJSONData(col.Type))
		if idx < numCol-1 {
			*jsonString += fmt.Sprintf(jsonSeperator)
		} else {
			*jsonString += fmt.Sprintf("\n")
		}
	}

	// See if there are any children
	if len(item.Children) > 0 {
		*jsonString += fmt.Sprintf(jsonSeperator)
		// Add child tables first
		for _, child := range item.Children {
			bpAddJSON(*child, defs, jsonString)
		}
	}

	// Close and append to bpData.SwaggerGeneratedStructs
	*jsonString += jsonObjectEnd
	if item.IsList {
		*jsonString += fmt.Sprintf(jsonListEnd)
	}

	return
}

//  bpGenerateStructurs
//    Use the schema definition found in bpDefs to create
//    Go structures and assign to bpData.SwaggerGeneratedStructs
//
//    Use the same schema to generate a formated dump
//    command to aid developer debugging and assign it to
//    bpData.DumpStructs
//
func bpGenerateStructurs(defs bpDef, output *bpData) error {
	order := defs.devineOrder()
	bpAddStruct(order, defs, output)
	return nil
}

// bpAddStruct
//
// Performs two tasks
//    - 1 Generates the structure as a string that is inserted
//        into the code blueprint.  This is the "tableString"
//        variable
//
//    - 2 Creates JSON sample data
//        One for insert, and one for updates
//
func bpAddStruct(item bpTableItem, defs bpDef, output *bpData) {
	table, _ := defs.tableByName(item.Name)

	// Start this table
	tableString := fmt.Sprintf(swaggerRoute, item.Name)
	tableString += fmt.Sprintf(structOpen, item.Name)

	// Only add the UUID if this is the parent table
	if item.Root {
		tableString += fmt.Sprintf("// %sUUID into JSONB\n\n",
			strcase.ToCamel(item.Name))

		tableString += fmt.Sprintf(structUUID,
			strcase.ToCamel(item.Name), "json",
			strings.ToLower(item.Name))
	}

	// See if there are any children
	if len(item.Children) > 0 {
		// Add child tables first
		for _, child := range item.Children {
			bpAddStruct(*child, defs, output)

			// Same as structField except type with be the suitable
			str := structSubstruct
			if child.IsList {
				str = structSubstructList
			}

			tableString += fmt.Sprintf(str,
				strcase.ToCamel(child.Name),
				strings.ToLower(child.Name),
				"json",
				strings.ToLower(child.Name))
		}
	}

	// Add this tables attributes
	for _, col := range table.Columns {

		// build json / yaml string
		importLine := col.MappedName
		if col.Modifiers != "" {
			importLine += "," + col.Modifiers
		}

		if col.Constraints != "" {
			importLine += "," + col.Constraints
		}

		tableString += fmt.Sprintf("// %s\n", strcase.ToCamel(col.Name))
		m := mappedTypes{}
		fieldType := m.inputToGoType(strings.ToLower(col.Type))

		//Deal with time types
		tableString += fmt.Sprintf(structField,
			strcase.ToCamel(col.Name),
			fieldType,
			"json",
			importLine)

	}

	// Close and append to bpData.SwaggerGeneratedStructs
	tableString += fmt.Sprintf(structClose)
	output.SwaggerGeneratedStructs += tableString

	return
}

/*
const structClose = "}\n"
// structField
//  name, type, encoding, json|yaml, encoded name, options
const structField = "\t%s %s\t`%s:%s%s`\n"
*/

// bpListItem provides information about a blueprint location
// and status
type bpListItem struct {
	Type         string // Type of blueprint, i.e. serverless
	Availability string // Availability ga, ....
	Name         string // Name of blueprint == directory name
	Path         string // Path to the blueprint
}

// bpLocation
type bpLocation struct {
	Name         string //Name of the blueprint file
	RelativePath string // Path relative to the current directory
}

type bpExplainItem struct {
	Name    string // Name of resource
	Content string // Text for explain document
}

type bpGenericResponseItem struct {
	Name    string // Name of resource
	Content string // Text for explain document
}

type bpDescribeItem struct {
	Type    string // Type of blueprint, i.e. serverless
	Name    string // Name of blueprint == directory name
	Content string // YAML configuration data
}

type bpExplainResponse struct {
	Blueprints []bpExplainItem
}

type bpCreateResponse struct {
	Blueprints []bpGenericResponseItem
}

type bpDescribeResponse struct {
	Blueprints []bpDescribeItem
}

type bpListResponse struct {
	Blueprints []bpListItem
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func (t bpExplainResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
}

func (t bpExplainResponse) RespondWithJSON() string {
	return t.RespondWithText() // One in the same for this type
}

func (t bpExplainResponse) RespondWithText() string {
	nl := ""
	for _, val := range t.Blueprints {
		nl += fmt.Sprintf("Name: %v\n", val.Name)
		nl += fmt.Sprintf("Content: %v\n", val.Content)
	}
	nl += "\n"
	fmt.Println(nl)
	return nl
}

func (t bpCreateResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
}

func (t bpCreateResponse) RespondWithJSON() string {
	return t.RespondWithText() // One in the same for this type
}

func (t bpCreateResponse) RespondWithText() string {
	nl := ""
	for _, val := range t.Blueprints {
		nl += fmt.Sprintf("Name: %v\n", val.Name)
		nl += fmt.Sprintf("Content: %v\n", val.Content)
	}
	nl += "\n"
	fmt.Println(nl)
	return nl
}

func (t bpDescribeResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
}

func (t bpDescribeResponse) RespondWithJSON() string {
	nl := "{'definitions': ["

	for _, val := range t.Blueprints {
		//body := make(map[interface{}]interface{})
		var body interface{}
		if err := yaml.Unmarshal([]byte(val.Content), &body); err != nil {
			fmt.Errorf("Marshaling JSON response failed %v\n", err)
			return err.Error()
		}

		body = convert(body)
		jb, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		}
		nl += string(jb)
		nl += ","
		fmt.Println(nl)
	}
	nl = string(nl[:len(nl)-1])
	nl += "]}"

	return nl
}

func (t bpDescribeResponse) RespondWithText() string {
	nl := ""

	// Case: no blueprint found
	if len(t.Blueprints) == 0 {
		msg := "Blueprint not found"
		fmt.Println(msg)
		return msg
	}

	for _, val := range t.Blueprints {
		nl += fmt.Sprintf("%v\n", val.Content)
		nl += fmt.Sprintf("---\n") //replies contains multiple documents
	}
	nl = string(nl[:len(nl)-4])
	fmt.Println(nl)
	return nl
}

func (t bpListResponse) RespondWithText() string {
	nl := fmt.Sprintf("%-15v %-20v %-20v\n", "Blueprint Type", "Name", "Release Status")
	for _, val := range t.Blueprints {
		nl += fmt.Sprintf("%-15v %-20v %-15v\n", val.Type, val.Name, val.Availability)
	}
	fmt.Println(nl)
	return nl
}

func (t bpListResponse) RespondWithJSON() string {
	jb, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jb))

	return string(jb)
}

func (t bpListResponse) RespondWithYAML() string {

	yb, err := yaml.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(yb))

	return string(yb)
}

// bpClone create the blueprint repository using git clone
//   branch: specify a branch to use, default is latest release
//           can also be changed by setting PR_BLUEPRINT_BRANCH
func bpClone(branch string) error {

	// Initialize the blueprint cache
	tc, err := NewBlueprintCache()

	if err.errno == tcSuccess {
		// Found blueprint cache, git is initialized
		return nil
	} else if err.errno == tcBadBlueprintDirectory {
		// Unable to locate or create the desired blueprint cache directory
		return errors.New("Unable to find or create the cache directory")
	}

	// Create the blueprint cache and the .pr_cache file
	return tc.CreateCache(gitclone, branch)
}

//bpPull pulls blueprints from a remote repository
//  pullOptions
//    all: default
//    microservices:
//    serverless:
//    crd:
//  org: GitHub organization
//  repo: GitHub repository
//  path: path to start in repository
//  client: a github client
//
func bpPull(pullOptions, org, repo, path, outdir string,
	client *github.Client) error {

	opts := github.RepositoryContentGetOptions{}

	// Either file or directory content will be nil
	// file, director, resp, err
	fileContent, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, repo, path, &opts)
	if err != nil {
		log.Println("client.Repositories.GetContents: ", err)
		return err
	}

	// TODO: what is this
	//fmt.Println(rsp.StatusCode)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if fileContent != nil {
		dstr, _ := base64.StdEncoding.DecodeString(*fileContent.Content)
		fp := outdir + "/" + *fileContent.Path
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if _, err := os.Create(fp); err != nil {
				msg := fmt.Errorf("Creating file %s failed with error %s\n", fp, err)
				return msg
			}
		}
		err = ioutil.WriteFile(fp, dstr, 0600)
		if err != nil {
			log.Printf("ioutil.WriteFile error %v\n", err.Error())
		}

	} else {
		for _, item := range directoryContent {
			// If it is a directory, create if necessary
			// Then walk it
			if *item.Type == "dir" {
				dn := outdir + "/" + *item.Path
				if _, err := os.Stat(dn); os.IsNotExist(err) {
					if err := os.MkdirAll(dn, os.ModePerm); err != nil {
						fmt.Errorf("Failed creating directory %s with error %s\n",
							dn, err)
						return err
					}
					fmt.Println("Blueprint directory created: ", dn)
				}
				_ = bpPull(pullOptions, org, repo, *item.Path, outdir, client)
			}

			// For files, request their content
			if *item.Type == "file" {
				_ = bpPull(pullOptions, org, repo, *item.Path, outdir, client)
			}

		}

	}
	return nil
}

// bpCreate
//   Reads all the blueprints for a given type/name and then
//   extracts just the file names.  If not specified on the
//   command line, use the default definition.yaml file
//   associated this blueprint.
//
//   Next, it reads the definitions file and maps the data
//   attributes to the input structure, bpData, used by all blueprints
//
//   Before blueprints can be executed, the dynamic code components
//   must be compiled and also mapped to bpData for:
//     - Go struts for user defined objects
//     - SQL Scripts for developer use
//     - SQL functions for testing
//
//   Next, iterate over the blueprints to generate
//
//   Last, run goimport to check for missing or unused import
//   statements and clean up any code formatting issues
//
//func bpCreate(rn string)  bpCreateResponse {
func bpCreate(rn string) (reply bpCreateResponse) {
	var filelist []bpLocation
	var importList []string

	// Returns a list of all files in the blueprint directory
	// including there path, .datamgr/Makefile ....
	inputBlueprintFiles, err := bpRead(bpFile)

	// Read the definition file
	defs := bpDef{}

	//TODO: should be a method not a function
	err = bpReadDefinitions(&defs)
	if err != nil {
		var errorReply bpGenericResponseItem = bpGenericResponseItem{
			Name:    rn,
			Content: "Failed reading defintions file",
		}

		reply.Blueprints = append(reply.Blueprints, errorReply)
		return (reply)
	}

	lb := startBlocksSpinner("Loading blocks")
	bpInputData := bpData{}
	err = bpDataMapper(defs, &bpInputData)

	var newEndPoint endpointConfig
	if epList, err := newEndPoint.loadFromDefinitions(defs); err != nil {
		msg := fmt.Errorf("Endpoints warning: [%v]'\n", err)
		log.Println(msg)
	} else {
		for _, ep := range epList {
			// TODO: future support for other request routers
			routes, err := ep.GenerateBlock(GorillaRouteBlocks)
			importList = append(importList, GorillaRouteBlocks.getImports()...)
			if err != nil {
				fmt.Println(err)
			}

			methods, err := ep.GenerateBlock(GorillaMethodBlocks)
			importList = append(importList, GorillaMethodBlocks.getImports()...)
			incSpinner()
			if err != nil {
				fmt.Println(err)
			}

			bpInputData.EndpointRoutes += string(routes)
			bpInputData.EndpointHandlers += string(methods)
		}
	}

	lb.Stop()

	for _, projectBlocks := range defs.Project.Blocks {

		nb, err := projectBlocks.loadBlock(
			projectBlocks.ID,
			projectBlocks.Metadata.Labels)

		if err != nil {
			log.Println("Error: ", err)
		}

		if _, e := nb.GenerateBlock(defs); e != nil {
			log.Println("Error: ", e)
		}
	}

	b := Logger{}
	var loggerImports []string
	// reads a list of logger blocks and returns
	// required imports
	if loggerImports, err = b.getLoggerImports(defs); err != nil {
		msg := fmt.Errorf("Failed loading import: [%s]\n", err)
		log.Println(msg)
	}
	importList = append(importList, loggerImports...)

	bpInputData.GoImports = flattenUniqueStrings(importList)

	if len(defs.TableList) > 0 {
		SQLToJSONBlock.GenerateBlock(defs)
	}

	// Given the list returned in inputBlueprintFiles
	// Create a bpLocation object with the name and path
	// And a filtered list of blueprint files
	ebps := startBlocksSpinner("Load and execute blueprint")
	var filteredBlueprintList []string
	for _, rec := range inputBlueprintFiles {
		nm := filepath.Base(rec)
		di := rec[len(bpDirSelected)+1 : len(rec)-len(nm)]
		li := bpLocation{Name: nm, RelativePath: di}

		// This covers two special cases where we want to create a
		// directory with the name of the microservice.
		//
		//   - 1 If the file name is "blueprint", it represents a
		//       directory to created. Create it but don't add it
		//       to the list of blueprints to process.
		//       {blueprint manifests/kubernetes/dev/}
		//
		//   - 2 If the path contains "/blueprint/", replace the word
		//       "blueprint" with the name of the microservice and
		//       add it to the list of blueprints to process.
		//       Note: The pattern "/blueprint/" can occur multiple
		//       times.
		//       {blueprint-deployment.yaml manifests/kubernetes/dev/blueprint/}

		if nm == PREFIX {
			nm = bpInputData.Name
			dirName := di + nm
			_ = createDirectory(dirName)
			continue
		}

		testString1 := "/" + PREFIX + "/"
		if strings.Contains(di, testString1) {
			replaceString := "/" + bpInputData.Name + "/"
			li.RelativePath = strings.ReplaceAll(li.RelativePath,
				testString1, replaceString)
		}

		// starts with blueprint"/"
		testString2 := PREFIX + "/"
		if strings.HasPrefix(di, testString2) {
			replaceString := bpInputData.Name + "/"
			li.RelativePath = strings.Replace(li.RelativePath,
				testString2, replaceString, 1)
		}

		// TODO: we need the equivalent of a .gitignore file
		//       and we can move this logic into the skipExistingHookFile
		//       function
		isHook := strings.Contains(strings.ToLower(li.Name),
			strings.ToLower(HOOK))
		if isHook {
			testName := strings.Replace(li.Name, PREFIX, bpInputData.Name, 1)
			var fileWithPath string
			if li.RelativePath == "" {
				fileWithPath = testName
			} else {
				fileWithPath = li.RelativePath + "/" + testName
			}
			e := skipExistingHookFile(fileWithPath)
			if e {
				//fmt.Printf("Skipping processing of hook file %v\n", fileWithPath)
				continue
			}
		}

		filelist = append(filelist, li)
		filteredBlueprintList = append(filteredBlueprintList, rec)
	}

	// Generate internal structures
	if len(defs.TableList) > 0 {
		err = bpGenerateStructurs(defs, &bpInputData)
		if err != nil {
			fmt.Println("Generating structures failed: ", err)
			os.Exit(-1)
		}

		// Generate JSON test data
		err = bpJSONData(defs, &bpInputData)
		if err != nil {
			fmt.Println("Generating JSON failed: ", err)
			os.Exit(-1)
		}

	}

	// Build the blueprint cache
	//blueprints, err = templates.New("").ParseFiles(inputBlueprintFiles...)
	blueprints, err = template.New("").ParseFiles(filteredBlueprintList...)
	if err != nil {
		fmt.Println("Blueprint parsing failed: ", err)
		os.Exit(-1)
	}

	//blueprints = blueprint.Must(blueprint.ParseFiles(inputBlueprintFiles...))
	//TODO: turn into a function
	var fn string
	incSpinner()
	//	fmt.Printf("Executing  %s blueprint for new service %s\n", defs.Info.ID, defs.Info.Name)
	for _, v := range filelist {
		if v.Name == bpDefinition {
			continue
		}
		// Replace generic string "blueprint" with the name of the service
		// If it is not in the name, the unmodified file name is used
		if strings.HasPrefix(v.Name, PREFIX) {
			fn = strings.Replace(v.Name, PREFIX, bpInputData.Name, 1)
		} else if strings.HasPrefix(v.Name, ORGANIZATION) {
			fn = strings.Replace(v.Name, ORGANIZATION, bpInputData.Organization, 1)
		} else {
			fn = v.Name
		}

		// Ensure the path to the file exists
		if v.RelativePath != "" {
			if _, err := os.Stat(v.RelativePath); os.IsNotExist(err) {
				err := os.MkdirAll(v.RelativePath, 0750)
				if err != nil {
					fmt.Println("Failed to make directory: ", v.RelativePath)
				}
			}
		}

		var defaultMode os.FileMode = 0660
		// Make sure shell scripts are created executable
		if strings.Contains(fn, ".sh") {
			defaultMode = 0755
		}

		file, err := os.OpenFile(v.RelativePath+fn,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, defaultMode)
		if err != nil {
			log.Fatal(err, v.RelativePath+fn)
		}

		bw := bufio.NewWriter(file)

		//fmt.Printf("executing %v writing to %v\n", fn, v.RelativePath+fn)
		rawDirectoryName := toInputBlueprintName(v.RelativePath, bpInputData.Name)
		err = blueprints.ExecuteTemplate(bw, rawDirectoryName+v.Name, bpInputData)
		if err != nil {
			fmt.Printf("Blueprint execution failed for: %v with error %v", v, err)
			os.Exit(-1)
		}
		if err := bw.Flush(); err != nil {
			fmt.Errorf("Flush for file %s failed with error %v\n",
				file.Name(), err)
			os.Exit(-1)
		}
		if err := file.Close(); err != nil {
			fmt.Errorf("Close for file %s failed with error %v\n",
				file.Name(), err)
			os.Exit(-1)
		}
	}

	ebps.Stop()

	var successReply bpGenericResponseItem = bpGenericResponseItem{
		Name: rn,
		Content: `Create succeeded run make to compile your applications
$ make <CR>`,
	}
	reply.Blueprints = append(reply.Blueprints, successReply)

	return reply
}

// toInputBlueprintName map a directory path to its name
// in the input blueprint directory
func toInputBlueprintName(path, name string) string {
	testString1 := "/" + name + "/"
	testString2 := name + "/"
	result := path

	if strings.Contains(path, testString1) {
		replaceString := "/" + PREFIX + "/"
		result = strings.ReplaceAll(path,
			testString1, replaceString)
	}

	if strings.HasPrefix(path, testString2) {
		replaceString := PREFIX + "/"
		result = strings.Replace(path, testString2, replaceString, 1)
	}

	return result
}

// bpEdit
//   use the $EDITOR envar to edit the named resource
//
func bpEdit(name string) {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}

	path, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal(editor, "not found")
	}

	fmt.Println(path, name)
	cmd := exec.Command(editor, name)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	}

	return
}

// bpReadDefinitions
//   Read the definition file and then validate it
//
func bpReadDefinitions(definitionsStruct *bpDef) error {

	msg := fmt.Sprintf("Reading definitions from: %s", bpDefFile)
	s := startBlocksSpinner(msg)
	df, err := os.Open(bpDefFile)
	if err != nil {
		fmt.Println("failed to open:", bpDefFile, ", error:", err)
	}
	defer df.Close()
	incSpinner()
	byteValue, e := ioutil.ReadAll(df)
	if e != nil {
		fmt.Println("read failed for ", bpDefFile)
		os.Exit(-1)
	}

	incSpinner()
	err = yaml.Unmarshal([]byte(byteValue), definitionsStruct)
	if err != nil {
		fmt.Println("Unmarshal faild", err)
		return err
	}

	definitionsStruct.DefinitionFile = bpDefFile
	definitionsStruct.DefinitionFileVersion = StaticDefinitionFileVersion

	incSpinner()
	errs := definitionsStruct.Validate()

	if errs > 0 {
		fmt.Println("definitionsStruct.Validate()")
		for _, item := range ErrList {
			fmt.Println(item.Error())
		}
		s.Stop()

		os.Exit(-1)
	}
	s.Stop()
	return nil
}

// bpDescribe Get default blueprint definitions file
func bpDescribe(bpListOption string, rn string) bpDescribeResponse {
	var response bpDescribeResponse

	// Get the list of blueprints
	rsp := bpGet("all", rn)

	// Load the definitions.yaml file in the blueprint directory
	for _, item := range rsp.Blueprints {
		fn := item.Path + "/" + item.Name + "/" + "definition.yaml"
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			continue
		}

		jf, err := os.Open(fn)
		if err != nil {
			fmt.Printf("failed to open: %v err %v", fn, err)
			continue
		}

		defer jf.Close()

		byteValue, _ := ioutil.ReadAll(jf)
		/*
				Type         string // Type of blueprint, i.e. Server less
				Name         string // Name of blueprint == directory name
			  Content      string // YAML configuration data
		*/
		nItem := bpDescribeItem{item.Type, item.Name, string(byteValue)}
		response.Blueprints = append(response.Blueprints, nItem)
	}

	return response
}

// bpExplain
//   is stored in docs/
//   There is a name.txt file for each blueprint
//
//   docs is defined by eTLD
func bpExplain(bpListOption string, rn string) bpExplainResponse {
	var response bpExplainResponse

	// read/check blueprint cache
	tc, te := NewBlueprintCache()
	if te.errno != tcSuccess {
		log.Fatalf("Failed to read blueprint cache, Got (%v)\n", te)
	}

	// Load explanation
	fn := tc.location.Location() + "/" + eTLD + "/" + rn + ".txt"
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return response
	}

	tf, err := os.Open(fn)
	if err != nil {
		fmt.Printf("failed to open: %v err %v", tf, err)
		return response
	}
	defer tf.Close()

	byteValue, _ := ioutil.ReadAll(tf)
	nItem := bpExplainItem{bpResourceName, string(byteValue)}
	response.Blueprints = append(response.Blueprints, nItem)

	return response
}

// bpRead(name)
//   read all files for the named blueprint
//   return a slice of paths with file names
//   "bpDir/TLD/SLD/blueprintFiles......."
func bpRead(bpName string) ([]string, error) {
	var bpFlLst []string

	bpRsp := bpGet(bpResourceName, bpFile)
	if len(bpRsp.Blueprints) == 0 {
		em := errors.New("Blueprint " + bpName + " not found")
		return bpFlLst, em
	}

	// TODO: allow namespace to include ga, .....
	if len(bpRsp.Blueprints) > 1 {
		em := errors.New("Error: Blueprint " + bpName + " not unique")
		return bpFlLst, em
	}

	bpItem := bpRsp.Blueprints[0]
	td := bpItem.Path + "/" + bpItem.Name

	bpDirSelected = td
	err := filepath.Walk(td,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Don't include directories, just their contents
			if info.IsDir() == false {
				//Don't include backup files
				if len(strings.Split(path, "~")) == 1 {
					bpFlLst = append(bpFlLst, path)
				}
			}

			// Special case where the directory name is blueprint
			if info.IsDir() == true && strings.HasSuffix(path, PREFIX) {
				bpFlLst = append(bpFlLst, path)
			}

			return nil
		})

	if err != nil {
		em := errors.New("Error: Unable to walk directory" + td)
		return bpFlLst, em
	}

	return bpFlLst, nil
}

// List available blueprints
//  bpListOption: TBD
//  rn: resource name if specified on the command line
func bpGet(bpListOption string, rn string) bpListResponse {
	bpTLD := []string{"crd", "microservices", "serverless"}
	bpSLD := []string{"ga", "experimental", "incubation"}
	var response bpListResponse
	tc, err := NewBlueprintCache()
	if err.errno != tcSuccess {
		log.Fatalf("Failed to read blueprint cache, Got (%v)\n", err)
	}

	for _, tld := range bpTLD {
		for _, sld := range bpSLD {
			dn := tc.location.Location() + "/" + tld + "/" + sld
			if _, err := os.Stat(dn); os.IsNotExist(err) {
				continue
			}
			f, err := os.Open(dn)
			if err != nil {
				continue
			}

			list, err := f.Readdir(-1)
			defer f.Close()

			if err != nil {
				continue
			}

			for _, fn := range list {
				nrec := bpListItem{tld, sld, fn.Name(), dn}
				// Skip empty directories initialized with a .nothing file
				if fn.Name() != ".nothing" {
					if rn != "" && fn.Name() != rn {
						//n is defined skip records that don't match
						continue
					}
					response.Blueprints = append(response.Blueprints, nrec)
				}
			}
		}
	}
	return response
}

func createDirectory(path string) error {
	var defaultMode os.FileMode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, defaultMode)
		if err != nil {
			fmt.Println("failed creating directory: ", path)
			return err
		}
	}

	return nil
}

func skipExistingHookFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

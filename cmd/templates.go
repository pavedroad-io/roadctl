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
	"time"

	"github.com/google/go-github/github"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

// Default template repository
var defaultOrg = "pavedroad-io"

// Default template directory on GitHub
var defaultRepo = "templates"
var defaultPath string

// Default template directory on local machine
var defaultTemplateDir = ".templates"
var repoType = "GitHub"
var tplFile string    // Name of the template file to use
var tplDir = "."      // Directory for generated code output
var tplDefFile string // Name of the definitions file used to generated templates
var tplDirSelected string

// The release branch stores released templates
// The latest and stable tags are used to select which release
const (
	gitTemplateBranch  = "release"
	gitLatestTag       = "latest"
	gitStableTag       = "stable"
	templateRepository = "https://github.com/pavedroad-io/templates"
	githubAPI          = "GitHub API"
	gitclone           = "git clone"
)

// TEMPLATE needs documentation.
const (
	tplResourceName = "templates"
	tplDefinition   = "definition.yaml"
	// TEMPLATE is the prefix to be replaced in front of the file name
	TEMPLATE = "template"
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

// Makefile constants
const (
	allWithFossa    string = "all: $(PREFLIGHT) $(FOSSATEST) compile check"
	allWithoutFossa string = "all: $(PREFLIGHT) compile check"

	checkWithSonar    string = "check: lint sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)"
	checkWithoutSonar string = "check: lint $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)"

	// Fossa has a build section and a lint section
	fossaSection string = `
$(FOSSATEST):
	fossa init
`
	fossaLint string = `
	@echo "  >  running FOSSA license scan."
	@FOSSA_API_KEY=$(FOSSA_API_KEY) fossa analyze
`
)

var templates *template.Template

// tplDirectory manages template directory locations
type tplDirectory struct {
	// Full path to the template directory
	location string

	// Is it initialized
	initialized bool

	// How we determined the location
	// default, so.GetEnv, command line option
	locationFrom string
}

const (
	tplCacheCreatedFile string = ".pr_cache_created"
	tplCacheUreatedFile string = ".pr_cache_updated"
	tplCacheCreatedWtih string = ".pr_cache_created_with"
)

// tplCache manages information about templates
//  stored in a template directory
type tplCache struct {
	// What directory is it in
	location *tplDirectory

	// Is it initialized
	initialized bool

	// git clone or github API
	initilazedFrom string

	// Track we we creatd and last updated
	created time.Time
	updated time.Time
}

type tplData struct {
	// Information about company and project
	Organization        string // Name of Organization
	OrgSQLSafe          string // Mapped to a safe name for using in SQL
	OrganazationInfo    string // Name of Organization
	OrganizationLicense string // Org license/copyright
	ProjectInfo         string // Project/service description
	SchedulerName       string // For worker polls specifies the type of
	//   scheduler to create

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

	// Integrations
	Badges            string // badges to include docs
	SonarKey          string
	SonarLogin        string
	SonarPrefix       string
	SonarCloudEnabled bool
	FOSSAEnabled      bool

	// Service and tpl-names
	Name         string //service name
	NameExported string //camel case with first letter cap
	TplName      string //template name
	DefFile      string //definition file used

	//PR license/copyright should be a function
	PavedroadInfo string //PR license/copyright

	//Swagger headers probably turn these into functions
	AllRoutesSwaggerDoc     string
	GetAllSwaggerDoc        string // swagger for list method
	GetSwaggerDoc           string // swagger for get method
	PutSwaggerDoc           string // swagger for put method
	PostSwaggerDoc          string // swagger for post method
	DeleteSwaggerDoc        string // swagger for delete method
	SwaggerGeneratedStructs string // swagger doc and go struct
	DumpStructs             string // Generic dump of given object type

	//JSON data
	PostJSON string // Sample data for a post
	PutJSON  string // Sample data for a put

	// Makefile options
	CheckBuildTarget  string //build line for check section
	AllBuildTarget    string //build line for check section
	FossaBuildSection string //build target for Fossa
	FossaLintSection  string //lint section for Fossa

}

//  tplDataMapper
//    Map data from definitions file to tplData structure
//    return error if required mappings are missing
//    TODO: jms
func tplDataMapper(defs tplDef, output *tplData) error {
	// Docker images names don't allow uppercase letters
	output.Name = strings.ToLower(defs.Info.Name)
	output.NameExported = strcase.ToCamel(defs.Info.Name)
	output.TplName = defs.Info.ID
	output.DefFile = tplDefFile
	output.OrganizationLicense = defs.Project.License
	output.Organization = defs.Info.Organization

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
	output.Badges = defs.BadgesToString()

	//Sonarcloud
	si := defs.findIntegration("sonarcloud")

	if si.Name != "" {
		output.SonarKey = si.SonarCloudConfig.Key
		output.SonarLogin = si.SonarCloudConfig.Login
		output.SonarPrefix = SONARPREFIX
		output.SonarCloudEnabled = si.Enabled
	}

	if output.SonarCloudEnabled {
		output.CheckBuildTarget = checkWithSonar
	} else {
		output.CheckBuildTarget = checkWithoutSonar
	}

	si = defs.findIntegration("fossa")
	if si.Name != "" {
		output.FOSSAEnabled = si.Enabled
	}

	if output.FOSSAEnabled {
		output.FossaBuildSection = fossaSection
		output.FossaLintSection = fossaLint
		output.AllBuildTarget = allWithFossa
	} else {
		output.FossaBuildSection = ""
		output.FossaLintSection = ""
		output.AllBuildTarget = allWithoutFossa
	}

	return nil
}

//  tplJSONData
//    Use the schema definition found in tplDefs to create
//    Sample JSON data files
//
func tplJSONData(defs tplDef, output *tplData) error {
	var jsonString string
	order := defs.devineOrder()
	tplAddJSON(order, defs, &jsonString)

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

// tplAddJSON
//
//   Creates JSON sample data
//
func tplAddJSON(item tplTableItem, defs tplDef, jsonString *string) {
	table, _ := defs.tableByName(item.Name)

	// Start this table
	if item.Root {
		*jsonString = fmt.Sprintf(jsonObjectStart)
	} else {
		*jsonString += fmt.Sprintf(jsonField, strings.ToLower(item.Name))
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
	for idx, col := range table.Columns {
		// Add it to the dynamic struct
		var sample interface{}
		switch col.Type {
		case "string":
			sample = RandomString(15)
		case "int", "integer", "int32", "int64":
			sample = RandomInteger(0, 254)
		case "number", "float", "float32", "float64":
			sample = RandomFloat()
		case "bool":
			sample = RandomBool()
		case "time":
			sample = time.Now().Format(time.RFC3339)
		}

		*jsonString += fmt.Sprintf(jsonField, strings.ToLower(col.Name))
		*jsonString += fmt.Sprintf(jsonValue, sample)
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
			tplAddJSON(*child, defs, jsonString)
		}
	}

	// Close and append to tplData.SwaggerGeneratedStructs
	*jsonString += jsonObjectEnd

	return
}

//  tplGenerateStructurs
//    Use the schema definition found in tplDefs to create
//    Go structures and assign to tplData.SwaggerGeneratedStructs
//
//    Use the same schema to generate a formated dump
//    command to aid developer debugging and assign it to
//    tplData.DumpStructs
//
func tplGenerateStructurs(defs tplDef, output *tplData) error {
	order := defs.devineOrder()
	tplAddStruct(order, defs, output)
	return nil
}

// tplAddStruct
//
// Performs two tasks
//    - 1 Generates the structure as a string that is inserted
//        into the code template.  This is the "tableString"
//        variable
//
//    - 2 Creates JSON sample data
//        One for insert, and one for updates
//
func tplAddStruct(item tplTableItem, defs tplDef, output *tplData) {
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
			tplAddStruct(*child, defs, output)

			// Same as structField except type with be the suitable
			tableString += fmt.Sprintf(structSubstruct,
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
		var fieldType string
		if col.Type == "time" {
			fieldType = "time.Time"
		} else {
			fieldType = strings.ToLower(col.Type)
		}

		//Deal with time types
		tableString += fmt.Sprintf(structField,
			strcase.ToCamel(col.Name),
			fieldType,
			"json",
			importLine)

	}

	// Close and append to tplData.SwaggerGeneratedStructs
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

// tplListItem provides information about a template location
// and status
type tplListItem struct {
	Type         string // Type of template, i.e. serverless
	Availability string // Availability ga, ....
	Name         string // Name of template == directory name
	Path         string // Path to the template
}

// tplLocation
type tplLocation struct {
	Name         string //Name of the template file
	RelativePath string // Path relative to the current directory
}

type tplExplainItem struct {
	Name    string // Name of resource
	Content string // Text for explain document
}

type tplDescribeItem struct {
	Type    string // Type of template, i.e. serverless
	Name    string // Name of template == directory name
	Content string // YAML configuration data
}

type tplExplainResponse struct {
	Templates []tplExplainItem
}

type tplDescribeResponse struct {
	Templates []tplDescribeItem
}

type tplListResponse struct {
	Templates []tplListItem
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

func (t tplExplainResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
}

func (t tplExplainResponse) RespondWithJSON() string {
	return t.RespondWithText() // One in the same for this type
}

func (t tplExplainResponse) RespondWithText() string {
	nl := ""
	for _, val := range t.Templates {
		nl += fmt.Sprintf("Name: %v\n", val.Name)
		nl += fmt.Sprintf("%v\n", val.Content)
	}
	nl += "\n"
	fmt.Println(nl)
	return nl
}

func (t tplDescribeResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
}

func (t tplDescribeResponse) RespondWithJSON() string {
	nl := "{'definitions': ["

	for _, val := range t.Templates {
		//body := make(map[interface{}]interface{})
		var body interface{}
		yaml.Unmarshal([]byte(val.Content), &body)

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

func (t tplDescribeResponse) RespondWithText() string {
	nl := ""
	for _, val := range t.Templates {
		nl += fmt.Sprintf("%v\n", val.Content)
		nl += fmt.Sprintf("---\n") //replies contains multiple documents
	}
	nl = string(nl[:len(nl)-4])
	fmt.Println(nl)
	return nl
}

func (t tplListResponse) RespondWithText() string {
	nl := fmt.Sprintf("%-15v %-20v %-20v\n", "Template Type", "Name", "Release Status")
	for _, val := range t.Templates {
		nl += fmt.Sprintf("%-15v %-20v %-15v\n", val.Type, val.Name, val.Availability)
	}
	fmt.Println(nl)
	return nl
}

func (t tplListResponse) RespondWithJSON() string {
	jb, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jb))

	return string(jb)
}

func (t tplListResponse) RespondWithYAML() string {

	yb, err := yaml.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(yb))

	return string(yb)
}

// Location returns the location of the template directory
// Initialize if necessary
func (t *tplDirectory) Location() string {

	if !t.initialized {
		err := t.initialize()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return t.location
}

// initialize a private function for initializing
//   the template directory location, not the
//   templates
func (t *tplDirectory) initialize() error {
	// Order of precedence
	//   - roadctl CLI
	//   - PR_TEMPLATE_DIR
	//   - defaultTemplateDir

	env := os.Getenv("PR_TEMPLATE_DIR")
	if templateDirectoryLocation != "" {
		t.location = templateDirectoryLocation
		t.locationFrom = "CLI"
	} else if env != "" && templateDirectoryLocation == "" {
		t.location = env
		t.locationFrom = "PR_TEMPLATE_DIR"
	} else {
		t.location = defaultTemplateDir
		t.locationFrom = "default"
	}

	if !t.initialized {
		if err := createDirectory(t.location); err != nil {
			log.Fatal(err.Error())
		}
		t.initialized = true
	}

	return nil
}

func (t *tplDirectory) getDefault() string {

	return ""
}

// New create a tplCache
// If it does not exists, initalize it using method specified
//   td: a tplDirectory type
//   method: GitHub API or git clone
func (tc *tplCache) New(td *tplDirectory, method string) error {
	tc.location = td
	switch method {
	case gitclone:
		//tc.Clone()
	case githubAPI:
		//tc.API()
	default:
		fmt.Println("error")
	}
	return nil
}

// tplClone create the template repository using git clone
//   branch: specify a branch to use, default is latest release
//           can also be changed by setting PR_TEMPLATE_BRANCH
func tplClone(branch string) error {

	t := &tplDirectory{}
	tc := &tplCache{}
	if dir := t.Location(); dir != "" {
		msg := fmt.Sprintf("Unable to open template directory (%v)\n", dir)
		return errors.New(msg)
	}
	return tc.New(t, gitclone)
}

//tplPull pulls templates from a remote repository
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
func tplPull(pullOptions, org, repo, path, outdir string,
	client *github.Client) error {

	opts := github.RepositoryContentGetOptions{}

	// Either file or directory content will be nil
	// file, director, resp, err
	fileContent, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, repo, path, &opts)
	if err != nil {
		//TODO: change to proper logging method
		fmt.Println("client.Repositories.GetContents: ", err)
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
			os.Create(fp)
		}
		err = ioutil.WriteFile(fp, dstr, 0644)
		if err != nil {
			fmt.Printf("ioutil.WriteFile error %v\n", err.Error())
		}

	} else {
		for _, item := range directoryContent {
			// If it is a directory, create if necessary
			// Then walk it
			if *item.Type == "dir" {
				dn := outdir + "/" + *item.Path
				if _, err := os.Stat(dn); os.IsNotExist(err) {
					os.MkdirAll(dn, os.ModePerm)
					fmt.Println("Template directory created: ", dn)
				}
				_ = tplPull(pullOptions, org, repo, *item.Path, outdir, client)
			}

			// For files, request their content
			if *item.Type == "file" {
				_ = tplPull(pullOptions, org, repo, *item.Path, outdir, client)
			}

		}

	}
	return nil
}

// tplCreate
//   Reads all the templates for a given type/name and then
//   extracts just the file names.  If not specified on the
//   command line, use the default definition.yaml file
//   associated this template.
//
//   Next, it reads the definitions file and maps the data
//   attributes to the input structure, tplData, used by all templates
//
//   Before templates can be executed, the dynamic code components
//   must be compiled and also mapped to tplData for:
//     - Go struts for user defined objects
//     - SQL Scripts for developer use
//     - SQL functions for testing
//
//   Next, iterate over the templates to generate
//
//   Last, run goimport to check for missing or unused import
//   statements and clean up any code formatting issues
//
func tplCreate(rn string) string {
	var filelist []tplLocation
	//var filenames []string

	// Returns a list of all files in the template directory
	// including there path, .datamgr/Makefile ....
	inputTemplateFiles, err := tplRead(tplFile)
	if err != nil {
		fmt.Println(err)
	}

	// Read the definition file
	defs := tplDef{}

	//TODO: should be a method not a function
	err = tplReadDefinitions(&defs)
	if err != nil {
		fmt.Println(err)
		return (err.Error())
	}

	tplInputData := tplData{}
	err = tplDataMapper(defs, &tplInputData)

	// Given the list returned in inputTemplateFiles
	// Create a tplLocation object with the name and path
	// And a filtered list of template files
	var filteredTemplateList []string
	for _, rec := range inputTemplateFiles {
		nm := filepath.Base(rec)
		di := rec[len(tplDirSelected)+1 : len(rec)-len(nm)]
		li := tplLocation{Name: nm, RelativePath: di}

		// This covers two special cases where we want to create a
		// directory with the name of the microservice.
		//
		//   - 1 If the file name is "template", it represents a
		//       directory to created. Create it but don't add it
		//       to the list of templates to process.
		//       {template manifests/kubernetes/dev/}
		//
		//   - 2 If the path contains "/template/", replace the word
		//       "template" with the name of the microservice and
		//       add it to the list of templates to process.
		//       Note: The pattern "/template/" can occur multiple
		//       times.
		//       {template-deployment.yaml manifests/kubernetes/dev/template/}

		if nm == TEMPLATE {
			nm = tplInputData.Name
			dirName := di + nm
			_ = createDirectory(dirName)
			continue
		}

		testString1 := "/" + TEMPLATE + "/"
		if strings.Contains(di, testString1) {
			replaceString := "/" + tplInputData.Name + "/"
			li.RelativePath = strings.ReplaceAll(li.RelativePath,
				testString1, replaceString)
		}

		// starts with template"/"
		testString2 := TEMPLATE + "/"
		if strings.HasPrefix(di, testString2) {
			replaceString := tplInputData.Name + "/"
			li.RelativePath = strings.Replace(li.RelativePath,
				testString2, replaceString, 1)
		}

		// TODO: we need the equivalent of a .gitignore file
		//       and we can move this logic into the skipExistingHookFile
		//       function
		isHook := strings.Contains(strings.ToLower(li.Name),
			strings.ToLower(HOOK))
		if isHook {
			testName := strings.Replace(li.Name, TEMPLATE, tplInputData.Name, 1)
			var fileWithPath string
			if li.RelativePath == "" {
				fileWithPath = testName
			} else {
				fileWithPath = li.RelativePath + "/" + testName
			}
			e := skipExistingHookFile(fileWithPath)
			if e {
				fmt.Printf("Skipping processing of hook file %v\n", fileWithPath)
				continue
			}
		}

		filelist = append(filelist, li)
		filteredTemplateList = append(filteredTemplateList, rec)
	}

	// If sonarcloud is configured validate token and project
	if tplInputData.SonarCloudEnabled {
		err = validateIntegrations(&tplInputData)
	}

	if err != nil {
		fmt.Println("Validating integrations failed: ", err)
		os.Exit(-1)
	}

	// Generate internal structures
	if len(defs.TableList) > 0 {
		err = tplGenerateStructurs(defs, &tplInputData)
		if err != nil {
			fmt.Println("Generating structures failed: ", err)
			os.Exit(-1)
		}

		// Generate JSON test data
		err = tplJSONData(defs, &tplInputData)
		if err != nil {
			fmt.Println("Generating JSON failed: ", err)
			os.Exit(-1)
		}

	}

	// Build the template cache
	//templates, err = template.New("").ParseFiles(inputTemplateFiles...)
	templates, err = template.New("").ParseFiles(filteredTemplateList...)
	if err != nil {
		fmt.Println("Template parsing failed: ", err)
		os.Exit(-1)
	}

	//templates = template.Must(template.ParseFiles(inputTemplateFiles...))
	//TODO: turn into a function
	var fn string
	for _, v := range filelist {
		if v.Name == tplDefinition {
			continue
		}
		// Replace generic string "template" with the name of the service
		// If it is not in the name, the unmodified file name is used
		if strings.HasPrefix(v.Name, TEMPLATE) {
			fn = strings.Replace(v.Name, TEMPLATE, tplInputData.Name, 1)
		} else if strings.HasPrefix(v.Name, ORGANIZATION) {
			fn = strings.Replace(v.Name, ORGANIZATION, tplInputData.Organization, 1)
		} else {
			fn = v.Name
		}

		// Ensure the path to the file exists
		if v.RelativePath != "" {
			if _, err := os.Stat(v.RelativePath); os.IsNotExist(err) {
				err := os.MkdirAll(v.RelativePath, 0766)
				if err != nil {
					fmt.Println("Failed to make directory: ", v.RelativePath)
				}
			}
		}

		var defaultMode os.FileMode = 0666
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

		fmt.Printf("executing %v writing to %v\n", fn, v.RelativePath+fn)
		rawDirectoryName := toInputTemplateName(v.RelativePath, tplInputData.Name)
		err = templates.ExecuteTemplate(bw, rawDirectoryName+v.Name, tplInputData)
		if err != nil {
			fmt.Printf("Template execution failed for: %v with error %v", v, err)
			os.Exit(-1)
		}
		bw.Flush()
		file.Close()
	}

	// Execute goimport for code formatting

	return ""
}

// toInputTemplateName map a directory path to its name
// in the input template directory
func toInputTemplateName(path, name string) string {
	testString1 := "/" + name + "/"
	testString2 := name + "/"
	result := path

	if strings.Contains(path, testString1) {
		replaceString := "/" + TEMPLATE + "/"
		result = strings.ReplaceAll(path,
			testString1, replaceString)
	}

	if strings.HasPrefix(path, testString2) {
		replaceString := TEMPLATE + "/"
		result = strings.Replace(path, testString2, replaceString, 1)
	}

	return result
}

// tplEdit
//   use the $EDITOR envar to edit the named resource
//
func tplEdit(name string) {
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

// tplReadDefinitions
//   Read the definition file and then validate it
//
func tplReadDefinitions(definitionsStruct *tplDef) error {

	fmt.Println("Reading definitions from: ", tplDefFile)
	df, err := os.Open(tplDefFile)
	if err != nil {
		fmt.Println("failed to open:", tplDefFile, ", error:", err)
	}
	defer df.Close()
	byteValue, e := ioutil.ReadAll(df)
	if e != nil {
		fmt.Println("read failed for ", tplDefFile)
		os.Exit(-1)
	}

	err = yaml.Unmarshal([]byte(byteValue), definitionsStruct)
	if err != nil {
		return err
	}

	errs := definitionsStruct.Validate()

	if errs != nil {
		fmt.Println("definitionsStruct.Validate()")
		for errs != nil {
			fmt.Println(errs.Error())
			errs = errs.nextError
		}
		os.Exit(-1)
	}
	return nil
}

// tplDescribe Get default template definitions file
func tplDescribe(tplListOption string, rn string) tplDescribeResponse {
	var response tplDescribeResponse

	// Get the list of templates
	rsp := tplGet("all", rn)

	// Load the definitions.yaml file in the template directory
	for _, item := range rsp.Templates {
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
				Type         string // Type of template, i.e. Server less
				Name         string // Name of template == directory name
			  Content      string // YAML configuration data
		*/
		nItem := tplDescribeItem{item.Type, item.Name, string(byteValue)}
		response.Templates = append(response.Templates, nItem)
	}

	return response
}

// tplExplain
//   is stored in docs/explain.txt
//   docs is defined by eTLD
//   explain is tplResourceName
func tplExplain(tplListOption string, rn string) tplExplainResponse {
	var response tplExplainResponse

	// Load explanation
	fn := defaultTemplateDir + "/" + eTLD + "/" + tplResourceName + ".txt"
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
	fmt.Println(string(fn))
	nItem := tplExplainItem{tplResourceName, string(byteValue)}
	response.Templates = append(response.Templates, nItem)

	return response
}

// tplRead(name)
//   read all files for the named template
//   return a slice of paths with file names
//   "tplDir/TLD/SLD/templateFiles......."
func tplRead(tplName string) ([]string, error) {
	var tplFlLst []string

	tplRsp := tplGet(tplResourceName, tplFile)
	if len(tplRsp.Templates) == 0 {
		em := errors.New("Template " + tplName + " not found")
		return tplFlLst, em
	}

	// TODO: allow namespace to include ga, .....
	if len(tplRsp.Templates) > 1 {
		em := errors.New("Error: Template " + tplName + " not unique")
		return tplFlLst, em
	}

	tplItem := tplRsp.Templates[0]
	td := tplItem.Path + "/" + tplItem.Name

	fmt.Println("Tpl dir", td)
	tplDirSelected = td
	err := filepath.Walk(td,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Don't include directories, just their contents
			if info.IsDir() == false {
				//Don't include backup files
				if len(strings.Split(path, "~")) == 1 {
					tplFlLst = append(tplFlLst, path)
				}
			}

			// Special case where the directory name is template
			if info.IsDir() == true && strings.HasSuffix(path, TEMPLATE) {
				fmt.Println(path)
				tplFlLst = append(tplFlLst, path)
			}

			return nil
		})

	if err != nil {
		em := errors.New("Error: Unable to walk directory" + td)
		return tplFlLst, em
	}

	return tplFlLst, nil
}

// List available templates
//  tplListOption: TBD
//  rn: resource name if specified on the command line
func tplGet(tplListOption string, rn string) tplListResponse {
	tplTLD := []string{"crd", "microservices", "serverless"}
	tplSLD := []string{"ga", "experimental", "incubation"}
	var response tplListResponse

	for _, tld := range tplTLD {
		for _, sld := range tplSLD {
			dn := defaultTemplateDir + "/" + tld + "/" + sld
			if _, err := os.Stat(dn); os.IsNotExist(err) {
				continue
			}
			f, err := os.Open(dn)
			if err != nil {
				continue
			}

			list, err := f.Readdir(-1)
			f.Close()

			if err != nil {
				continue
			}

			for _, fn := range list {
				nrec := tplListItem{tld, sld, fn.Name(), dn}
				// Skip empty directories initialized with a .nothing file
				if fn.Name() != ".nothing" {
					if rn != "" && fn.Name() != rn {
						//n is defined skip records that don't match
						continue
					}
					response.Templates = append(response.Templates, nrec)
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
	} else {
		return true
	}
}

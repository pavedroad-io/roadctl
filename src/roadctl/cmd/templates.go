/*
Copyright © 2019 PavedRoad

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

package cmd

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
	"path/filepath"
	// "reflect"
	"github.com/google/go-github/github"
	"github.com/iancoleman/strcase"
	"strings"
	"text/template"
	"time"
	// "github.com/ompluscator/dynamic-struct"
	"gopkg.in/yaml.v2"
)

// Default template repository
var defaultOrg = "pavedroad-io"

// Default template directory on GitHub
var defaultRepo = "templates"
var defaultPath = ""

// Default template directory on local machine
var defaultTemplateDir = ".templates"
var repoType = "GitHub"
var tplFile = ""      // Name of the template file to use
var tplDir = "."      // Directory for generated code output
var tplDefFile string // Name of the definitions file used to generated tempaltes
var tplDirSelected = ""

//TEMPLATE needs documentation.
const (
	tplResourceName = "templates"
	tplDefinition   = "definition.yaml"
	TEMPLATE        = "template"
	ORGANIZATION    = "organization"
	prCopyright     = `
//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
//`
)
const strcutComment = `
//
//
//`

// JSON formaters
const (
	jsonObjectStart = "{\n"
	jsonObjectEnd   = "}"
	jsonListStart   = "["
	jsonListEnd     = "]"
	jsonSeperator   = ",\n"
	jsonField       = "\"%v\": "
	jsonValue       = "\"%v\"" // If new object, or last field strip the comma
)

const swaggerRoute = "// swagger:response %s\n"
const structOpen = "type %s struct {\n"

// structUUID
// name of table, type of data json|yaml
const structUUID = "\t%sUUID string `%s:%suuid`\n"
const structClose = "}\n\n"

// structField
//  name, type, encoding, json|yaml, encoding options
const structField = "\t%s %s\t`%s:%s`\n"

// structSubStruct
// Same as structField except:
//  type will be the subtable
//  No options
const structSubstruct = "\t%s %s\t`%s:%s`\n"

var templates *template.Template

type tplData struct {
	// Information about company and project
	Version             string
	Organization        string // Name of Organization
	OrganazationInfo    string // Name of Organization
	OrganizationLicense string //Org lic/copyright
	ProjectInfo         string // Project/service description
	MaintainerName      string
	MaintainerEmail     string
	MaintainerSlack     string
	MaintainerWeb       string

	// Integrations
	Badges     string // badges to include docs
	SonarKey   string
	SonarLogin string

	// Service and tpl-names
	Name         string //service name
	NameExported string //camal case with first letter cap
	TplName      string //template name
	DefFile      string //definition file used

	//PR lic/copyright should be a function
	PavedroadInfo string //PR lic/copyright

	//Swagger headers probably turn these into functions
	AllRoutesSwaggerDoc     string
	GetAllSwaggerDoc        string // swagger for list method
	GetSwaggerDoc           string // swagger for get method
	PutSwaggerDoc           string // swagger for put method
	PostSwaggerDoc          string // swagger for post method
	DeleteSwaggerDoc        string // swagger for delete method
	SwaggerGeneratedStructs string // swagger doc and go structs
	DumpStructs             string // Generic dumb of given object type

	//JSON data
	PostJSON string // Smaple data for a post
}

//  tplDataMapper
//    Map data from definitions file to tplData structure
//    return error if required mappings are missing
//
func tplDataMapper(defs tplDef, output *tplData) error {
	//fmt.Println(defs)
	output.Name = defs.Info.Name
	output.NameExported = strcase.ToCamel(defs.Info.Name)
	output.TplName = defs.Info.ID
	output.DefFile = tplDefFile
	output.Version = defs.Info.Version
	output.OrganizationLicense = defs.Project.License
	output.Organization = defs.Info.Organization
	output.ProjectInfo = defs.Project.Description
	output.MaintainerName = defs.Project.Maintainer.Name
	output.MaintainerEmail = defs.Project.Maintainer.Email
	output.MaintainerWeb = defs.Project.Maintainer.Web
	output.MaintainerSlack = defs.Project.Maintainer.Slack
	output.PavedroadInfo = prCopyright

	// CI integrations
	output.Badges = defs.BadgesToString()

	//Sonarcloud
	si := defs.findIntegration("sonarcloud")

	output.SonarKey = si.SonarCloudConfig.Key
	output.SonarLogin = si.SonarCloudConfig.Login
	return nil
}

//  tplJSONData
//    Use the schema definition found in tplDefs to create
//    Sample JSON data files
func tplJSONData(defs tplDef, output *tplData) error {
	var jsonString string
	order := defs.devineOrder()
	tplAddJSON(order, defs, &jsonString)

	//Make it pretty
	var pj bytes.Buffer
	err := json.Indent(&pj, []byte(jsonString), "", "\t")
	if err != nil {
		log.Fatal("Failed to generaet json data with ", jsonString)
	}
	output.PostJSON = string(pj.String())
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
		//TODO: validate column attributes
		// required attribute
		// no reserved go words

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
	//fmt.Println(output.SwaggerGeneratedStructs)
	return nil
}

// tplAddStruct
//
// Performs two tasks
//    - 1 Generates the structure as a string that is inserted
//        into the code template.  This is the tableString
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

			// Same as structField except type with be the subtable
			tableString += fmt.Sprintf(structSubstruct,
				strcase.ToCamel(child.Name),
				strings.ToLower(child.Name),
				"json",
				strings.ToLower(child.Name))
		}
	}

	// Add this tables attributes
	for _, col := range table.Columns {
		//TODO: validate column attributes
		// required attribute
		// no reserved go words

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
	RelativePath string // Path realative to the current directory
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
		//fmt.Println(body)

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
		nl += fmt.Sprintf("---\n") //replies contains multip documents
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

//tplPull pulls templates from a remote repository
//  pullOptions
//    all: default
//    microservices:
//    serverless:
//    crd:
//  org: GitHub orginization
//  repo: GitHub repository
//  path: path to start in repository
func tplPull(pullOptions, org, repo, path, outdir string) error {
	//TODO: make this an authenticated request/client
	client := github.NewClient(nil)

	opts := github.RepositoryContentGetOptions{}

	// Either file or directory content will be nil
	// file, director, resp, err
	fileContent, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, repo, path, &opts)

	if err != nil {
		//TODO: change to proper logging method
		fmt.Println(err)
		return err
	}
	//fmt.Println(rsp.StatusCode)
	//		if err != nil {
	//			log.Println(err)
	//			os.Exit(1)
	//		}

	if fileContent != nil {
		dstr, _ := base64.StdEncoding.DecodeString(*fileContent.Content)
		fp := outdir + "/" + *fileContent.Path
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			os.Create(fp)
		}
		err = ioutil.WriteFile(fp, dstr, 0644)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		//fmt.Println(directoryContent)
		for _, item := range directoryContent {
			// If it is a directory, create if necessary
			// Then walk it
			if *item.Type == "dir" {
				dn := outdir + "/" + *item.Path
				if _, err := os.Stat(dn); os.IsNotExist(err) {
					os.MkdirAll(dn, os.ModePerm)
					fmt.Println("Template directory created: ", dn)
				}
				_ = tplPull(pullOptions, org, repo, *item.Path, outdir)
			}

			// For files, request their content
			if *item.Type == "file" {
				_ = tplPull(pullOptions, org, repo, *item.Path, outdir)
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
//     - Go structs for user defined objects
//     - SQL Scripts for developer use
//     - SQL functions for testing
//
//   Next, iterate over the templates to generate
//
//   Last, run goimport to check for missing or unused import
//   statements and clean up any code formating issues
//
func tplCreate(rn string) string {
	var filelist []tplLocation
	//var filenames []string

	//Get back a list of templates for requrested template name
	tplRsp, err := tplRead(tplFile)
	if err != nil {
		fmt.Println(err)
	}

	// Get a list of the files names
	for _, rec := range tplRsp {
		nm := filepath.Base(rec)
		di := rec[len(tplDirSelected)+1 : len(rec)-len(nm)]
		li := tplLocation{Name: nm, RelativePath: di}

		// If definitions file, save it, otherwise add to filenames
		// to process as templates
		if nm == tplDefinition && tplDefFile == "" {
			fmt.Println("rec", rec)
			tplDefFile = rec
		} else {
			filelist = append(filelist, li)
		}
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

	// Generate internal structures
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

	// Build the template cache
	//templates, err = template.New("").ParseFiles(tplRsp...)
	templates, err = template.New("").ParseFiles(tplRsp...)
	if err != nil {
		fmt.Println("Template parsing failed: ", err)
		os.Exit(-1)
	}

	//templates = template.Must(template.ParseFiles(tplRsp...))
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

		// Ensure the path to the file exissts
		if v.RelativePath != "" {
			if _, err := os.Stat(v.RelativePath); os.IsNotExist(err) {
				err := os.MkdirAll(v.RelativePath, 0766)
				if err != nil {
					fmt.Println("Failed to make directory: ", v.RelativePath)
				}
			}
		}

		var defaultMode os.FileMode = 0666
		// Make sure shell scripts are created execuatble
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
		err = templates.ExecuteTemplate(bw, v.Name, tplInputData)
		if err != nil {
			fmt.Printf("Template execution failed for: %v with error %v", v, err)
			os.Exit(-1)
		}
		bw.Flush()
		file.Close()
	}

	// Execute goimport for code formating

	//fmt.Println(tplRsp)
	return ""
}

// tplReadDefinitions
//   Read the definition file and then validate it
//
func tplReadDefinitions(definitionsStruct *tplDef) error {

	fmt.Println("Reading defintions from: ", tplDefFile)
	df, err := os.Open(tplDefFile)
	if err != nil {
		fmt.Println("failed to open: %v err %v", df, err)
	}
	defer df.Close()
	byteValue, e := ioutil.ReadAll(df)
	if e != nil {
		fmt.Println("read failed for " + tplDefFile)
		os.Exit(-1)
	}

	//defData := make(map[interface{}]interface{})
	//err = yaml.Unmarshal(yamlMap, defs)
	//yaml.Unmarshal([]byte(byteValue), &defData)
	err = yaml.Unmarshal([]byte(byteValue), definitionsStruct)
	if err != nil {
		return err
	}

	//definitionsStruct.Validate()
	//will return a list of validation errors
	//exit after printing

	errs := definitionsStruct.Validate()

	//if lens(errs) > 0 {
	//	for _, v := range errs {
	//		v.Error()
	//	}
	//	os.Exit(-1)
	//}
	if errs != nil {
		for errs != nil {
			fmt.Println(errs.Error())
			errs = errs.nextError
		}
		os.Exit(-1)
	}
	return nil
}

// tplDescribe
func tplDescribe(tplListOption string, rn string) tplDescribeResponse {
	var response tplDescribeResponse

	// Get the list of templates
	rsp := tplGet("all", rn)

	// Load the defitions.yaml file in the template directory
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
				Type         string // Type of template, i.e. serverless
				Name         string // Name of template == directory name
			  Content      string // YAML configuration data
		*/
		nItem := tplDescribeItem{item.Type, item.Name, string(byteValue)}
		response.Templates = append(response.Templates, nItem)
		/*
			    //fmt.Println(string(byteValue))

			    //var body interface{}
			    body := make(map[interface{}]interface{})
					yaml.Unmarshal([]byte(byteValue), &body)

					//fmt.Println(result)
					fmt.Println(body["id"])
					fmt.Println(body["api-version"])
		*/
	}

	return response
}

// tplExplain
//   is stored in docs/explain.txt
//   docs is defined by eTLD
//   explain is tplResourceName
func tplExplain(tplListOption string, rn string) tplExplainResponse {
	var response tplExplainResponse

	// Load explaination
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

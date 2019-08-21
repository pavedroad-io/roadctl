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
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	_ "reflect"
	"strings"
  "github.com/iancoleman/strcase"
	"text/template"
)

// Default template repository
var defaultOrg string = "pavedroad-io"
var defaultRepo string = "templates"
var defaultPath string = ""
var defaultTemplateDir string = "templates"
var repoType string = "GitHub"
var tplFile string = "" // Name of the template file to use
var tplDir string = "." // Directory for generated code output
var tplDefFile string = "" // Directory for generated code output

const tplResourceName = "templates"
const tplDefinition = "definition.yaml"
const prCopyright = `
//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
//`

var templates *template.Template

type tplData struct {
  // Information about company and project
  Version string
  Organization string // Name of Organization
  OrganazationInfo string // Name of Organization
  OrganizationLicense string //Org lic/copyright
  ProjectInfo string // Project/service description
  Maintain string 
  MaintainerEmail string
  Maintainer string

  // Service and tpl-names
  Name string //service name
  NameExported string //camal case with first letter cap
  TplName string //template name

  //PR lic/copyright should be a function
  PavedroadInfo string //PR lic/copyright

  //Swagger headers probably turn these into functions
  AllRoutesSwaggerDoc string
  GetAllSwaggerDoc string // swagger for list method
  GetSwaggerDoc string // swagger for get method
  PutSwaggerDoc string // swagger for put method
  PostSwaggerDoc string // swagger for post method
  DeleteSwaggerDoc string // swagger for delete method
  SwaggerGeneratedStructs string //swagger doc and go structs
  DumpStructs string //Generic dumb of given object type
}

//  tplDataMapper
//    Map data from definitions file to tplData structure
//    return error if required mappings are missing
//
func tplDataMapper(defs tplDef, output *tplData) error {
  output.Name = defs.Name
  output.NameExported = strcase.ToCamel(defs.Name)
  output.TplName = defs.ID
  output.Version = defs.Version
  output.OrganizationLicense = defs.Project.License
  output.Organization = defs.Organization
  output.ProjectInfo = defs.Project.Description
  output.Maintain = "Contant: "
  output.MaintainerEmail = defs.Project.MaintainerEmail
  output.Maintainer = defs.Project.Maintainer
  output.PavedroadInfo = prCopyright

  return nil
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


  return nil
}

// tplListItem provides information about a template location
// and status
type tplListItem struct {
	Type         string // Type of template, i.e. serverless
	Availability string // Availability ga, ....
	Name         string // Name of template == directory name
	Path         string // Path to the template
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
		//line := strings.Repeat("-", 80)
		//nl += fmt.Println(string(line))
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
		fmt.Println(err)
		return err
	} else {
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
	var filenames []string

	//Get back a list of templates for requrested template name
	tplRsp, err := tplRead(tplFile)
	if err != nil {
		fmt.Println(err)
	}

	// Get a list of the files names
	for _, rec := range tplRsp {
		nm := filepath.Base(rec)
		// If definitions file, save it, otherwise add to filenames
		// to process as templates
		if nm == tplDefinition && tplDefFile == ""{
      fmt.Println("rec", rec)
			tplDefFile = rec
		} else {
			filenames = append(filenames, nm)
		}
	}

  // Read the definition file
  defs := tplDef{}
  err = tplReadDefinitions(&defs)
  if err != nil {
    fmt.Println(err)
    return(err.Error())
  }
  tplInputData := tplData{}
  err = tplDataMapper(defs, &tplInputData)

  // Generate internal structures
  err = tplGenerateStructurs(defs, &tplInputData)
  if err != nil {
    fmt.Println("Generating structures failed: ", err)
    os.Exit(-1)
  }

  // Generate SQL scripts

  // Generate SQL code functions

	// Build the template cache
	//templates, err = template.New("").ParseFiles(tplRsp...)
	templates, err = template.ParseFiles(tplRsp...)
  if err != nil {
    fmt.Println("Template parsing failed: ", err)
    os.Exit(-1)
  }

  //templates = template.Must(template.ParseFiles(tplRsp...))
  //TODO: turn into a function
  var fn string
  for _, v := range filenames {
    // Replace generic string "template" with the name of the service
    fn = strings.Replace(v, "template", tplInputData.Name, 1)

    file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        log.Fatal(err)
    }

    bw := bufio.NewWriter(file)

    fmt.Printf("executing %v using template %v\n", fn, v)
    err = templates.ExecuteTemplate(bw, v, tplInputData)
    if err != nil {
      fmt.Printf("Template execution failed for: %v with error %v", v, err)
      os.Exit(-1)
    }
    bw.Flush()
    file.Close()
  }
  
  // Execut goimport

	//fmt.Println(tplRsp)
	return ""
}

// tplReadDefinitions
//   Read the definition file
//
func tplReadDefinitions(definitionsStruct *tplDef) (error) {

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
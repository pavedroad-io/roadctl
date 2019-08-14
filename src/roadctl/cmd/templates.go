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
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	_ "reflect"
)

// Default template repository
var defaultOrg string = "pavedroad-io"
var defaultRepo string = "templates"
var defaultPath string = ""
var defaultTemplateDir string = "templates"
var repoType string = "GitHub"

//
type tplListItem struct {
	Type         string // Type of template, i.e. serverless
	Availability string // Availability ga, ....
	Name         string // Name of template == directory name
	Path         string // Path to the template
}

type tplDescribeItem struct {
	Type    string // Type of template, i.e. serverless
	Name    string // Name of template == directory name
	Content string // YAML configuration data
}

type tplDescribeResponse struct {
	Templates []tplDescribeItem
}

type tplListResponse struct {
	Templates []tplListItem
}

func (t tplDescribeResponse) RespondWithYAML() string {
	return t.RespondWithText() // One in the same for this type
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

	jb, err := yaml.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jb))

	return string(jb)
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

//

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

// List available templates
//  tplListOption: tBD
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

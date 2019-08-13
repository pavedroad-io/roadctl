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
  "gopkg.in/yaml.v2"
	_ "errors"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"os"
	_"reflect"
)

// Default template repository
var defaultOrg string = "pavedroad-io"
var defaultRepo string = "templates"
var defaultPath string = ""
var defaultTemplateDir string = "templates"
var repoType string = "GitHub"

// 
type tplListItem struct {
  Type  string
  Availability  string
  Name  string
}

type tplListResponse struct {
  Templates []tplListItem
}

func (t tplListResponse) RespondWithText() string {

  nl := fmt.Sprintf("%-15v %-20v %-20v\n","Template Type", "Name", "Release Status")
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

// List available templates
func tplGet(tplListOption string) tplListResponse {
  tplTLD := []string{"crd","microservices","serverless"}
  tplSLD := []string{"ga","experimental","incubation"}
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
         nrec := tplListItem{tld,sld,fn.Name()}
         response.Templates = append(response.Templates, nrec)
      }
    }
  }
  return response
}

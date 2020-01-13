// Package cmd module from cobra
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

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new resource",
	Long: `create a new resource taking input from stdin or a file
For example:

roadctl create templates -t template-name -d option-directory`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreate(cmd, args)
	},
}

// runCreate validates and then executes a creation
// of resources.  For now, it only supports creating
// one resource at at time
func runCreate(cmd *cobra.Command, args []string) string {
	msg := "Failed creating resource"
	r := args[0]

	if err := isValidResourceType(r); err == nil {
		if tplFile == "" {
			fmt.Println("Usage: roadctl create templates -t templateName")
			fmt.Println("       --template or -t option is required")
			return msg
		}
		return createResource(r)
	}
	return msg
}

func createResource(rn string) string {

	switch rn {
	case "environments":
		fmt.Println("no environments found")
		return ""
	case "builders":
		fmt.Println("no builders found")
		return ""
	case "taggers":
		fmt.Println("no taggers found")
		return ""
	case "tests":
		fmt.Println("no tests found")
		return ""
	case "templates":
		return tplCreate(rn)
	case "integrations":
		fmt.Println("no integrations found")
		return ""
	case "artifacts":
		fmt.Println("no artifacts found")
		return ""
	case "providers":
		fmt.Println("no providers found")
		return ""
	case "deployments":
		fmt.Println("no deployments found")
		return ""
	}

	return ""
}

func init() {
	rootCmd.AddCommand(createCmd)

	//Set up expected command line flags

	//tplfile defined in templates.go
	//Required!
	createCmd.Flags().StringVar(&tplFile, "template", "",
		"Template file name to use")

	//tplDir  defined in templates.go
	//this flag is not in documentaion 01/09/2020
	//Not Required!
	createCmd.Flags().StringVar(&tplDir, "directory", "",
		"Directory to generate output to")

	//tplDefFile defined in templates.go
	//Required!
	//Expexted YAML originaly generated from $roadctl desrcibe templates
	createCmd.Flags().StringVar(&tplDefFile, "definition", "",
		"Service definition file to use")
}

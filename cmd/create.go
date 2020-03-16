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

//Note: above and below blank lines required for golint.
//Related to required documentation format for packages.

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new resource",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requirs a resource type with an optional resource name")
		}
		return nil
	},
	Long: `create a new resource taking input from stdin or a file
For example:

roadctl create templates template-name -f definition.yaml`,
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

	if len(args) != 2 {
		fmt.Println("Usage: roadctl create templates templateName -f definiton.yaml")
		fmt.Printf("       templateName missing\n")
		return msg
	}

	tplFile = args[1]

	if tplDefFile == "" {
		fmt.Println("Usage: roadctl create templates templateName -f definiton.yaml")
		fmt.Printf("       -f definitions file  missing\n")
		return msg
	}

	if err := isValidResourceType(r); err == nil {
		if tplFile == "" {
			fmt.Println("Usage: roadctl create templates templateName -f templateName")
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
	// createCmd.Flags().StringVarP(&tplFile, "template", "t",
	//	"datamgr", "Template file name to use")

	// tplDefFile defined in templates.go
	// Expected YAML originally generated from $roadctl describe templates >> myservice.yaml
	createCmd.Flags().StringVarP(&tplDefFile, "file", "f",
		"", "Service definition file to use")
}

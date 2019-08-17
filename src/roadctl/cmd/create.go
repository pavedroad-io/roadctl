/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
			fmt.Println("Usage: roadctl create emplates -t templateName")
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
	case "tool-network":
		fmt.Println("no tool-network found")
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

	createCmd.Flags().StringVar(&tplFile, "template", "t",
		"Template file name to use")
	createCmd.Flags().StringVar(&tplDir, "directory", "dn",
		"Directory to generate output to")
}

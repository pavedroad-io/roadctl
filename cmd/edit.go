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

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit the configuration for the specified resource",
	Long: `Invoke the configured $EDITOR and load the current configuration
  for the named resource.
  For example:
    roadctl edit template foo`,
	Run: func(cmd *cobra.Command, args []string) {
		runEdit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

// runEdit
func runEdit(cmd *cobra.Command, args []string) {
	resources := strings.Split(args[0], ",")

	for _, r := range resources {
		if err := isValidResourceType(r); err == nil {
			if len(args) > 1 {
				editResource(r, args[1])
			} else {
				fmt.Println("resource name required")
			}
		} else {
			fmt.Println("Edit failed: ", err)
		}
	}
}

func editResource(r, n string) {
	switch r {
	case "environments":
		fmt.Println("no environments found")
		return
	case "builders":
		fmt.Println("no builders found")
		return
	case "taggers":
		fmt.Println("no taggers found")
		return
	case "tests":
		fmt.Println("no tests found")
		return
	case "templates":
		tplEdit(n)
		return
	case "integrations":
		fmt.Println("no integrations found")
		return
	case "artifacts":
		fmt.Println("no artifacts found")
		return
	case "providers":
		fmt.Println("no providers found")
		return
	case "deployments":
		fmt.Println("no deployments found")
		return
	}

	return
}

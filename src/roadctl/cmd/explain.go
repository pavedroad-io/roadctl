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

//Note: above and below blank lines required for golint.
//Related to required documentation format for packages.

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var eTLD = "docs"

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "return documentation about a resource",
	Long: `Return documentation about the structure of a resource
  For example:
    roadctl explain template`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a resource type with an optional resource name")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		runExplain(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)

}

func explainResource(r, n string) Response {
	var rsp Response
	switch r {
	case "environments":
		fmt.Println("no environments found")
		return nil
	case "builders":
		fmt.Println("no builders found")
		return nil
	case "taggers":
		fmt.Println("no taggers found")
		return nil
	case "tests":
		fmt.Println("no tests found")
		return nil
	case "templates":
		rsp = tplExplain("all", n)
		return rsp
	case "integrations":
		fmt.Println("no integrations found")
		return nil
	case "artifacts":
		fmt.Println("no artifacts found")
		return nil
	case "providers":
		fmt.Println("no providers found")
		return nil
	case "deployments":
		fmt.Println("no deployments found")
		return nil
	}

	return nil
}

// runExplain validates and then executes a describe command
//
func runExplain(cmd *cobra.Command, args []string) {
	replies := []Response{}
	var reply Response

	resources := strings.Split(args[0], ",")

	for _, r := range resources {
		if err := isValidResourceType(r); err == nil {
			if len(args) > 1 {
				reply = explainResource(r, args[1])
			} else {
				reply = explainResource(r, "")
			}
		} else {
			fmt.Println(err)
		}
		replies = append(replies, reply)
	}

	for _, r := range replies {
		//Hack until all assets return a Response type
		// fmtFlag is just "f" if not specified
		if r != nil {
			switch strings.ToLower(fmtFlag) {
			case "text":
				r.RespondWithText()
				break
				// Only support text replies for now
				/*
					case "yaml":
						r.RespondWithYAML()
						break
					case "json":
						r.RespondWithJSON()
						break
				*/
			default:
				r.RespondWithText()
			}
		}
	}
}

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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe provides detailed information about a resource",
	Long: `Provided detailed information about one or more resources.
For example:

  roadctl describe blueprint datamgr`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a resource type with an optional resource name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		runDescribe(cmd, args)
	},
}

// runDescribe validates and then executes a describe command
//
func runDescribe(cmd *cobra.Command, args []string) {
	replies := []Response{}
	var reply Response

	resources := strings.Split(args[0], ",")

	for _, r := range resources {
		if err := isValidResourceType(r); err == nil {
			if len(args) > 1 {
				reply = describeResource(r, args[1])
			} else {
				reply = describeResource(r, "")
			}
		} else {
			fmt.Println("Describe failed: ", err)
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
			case "yaml":
				r.RespondWithYAML()
				break
			case "json":
				r.RespondWithJSON()
				break
			default:
				r.RespondWithText()
			}
		}
	}
}

func describeResource(r, n string) Response {
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
	case "blueprints":
		rsp = bpDescribe("all", n)
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

func init() {
	rootCmd.AddCommand(describeCmd)
	/*
	  describeCmd.SetUsageBlueprint("usage: foobar")
	*/
	describeCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return errors.New("Usage: roadctl VERB NOUN NAME")
	})

}

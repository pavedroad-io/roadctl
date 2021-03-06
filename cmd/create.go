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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resource",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a resource type with an optional resource name")
		}
		return nil
	},
	Long: `create a new resource taking input from stdin or a file
For example:

roadctl create blueprints blueprints-name -f definition.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreate(cmd, args)
	},
}

// runCreate validates and then executes a creation
// of resources.  For now, it only supports creating
// one resource at at time
func runCreate(cmd *cobra.Command, args []string) bpCreateResponse {
	var reply bpCreateResponse
	r := args[0]

	if len(args) != 2 {
		var errorResponse = bpGenericResponseItem{
			Name: "create usage",
			Content: `
			Usage: creating blueprints requires two arguments
			       a blueprint name and a definitions files specified with -f
				  
				   syntax: roadctl create blueprints NAME -f DEFINITIONS.YAML
			       `,
		}
		reply.Blueprints = append(reply.Blueprints, errorResponse)
		return reply
	}

	bpFile = args[1]

	if err := isValidResourceType(r); err == nil {
		reply = bpCreate(r)
	}

	if len(reply.Blueprints) > 0 {
		switch strings.ToLower(fmtFlag) {
		case "text":
			reply.RespondWithText()
			break
			// Only support text replies for now
		case "yaml":
			reply.RespondWithYAML()
			break
		case "json":
			reply.RespondWithJSON()
			break
		default:
			reply.RespondWithText()
		}
	}

	return reply
}

func createResource(rn string) bpCreateResponse {

	var notImplemented bpCreateResponse
	var notImplementedItem bpGenericResponseItem = bpGenericResponseItem{
		Name:    "Resource not implemented",
		Content: "",
	}
	notImplemented.Blueprints = append(notImplemented.Blueprints, notImplementedItem)

	switch rn {
	case "environments":
		fmt.Println("no environments found")
		return notImplemented
	case "builders":
		fmt.Println("no builders found")
		return notImplemented
	case "taggers":
		fmt.Println("no taggers found")
		return notImplemented
	case "tests":
		fmt.Println("no tests found")
		return notImplemented
	case "blueprints":
		createResource(rn)
	case "integrations":
		fmt.Println("no integrations found")
		return notImplemented
	case "artifacts":
		fmt.Println("no artifacts found")
		return notImplemented
	case "providers":
		fmt.Println("no providers found")
		return notImplemented
	case "deployments":
		fmt.Println("no deployments found")
		return notImplemented
	}

	return notImplemented
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&bpDefFile, "file", "f",
		"", "Service definition file to use")
}

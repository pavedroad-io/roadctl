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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	prHome = ".pavedroad.d"
)

var repository string
var branch string

// docCmd represents the doc command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize roadctl development environment",
	Long: `Initialize roadctl development environment
  Options:
    --api  Use GitHub API instead of clone (not recommended)
	       You must also specify a GitHub authentication method`,
	Run: func(cmd *cobra.Command, args []string) {
		runInit(cmd, args)
	},
}

// runGet validates and then executes a get command
func runInit(cmd *cobra.Command, args []string) {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("init: ", err)
		os.Exit(1)
	}

	phome := home + "/" + prHome
	createDirectory(phome)

	// See if we need to update or initialize the template repository
	apiTrue := cmd.Flag("api")
	initTemplates(apiTrue)
}

// initTemplates: Download templates from GitHub
// If the template dir is location, you can prefix
// with "_" or "." to have go skip them when compiling
//
func initTemplates(api *pflag.Flag) {
	fmt.Println("Initializing template repository")
	if api.Value.String() == "true" {
		// Create template dir if necessary
		if _, err := os.Stat(defaultTemplateDir); os.IsNotExist(err) {
			fmt.Println("defaultTemplateDir")
			os.MkdirAll(defaultTemplateDir, os.ModePerm)
		}
		client := getClient()
		tplPull("all", defaultOrg, defaultRepo, defaultPath, defaultTemplateDir, client)
	} else {
		tplClone(branch)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("api", "a", false, "Initialize template repository using GitHub API")
	initCmd.Flags().StringVarP(&repository, "repo", "r", "https://github.pavedroad-io/templates",
		"Change default repository for templates")
	initCmd.Flags().StringVarP(&branch, "branch", "b", "release", "Branch to clone (default release)")
}

// Package cmd module from cobra
package cmd

/*
<<<<<<< HEAD:src/roadctl/cmd/config.go
Copyright © 2019 PavedRoad <info@pavedroad.io>
=======
Copyright © 2019 PavedRoad <info@pavedroad>
>>>>>>> master:cmd/config.go

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
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage roadctl global configuration options",
	Long: `Allows you to manage global configurations

  For example:

  roadctl config environment test`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

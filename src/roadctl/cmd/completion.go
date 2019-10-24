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

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion scripts",
	Long: `To load completion run

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile

or rename to ~.bash_completion
`,
	Run: func(cmd *cobra.Command, args []string) {
		zshTrue := cmd.Flag("zsh")
		if zshTrue.Value.String() == "true" {
			fmt.Println("Writing: " + "roadctlZshCompletion.sh")
			rootCmd.GenZshCompletionFile("roadctlZshCompletion.sh")
		} else {
			fmt.Println("Writing: " + "roadctlBashCompletion.sh")
			rootCmd.GenBashCompletionFile("roadctlBashCompletion.sh")
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("zsh", "z", false, "Generate zsh completion")
}

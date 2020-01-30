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
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	ValidArgs: []string{"bash", "zsh"},
	Use:       "completion bash|zsh",
	Short:     "Generate completion scripts on stdout",
	Long: `To create a roadctl completion file:
  roadctl completion bash > roadctl.bash
Or:
  roadctl completion zsh > _roadctl
Then move the file to the appropriate completion directory

Or to load the completion code into the current shell:
  source <(roadctl completion bash)
Or:
  source <(roadctl completion zsh)
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Missing completion type of bash or zsh\n")
		} else if args[0] == "bash" {
			rootCmd.GenBashCompletion(os.Stdout)
		} else if args[0] == "zsh" {
			rootCmd.GenZshCompletion(os.Stdout)
		} else {
			fmt.Fprintf(os.Stderr, "Unsupported completion type %s\n", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

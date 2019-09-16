/*
Copyright © 2019 PavedRoad

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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var fmtFlag string = "text"
var debugFlag string = "info"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "roadctl",
	Short: "CLI for PavedRoad Development Kit (DevKit)",
	Long: `roadctl allows you to work with the PavedRoad CNCF low-code environment and the associated CI/CD pipeline

  Usage: roadctl [command] [TYPE] [NAME] [flags]

  TYPE specifies a resource type
  NAME is the name of a resource
  flags specify options`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initConstants)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.roadctl.yaml)")
	rootCmd.PersistentFlags().StringVar(&fmtFlag, "format", "f", "Output format: text(default)|json|yaml")
	rootCmd.PersistentFlags().StringVar(&debugFlag, "debug", "d", "Debug level: info(default)|warm|error|critical")
}

// initConstants populates global slices of types
func initConstants() {
	// Types or resouces command can act on
	initResourcetypes()

	return
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".roadctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".roadctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

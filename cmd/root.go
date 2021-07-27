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
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var fmtFlag string = "text"
var debugFlag string = "info"
var userName string
var userPassword string
var userAccessToken string
var blueprintDirectoryLocation string
var rclog *log.Logger
var rcLogFile *os.File
var rcLogFileName string = "debug.log"

// GitTag
var GitTag string

// Version
var Version string

// Build
var Build string

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
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initConstants)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.roadctl.yaml)")
	rootCmd.PersistentFlags().StringVar(&fmtFlag, "format", "text", "Output format: text(default)|json|yaml")
	rootCmd.PersistentFlags().StringVar(&debugFlag, "debug", "info", "Debug level: info(default)|warm|error|critical")

	// For API calls that require authentication
	rootCmd.PersistentFlags().StringVar(&userName, "user", "", "HTTP basic auth user name")
	rootCmd.PersistentFlags().StringVar(&userPassword, "password", "", "HTTP basic auth password")
	rootCmd.PersistentFlags().StringVar(&userAccessToken, "token", "", "OAUTH access token")

	// Blueprint directory
	//   Don't set the default so we know if it was set by command line vs env
	rootCmd.PersistentFlags().StringVar(&blueprintDirectoryLocation, "blueprints", "", "Set the location of the directory holding roadctl blueprints")
}

// initConstants populates global slices of types
func initConstants() {
	// Types or resources command can act on
	initResourcetypes()

	// TODO: Move this into viper configuration
	initAuthentication()

	return
}

// initAuthentication: Set global authentication variables
//   OAUTH token takes precedents over basic authentication
//   so quit if we find that as a command line option or
//   environment variable
//
func initAuthentication() {
	if userAccessToken != "" {
		return
	}
	//No dashes in environment name, use underscore
	//Prefix for environment name with service (GH GitHub)
	envVar := os.Getenv("GH_ACCESS_TOKEN")

	if envVar != "" {
		userAccessToken = envVar
		return
	}

	// Command line rules over envVar and we have both
	if userName != "" && userPassword != "" {
		return
	}

	//No dashes in environment name, use underscore
	//Prefix for environment name with service (GH GitHub)
	envName := os.Getenv("GH_USER_NAME")
	envPass := os.Getenv("GH_USER_PASSWORD")

	if envPass != "" {
		userPassword = envPass
	}

	// Allow for case where userName is from command
	// line and password is from envVar
	if userName != "" {
		return
	}

	if envName != "" {
		userPassword = envPass
	}
	return
}

func initLogging() {

	var rcLogFile, err1 = os.Create(rcLogFileName)

	if err1 != nil {
		panic(err1)
	}

	rclog = log.New(rcLogFile, "roadctl: ", log.LstdFlags|log.Lshortfile)

	rclog.SetOutput(rcLogFile)
	fmt.Println("Logging to " + rcLogFileName)

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
			fmt.Println("initConfig: ", err)
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

// Package cmd this files handels functions related to integratiosn
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sonarcloud "github.com/pavedroad-io/integrations/sonarcloud/cmd"
)

// validateIntegrations for all valid integrations
// verify the configurations are good or complete them
func validateIntegrations(config *tplData) error {

	err := checkSonarCloud(config)

	return err
}

// checkSonarCloud validates project key and generates
// a token if not present
func checkSonarCloud(config *tplData) error {
	scClient := sonarcloud.SonarCloudClient{}

	// TODO: Make this part of the New() method
	var token string
	envVar := os.Getenv("SONARCLOUD_TOKEN")
	if envVar != "" {
		token = envVar
	} else {
		log.Println("Need SONARCLOUD_TOKEN set to run tests")
		os.Exit(-1)
	}

	err := scClient.New(token, 10)

	if err != nil {
		log.Println("failed to create New SonarCloudClient")
		return err
	}

	private := false
	err = ensureSonarCloudKeyExists(scClient,
		config.Organization, config.SonarKey, config.ProjectInfo,
		private)

	if config.SonarLogin == "" {
		var tokenName string
		if config.Organization != "" {
			tokenName = config.Organization + "_" + config.SonarKey + " CI key"
		} else {
			tokenName = config.SonarKey + " CI key"
		}

		rsp, err := scClient.CreateToken(tokenName)
		if err != nil {
			log.Printf("Failed to create token, err %v", err)
			return err
		}

		token, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Printf("Failed to read response body, err %v", err)
			return err
		}

		var tk sonarcloud.NewTokenResponse
		err = json.Unmarshal(token, &tk)
		if err != nil {
			log.Printf("Failed to Unmarshal token, Err %v\n", err)
			return err
		}

		config.SonarLogin = tk.Token
	}

	return nil
}

// ensureSonarCloudKeyExists search for all posible
// combinations:
//   key in global namespace == key
//   key in private namespace == organization_key
// Create in global namespace is it doesn't exist
func ensureSonarCloudKeyExists(client sonarcloud.SonarCloudClient,
	org, key, name string,
	public bool) error {
	possibleNames := []string{key, org + "_" + key}

	for _, v := range possibleNames {
		fmt.Println(v)
		_, err := client.GetProject(org, key)

		//
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				// Continue because the records already exists
				return nil
			}
			// Else, we can't connect to the server
			log.Println("failed conneting to SonarCloud.io")
			return err

		}
	}

	// Make the project in the global namespace
	// A prefix will be added to ensure uniqueness
	visability := "public"

	if !public {
		visability = "private"
	}

	p := sonarcloud.NewProject{
		Organization: org,
		Name:         name,
		Project:      key,
		Visibility:   visability,
	}

	_, err := client.CreateProject(p)

	if err != nil {
		log.Println("Failed creating project:", key)
		return err
	}

	return nil
}

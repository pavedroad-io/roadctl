// Package cmd this files handles functions related to integrations
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sonarcloud "github.com/pavedroad-io/integrations/sonarcloud/cmd"
)

// validateIntegrations for all valid integrations
// verify the configurations are good or complete them
func validateIntegrations(config *bpData) error {

	err := checkSonarCloud(config)

	return err
}

// scConnect return a SonarCloudClient or error
func scConnnect() (client sonarcloud.SonarCloudClient, err error) {
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

	err = scClient.New(token, 10)

	if err != nil {
		log.Println("failed to create New SonarCloudClient")
		return scClient, err
	}
	return scClient, nil

}

// scBadges, give a list of valid SonarCloud badges return a
//   list of emeddable Markdown links
//   Always include the quality gate badge if SonarCloud is enabled
func scBadges(
	config *bpData,
	badgeList []string) (resultList []string, err error) {

	client, err := scConnnect()
	if err != nil {
		fmt.Println("Connection failed: ", err)
		return resultList, err
	}

	pk := config.SonarPrefix + config.Name

	if len(badgeList) > 0 {
		for _, v := range badgeList {
			var found int
			if found = client.NameToEnum(v); found == -1 {
				fmt.Println("Not found: ", v)
			}
			// metric, project key, branch
			rsp, err := client.GetMetric(found, pk, "")
			if err != nil {
				fmt.Println(err)
				continue
			}
			badge, _ := ioutil.ReadAll(rsp.Body)

			resultList = append(resultList, string(badge))
		}
	}

	return resultList, err

}

// checkSonarCloud validates project key and generates
// a token if not present
func checkSonarCloud(config *bpData) error {
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

	private := true
	err = ensureSonarCloudKeyExists(scClient,
		config.Organization, config.SonarKey, config.Name,
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

// ensureSonarCloudKeyExists search for all possible
// combinations:
//   key in global namespace == PavedRoad_{{key}}
//   key in private namespace == organization_key
// Create in global namespace is it doesn't exist
func ensureSonarCloudKeyExists(client sonarcloud.SonarCloudClient,
	org, key, name string,
	public bool) error {

	resp, err := client.GetProject(org, key)

	if err != nil {
		// Continue because we can't connect to the server
		log.Println("failed connecting to SonarCloud.io")
		return err
	}

	// A 200 doesn't mean we found the record
	project, err := ioutil.ReadAll(resp.Body)
	var prj sonarcloud.ProjectSearchResponse
	err = json.Unmarshal(project, &prj)

	if err != nil {
		fmt.Println("Failed to Unmarshal sonarcloud search response")
		fmt.Println(err)
		return err
	}

	// If it is creater than one the project has already been created
	if len(prj.Components) >= 1 {
		fmt.Printf("Project %s already exists, skip creating\n", key)
		return nil
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

	_, err = client.CreateProject(p)

	if err != nil {
		log.Printf("Failed creating project %s\n", name)
		log.Println(err)
		return err
	}

	return nil
}

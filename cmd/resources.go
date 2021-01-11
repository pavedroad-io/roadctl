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
)

var resourceTypes []string

func initResourcetypes() {
	resourceTypes = []string{"environments",
		"builders",
		"packagers",
		"taggers",
		"tests",
		"blueprints",
		"integrations",
		"artifacts",
		"providers",
		"deployments"}
}

func isValidResourceType(typeToCheck string) error {
	for _, val := range resourceTypes {
		if val == typeToCheck {
			return nil
		}
	}
	msg := fmt.Sprintf("resouce must be on of %v", resourceTypes)
	return errors.New(msg)
}

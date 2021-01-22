// Package cmd from cobra
package cmd

/*
Copyright © 2019,2020,2021 PavedRoad <info@pavedroad.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

response: defines mandatory response methods a resource MUST implement
  RespondWithText:
  RespondWithYAML:
  RespondWithJSON:
*/

import (
	// used for debug
	_ "errors"
	// used for debug
	_ "fmt"
)

// Response must include all three
type Response interface {
	RespondWithJSON() string
	RespondWithYAML() string
	RespondWithText() string
}

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
package main

import (
	"github.com/pavedroad-io/roadctl/cmd"
)

// GitTag contains current git tab for this repository
var GitTag string

// Version contains version specified in definitions file
var Version string

// Build holds latest git commit hash in short form
var Build string

func main() {
	cmd.GitTag = GitTag
	cmd.Version = Version
	cmd.Build = Build
	cmd.Execute()
}

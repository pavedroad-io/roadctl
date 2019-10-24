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
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// getClient: return a githubClient with:
//   TokenAuth, BasicAuth, or no authentication
//
func getClient() *github.Client {
	var client *github.Client

	if userAccessToken != "" {
		// User oAuth token if available
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: userAccessToken},
		)
		// Construct *http.Client to pass to github.NewClient
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else if userName != "" && userPassword != "" {
		// Use basic auth transport
		bat := github.BasicAuthTransport{
			Username: strings.TrimSpace(userName),
			Password: strings.TrimSpace(userPassword),
		}
		client = github.NewClient(bat.Client())
	} else {
		// Use an unauthenticated client
		client = github.NewClient(nil)
	}

	return client
}

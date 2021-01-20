// cmd Corbra CLI
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
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
)

// The release branch stores released blueprints
// The latest and stable tags are used to select which release
const (
	gitBlueprintBranch  plumbing.ReferenceName = "refs/heads/release"
	gitLatestTag        plumbing.ReferenceName = "refs/heads/latest"
	gitStableTag        plumbing.ReferenceName = "refs/heads/stable"
	blueprintRepository                        = "https://github.com/pavedroad-io/blueprints"
	githubAPI                                  = "GitHub API"
	gitclone                                   = "git clone"

	// File that holds meta-data about a cache
	bpCacheFileName string = ".pr_cache"
)

// bpDirectory manages blueprint directory locations
type bpDirectory struct {
	// Full path to the blueprint directory
	location string

	// Is it initialized
	initialized bool

	// How we determined the location
	// default, so.GetEnv, command line option
	locationFrom string
}

// bpCache manages information about blueprints
//  stored in a blueprint directory
type bpCache struct {
	// What directory is it in
	location *bpDirectory

	// Persist information to disk in a cache file
	CacheFile bpCacheFile `json:"cache_file"`
}

// errno constants for bpCacheError
const (
	tcBadBlueprintDirectory = iota
	tcNoCacheFile
	tcBadCacheFile
	tcSuccess
)

// errmsg constants for bpCacheError
const (
	TcBadBlueprintDirectory = "Unable to create blueprint directory, Got (%v)\n"
	TcNoCacheFile           = "Cache file not found (%v)\n"
	TcBadCacheFile          = "Bad cache file (%v)\n"
)

// bpCacheError
type bpCacheError struct {
	errno  int
	errmsg string
}

func (tc *bpCacheError) Error() string {
	return tc.errmsg
}

// bpCacheFile Store information to disk for later access
type bpCacheFile struct {
	// Cache file syntax version
	Version string `json:"version"`

	// Time it was created on disk
	Created time.Time `json:"created"`

	// last updated
	Updated time.Time `json:"updated"`

	// Is it initialized
	Initialized bool `json:"initialized"`

	// git clone or github API
	InitializedFrom string `json:"initialized_from"`

	// git branch
	Branch string `json:"branch"`

	// git tags
	Tags []string `json:"tags"`
}

// Location returns the location of the blueprint directory
// Initialize if necessary
func (t *bpDirectory) Location() string {

	if !t.initialized {
		err := t.initialize()
		if err != nil {
			log.Fatal(err.Error())
			return ""
		}
	}
	return t.location
}

// initialize a private function for initializing
//   the blueprint directory location, not the
//   blueprints
func (t *bpDirectory) initialize() error {
	// Order of precedence
	//   - roadctl CLI
	//   - PR_BLUEPRINT_DIR
	//   - defaultBlueprintDir

	env := os.Getenv("PR_BLUEPRINT_DIR")

	if blueprintDirectoryLocation != "" {
		// Make sure they
		if strings.Contains(blueprintDirectoryLocation, "/"+defaultBlueprintDir) {
			fmt.Printf("Don't include /blueprint with --blueprint option\n")
			if string(blueprintDirectoryLocation[len(blueprintDirectoryLocation)-1]) == "/" {
				fmt.Printf("Use: --blueprint %s\n",
					strings.TrimSuffix(blueprintDirectoryLocation, "/"+defaultBlueprintDir+"/"))
			} else {
				fmt.Printf("Use: --blueprint %s\n",
					strings.TrimSuffix(blueprintDirectoryLocation, "/"+defaultBlueprintDir))
			}
			os.Exit(-1)

		}
		t.location = blueprintDirectoryLocation + "/" + defaultBlueprintDir
		t.locationFrom = "CLI"
		if defaultBlueprintDir != t.location {
			defaultBlueprintDir = t.location
		}
	} else if env != "" && blueprintDirectoryLocation == "" {
		t.location = env + "/" + defaultBlueprintDir
		t.locationFrom = "PR_BLUEPRINT_DIR"
		if defaultBlueprintDir != t.location {
			defaultBlueprintDir = t.location
		}
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("error setting home directory")
		}
		home = home + "/" + prHome + "/" + defaultBlueprintDir

		t.location = home
		t.locationFrom = "default"
	}

	if !t.initialized {
		if err := createDirectory(t.location); err != nil {
			log.Fatal(err.Error())
		}
		t.initialized = true
	}

	return nil
}

func (t *bpDirectory) getDefault() string {

	return ""
}

// New create a bpCache
// If it does not exists, initialize it using method specified
//   td: a bpDirectory type
//   method: GitHub API or git clone
func (tc *bpCache) CreateCache(method, branch string) error {
	log.Println("Cloning with method: ", method)
	switch method {
	case gitclone:
		tc.Clone(branch)
	case githubAPI:
		//tc.API()
	default:
		fmt.Println("error")
	}
	return nil
}

// NewBlueprintCache read the blueprint directory and it's meta-data
func NewBlueprintCache() (*bpCache, bpCacheError) {
	t := &bpDirectory{}
	tc := &bpCache{location: t}
	te := bpCacheError{errno: tcSuccess}

	if dir := t.Location(); dir == "" {
		te.errno = tcBadBlueprintDirectory
		te.errmsg = fmt.Sprintf(TcBadBlueprintDirectory, dir)
		return tc, te
	}
	err := tc.readCache()
	if err != nil {
		te.errno = tcNoCacheFile
		te.errmsg = fmt.Sprintf(TcNoCacheFile, err)
		return tc, te
	}

	return tc, te
}

// Clone Do a git clone if the repository doesn't exist
// in the desired directory
func (tc *bpCache) Clone(branch string) error {

	// See if it already has been cloned
	if e := tc.readCache(); e == nil {
		log.Println("Cache exists exiting")
		return e
	}

	/* TODO: Fix this
	   cb := gitBlueprintBranch
	   if branch != gitBlueprintBranch.String() {
	       cb = plumbing.ReferenceName(branch)
	   }
	*/

	_, err := git.PlainClone(defaultBlueprintDir, false, &git.CloneOptions{
		URL:           blueprintRepository,
		ReferenceName: "refs/heads/release",
		//ReferenceName: cb,
	})

	if err != nil {
		log.Fatal(err)
	}

	tc.CacheFile.Created = time.Now()
	tc.CacheFile.Updated = time.Now()
	tc.CacheFile.InitializedFrom = "clone"
	tc.CacheFile.Version = "1.0.0alpha"
	tc.CacheFile.Initialized = true
	tc.CacheFile.Branch = gitBlueprintBranch.String()

	if err = tc.writeCache(defaultBlueprintDir); err != nil {
		log.Fatal(err)
	}

	return nil
}

// writeCache save cache file in repository
//   dir is set if not using tc.location
func (tc *bpCache) writeCache(dir string) error {
	var cfn string

	fileData, err := yaml.Marshal(tc.CacheFile)
	if err != nil {
		return err
	}

	if dir != "" {
		cfn = dir + "/" + bpCacheFileName
	} else {
		cfn = tc.location.Location() + "/" + bpCacheFileName
	}

	err = ioutil.WriteFile(cfn, fileData, 0644)
	if err != nil {
		log.Fatal(err, cfn)
		return err
	}

	return nil
}

// readCache Opens the cache file and reads contents
func (tc *bpCache) readCache() error {
	fl := tc.location.Location() + "/" + bpCacheFileName
	if _, err := os.Stat(fl); os.IsNotExist(err) {
		return errors.New("Not found")
	}

	df, err := os.Open(fl)
	if err != nil {
		log.Println("failed to open:", fl, ", error:", err)
	}

	defer df.Close()
	byteValue, e := ioutil.ReadAll(df)
	if e != nil {
		log.Fatal("read failed for ", fl)
	}

	err = yaml.Unmarshal(byteValue, &tc.CacheFile)
	if err != nil {
		return err
	}

	return nil
}

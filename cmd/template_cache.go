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
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
)

// The release branch stores released templates
// The latest and stable tags are used to select which release
const (
	gitTemplateBranch  plumbing.ReferenceName = "refs/heads/release"
	gitLatestTag       plumbing.ReferenceName = "refs/heads/latest"
	gitStableTag       plumbing.ReferenceName = "refs/heads/stable"
	templateRepository                        = "https://github.com/pavedroad-  io/templates"
	githubAPI                                 = "GitHub API"
	gitclone                                  = "git clone"

	// File that holds meta-data about a cache
	tplCacheFileName string = ".pr_cache"
)

// tplDirectory manages template directory locations
type tplDirectory struct {
	// Full path to the template directory
	location string

	// Is it initialized
	initialized bool

	// How we determined the location
	// default, so.GetEnv, command line option
	locationFrom string
}

// tplCache manages information about templates
//  stored in a template directory
type tplCache struct {
	// What directory is it in
	location *tplDirectory

	// Persist information to disk in a cache file
	CacheFile tplCacheFile `json:"cache_file"`
}

// errno constants for tplCacheError
const (
	tcBadTemplateDirectory = iota
	tcNoCacheFile
	tcBadCacheFile
	tcSuccess
)

// errmsg constants for tplCacheError
const (
	TcBadTemplateDirectory = "Unable to create template directory, Got (%v)\n"
	TcNoCacheFile          = "Cache file not found (%v)\n"
	TcBadCacheFile         = "Bad cache file (%v)\n"
)

// tplCacheError
type tplCacheError struct {
	errno  int
	errmsg string
}

func (tc *tplCacheError) Error() string {
	return tc.errmsg
}

// tplCacheFile Store information to disk for later access
type tplCacheFile struct {
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

// Location returns the location of the template directory
// Initialize if necessary
func (t *tplDirectory) Location() string {

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
//   the template directory location, not the
//   templates
func (t *tplDirectory) initialize() error {
	// Order of precedence
	//   - roadctl CLI
	//   - PR_TEMPLATE_DIR
	//   - defaultTemplateDir

	env := os.Getenv("PR_TEMPLATE_DIR")
	if templateDirectoryLocation != "" {
		t.location = templateDirectoryLocation
		t.locationFrom = "CLI"
	} else if env != "" && templateDirectoryLocation == "" {
		t.location = env
		t.locationFrom = "PR_TEMPLATE_DIR"
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("error setting home directory")
		}
		home = home + "/" + prHome + "/" + defaultTemplateDir

		t.location = home
		// TODO: remove this hack once defaultTemplateDir is removed
		// For now, avoid duplicating the location string
		if home != defaultTemplateDir {
			defaultTemplateDir = home
		}

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

func (t *tplDirectory) getDefault() string {

	return ""
}

// New create a tplCache
// If it does not exists, initialize it using method specified
//   td: a tplDirectory type
//   method: GitHub API or git clone
func (tc *tplCache) CreateCache(method, branch string) error {
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

// NewTemplateCache read the template directory and it's meta-data
func NewTemplateCache() (*tplCache, tplCacheError) {
	t := &tplDirectory{}
	tc := &tplCache{location: t}
	te := tplCacheError{errno: tcSuccess}

	if dir := t.Location(); dir == "" {
		te.errno = tcBadTemplateDirectory
		te.errmsg = fmt.Sprintf(TcBadTemplateDirectory, dir)
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
func (tc *tplCache) Clone(branch string) error {

	// See if it already has been cloned
	if e := tc.readCache(); e == nil {
		log.Println("Cache exists exiting")
		return e
	}

	/* TODO: Fix this
	   cb := gitTemplateBranch
	   if branch != gitTemplateBranch.String() {
	       cb = plumbing.ReferenceName(branch)
	   }
	*/

	_, err := git.PlainClone(defaultTemplateDir, false, &git.CloneOptions{
		URL:           templateRepository,
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
	tc.CacheFile.Branch = gitTemplateBranch.String()

	if err = tc.writeCache(defaultTemplateDir); err != nil {
		log.Fatal(err)
	}

	return nil
}

// writeCache save cache file in repository
//   dir is set if not using tc.location
func (tc *tplCache) writeCache(dir string) error {
	var cfn string

	fileData, err := yaml.Marshal(tc.CacheFile)
	if err != nil {
		return err
	}

	if dir != "" {
		cfn = dir + "/" + tplCacheFileName
	} else {
		cfn = tc.location.Location() + "/" + tplCacheFileName
	}

	err = ioutil.WriteFile(cfn, fileData, 0644)
	if err != nil {
		log.Fatal(err, cfn)
		return err
	}

	return nil
}

// readCache Opens the cache file and reads contents
func (tc *tplCache) readCache() error {
	fl := tc.location.Location() + "/" + tplCacheFileName
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

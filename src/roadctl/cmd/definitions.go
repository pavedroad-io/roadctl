// definitions structure

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

//Note: above and below blank lines required for golint.
//Related to required documentation format for packages.

import (
	"errors"
	"fmt"
	"os"
	//	"reflect"
	"regexp"
	"strings"
)

// Tables require additional documentaion.
type Tables struct {
	// TableName is the name of the table to create.
	//   This is really just a sub-object
	//   It could be meta-data, user, blah...
	//
	TableName string `yaml:"table-name"`

	// TableType
	//   JSONB    Supported today
	//   SQLTable Future
	//
	//   If TableRole is Secondary, TableType is ignored
	//   if present and the table will be added as a sub
	//   strucuter
	//
	TableType string `yaml:"table-type"`

	// ParentTable
	//   Nest this table under the named parent
	ParentTable string `yaml:"parent-tables"`

	// A list of table columes or object attributes
	Columns []struct {
		// Name of the column
		Name string `yaml:"name"`

		// Modifiers to apply when marshalling
		// omitempty or string
		//
		Modifiers string `yaml:"modifiers"`

		// MappedName to use when marshalling
		// If empty, map to lowercase of name
		//
		MappedName string `yaml:"mapped-name"`

		// Contraints such as required
		// Valid swagger 2.0 validation
		//
		Constraints string `yaml:"constraints"`

		// TODO: need to map this to array/map for validation
		// Type such as boolen, string, float, or integer
		//   valid types: Stick to JSON for now
		//     string
		//     number
		//     integer
		//     boolen
		//     time
		//     null

		Type string `yaml:"type"`
	} `yaml:"columns"`
}

// Community require additional documentaion.
type Community struct {
	CommunityFiles []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
		Src  string `yaml:"src"`
		Md5  string `yaml:"md5,omitempty"`
	} `yaml:"community-files"`
	Description string `yaml:"description"`
}

// Info is holds API help information.
type Info struct {
	APIVersion    string `yaml:"api-version"`
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Organization  string `yaml:"organization"`
	ReleaseStatus string `yaml:"release-status"`
	Version       string `yaml:"version"`
}

// Dependencies require additional documentaion.
type Dependencies []struct {
	Command          string      `yaml:"command"`
	Comments         string      `yaml:"comments"`
	DockerCockroahdb interface{} `yaml:"docker-cockroahdb,omitempty"`
	Image            string      `yaml:"image"`
	Name             string      `yaml:"name"`
	Ports            []struct {
		External string `yaml:"external"`
		Internal string `yaml:"internal"`
	} `yaml:"ports"`
	Volumes     []interface{} `yaml:"volumes"`
	DockerKafka interface{}   `yaml:"docker-kafka,omitempty"`
	Topics      []string      `yaml:"topics,omitempty"`
}

// Maintainer require additional documentaion.
type Maintainer struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
	Slack string `yaml:"slack"`
	Web   string `yaml:"web"`
}

// ProjectFiles require additional documentaion.
type ProjectFiles struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Src         string `yaml:"src"`
}

//Badges require additional documentaion.
type Badges struct {
	Enable bool   `yaml:"enable"`
	Link   string `yaml:"link"`
	Name   string `yaml:"name"`
}

// ConfigurationFile require additional documentaion.
type ConfigurationFile struct {
	ArtifactsDir string `yaml:"artifacts-dir"`
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Src          string `yaml:"src"`
}

// Integrations require additional documentaion.
type Integrations struct {
	Badges []Badges `yaml:"badges,omitempty"`
	Name   string   `yaml:"name"`
	// TODO: Needs to be more generic
	//
	SonarCloudConfig struct {
		// A sonarcloud access token
		Login string `yaml:"login"`
		// Project key should be same a the name
		Key     string `yaml:"key"`
		Options struct {
			Badges   []Badges `yaml:"badges"`
			Coverage struct {
				Enable bool   `yaml:"enable"`
				Report string `yaml:"report"`
			} `yaml:"coverage"`
			GoSec struct {
				Enable bool   `yaml:"enable"`
				Report string `yaml:"report"`
			} `yaml:"go-sec"`
			Lint struct {
				Enable bool   `yaml:"enable"`
				Report string `yaml:"report"`
			} `yaml:"lint"`
		} `yaml:"options"`
	} `yaml:"sonar-cloud-config,omitempty"`
	ConfigurationFile ConfigurationFile `yaml:"configuration-file,omitempty"`
}

// Project require additional documentaion.
type Project struct {
	Description  string         `yaml:"description"`
	Dependencies Dependencies   `yaml:"dependencies"`
	License      string         `yaml:"license"`
	Maintainer   Maintainer     `yaml:"maintainer"`
	ProjectFiles []ProjectFiles `yaml:"project-files"`
	Integrations []Integrations `yaml:"integrations"`
}

type tplDef struct {
	TableList []Tables  `yaml:"tables"`
	Community Community `yaml:"community"`
	Info      Info      `yaml:"info"`
	Project   Project   `yaml:"project"`
}

// Define constants for error types
// iota starts at 0, next constant is iota + 1
const (
	INVALIDCOLUMNTYPE = iota
	INVALIDCONTENT
	INVALIDTABLENAME
	NOPARENT
	NOMAPPEDNAME
	INVALIDTABLETYPE
	INVALIDCOLUMNNAME
	INVALIDCONSTRAINT
	INVALIDMODIFIER
	INVALIDARRAY
	UNKOWN
)

// tblDefError
// structure for returning table defintion errors
//
type tblDefError struct {
	errorType    int
	errorMessage string
	tableName    string
	nextError    *tblDefError
}

// ErrList is a list of table definition error messages.
// Processing is done on all errors instead of exiting for
// every error.
var ErrList *tblDefError //Prior comment for exported format.

// LastErr is the last error message on ErrList.
var LastErr *tblDefError //Prior comment for exported format.

// tlbDefError
// implements error.Error() interface
//
func (d *tblDefError) Error() string {
	e := fmt.Sprintf("Table: %v, Error number: %v, %v\n",
		d.tableName,
		d.errorType,
		d.errorMessage)

	return e
}

func (d *tplDef) setErrorList(msgNum int, msg string, tName string) {
	e := tblDefError{
		errorType:    msgNum,
		errorMessage: msg,
		tableName:    tName,
	}
	if LastErr == nil {
		LastErr = &e
		ErrList = &e
	} else {
		LastErr.nextError = &e
		LastErr = &e

	}
}

// tplTableItem
type tplTableItem struct {
	// The name of this table
	Name string

	Root bool
	// Children: a list of table items containing
	//           child tables
	Children []*tplTableItem
}

// devineOrder: Determine primary table and its
// children.  Generate an error if no primary is
// found or more than one primary is found
// TODO: The above logic needs to be specific to
//       the type of service build built
func (d *tplDef) devineOrder() tplTableItem {
	ptName := ""

	// Get primary table and make sure it is the only primary
	x := d.findTables(ptName)
	if len(x) == 0 {
		fmt.Println("No primary table found")
		os.Exit(-1)
	} else if len(x) > 1 {
		fmt.Println("More than primary table found: ", len(x))
		os.Exit(-1)
	} else {
		pt := x[0]
		ptName = pt.Name
	}

	d.addChildren(&x[0])
	//d.walkOrder(x[0])

	return x[0]
}

// walkOrder: Given a parent, print out all of its
//   children
func (d *tplDef) walkOrder(item tplTableItem) {

	if len(item.Children) > 0 {
		for _, v := range item.Children {
			fmt.Printf("Parent: %v Child: %v\n", item.Name, v.Name)
			d.walkOrder(*v)
		}
	} else {
		fmt.Printf("Parent: %v Child: no children\n", item.Name)
	}
	return
}

// addChildren: Add children to a parent, then
//  add any children they may have recursivley
func (d *tplDef) addChildren(parent *tplTableItem) {

	c := d.findTables(parent.Name)

	if len(c) == 0 {
		return
	}

	for _, v := range c {
		parent.Children = append(parent.Children, &v)
		d.addChildren(&v)
	}
	/*
		     Removed as duplicate.
		     No noted reason for this additional processing.

			         for _, v := range c {
					parent.Children = append(parent.Children, &v)
					d.addChildren(&v)
				}
	*/
	return

}

// findTables: Find primary parent table, or
// children for a given table
func (d *tplDef) findTables(parent string) []tplTableItem {
	rlist := []tplTableItem{}
	//	tlist := d.tables()

	for _, t := range d.TableList {
		if t.ParentTable == parent {
			c := make([]*tplTableItem, 0, 20)
			var isRoot = false
			if parent == "" {
				isRoot = true
			}
			newrec := tplTableItem{t.TableName, isRoot, c}
			//fmt.Println(newrec)
			rlist = append(rlist, newrec)
		}
	}

	return rlist
}

// tables(): return a pointer(a copy?) to definitions Tables
func (d *tplDef) tables() []Tables {
	return d.TableList
}

// Search for Table by name
func (d *tplDef) tableByName(name string) (Tables, error) {
	e := Tables{}
	for _, v := range d.TableList {
		if v.TableName == name {
			return v, nil
		}
	}
	return e, errors.New("table not found")
}

func (d *tplDef) badges() []Badges {
	var badgelist []Badges
	for _, rec := range d.Project.Integrations {
		if len(rec.Badges) > 0 {
			badgelist = append(badgelist, rec.Badges...)
		}
		if strings.ToLower(rec.Name) == "sonarcloud" &&
			len(rec.SonarCloudConfig.Options.Badges) > 0 {
			badgelist = append(badgelist, rec.SonarCloudConfig.Options.Badges...)
		}
	}
	return badgelist
}

func (d *tplDef) findIntegration(name string) Integrations {
	for _, rec := range d.Project.Integrations {
		if strings.ToLower(rec.Name) == strings.ToLower(name) {
			return rec
		}
	}
	a := Integrations{}
	return a
}

func (d *tplDef) BadgesToString() string {
	badges := ""
	for _, b := range d.badges() {
		if b.Enable == true {
			badges += b.Link + "\n"
		}
	}
	return badges
}

//Valide the table(s) definition

func (d *tplDef) Validate() *tblDefError {

	// Do all tables and report all potential errors
	for _, t := range d.tables() {
		// Metadata validation
		d.validateTableMetaData(t)

		//Table definition validation
		d.validateTableColumns(t)

	}
	return ErrList
}

// validateTableMetaData
// name not blank, type is supported, parent table exists
// Table name only have allowed characters
// **All table names should be unique (not case sensitive)
// Table name length
func (d *tplDef) validateTableMetaData(t Tables) *tblDefError {

	var validTypes = []string{"JSONB"}
	const maxLen = 60

	// Make sure table name is set
	if t.TableName == "" {
		d.setErrorList(INVALIDTABLENAME, "Missing table name", "")
	} else {

		if len(t.TableName) > maxLen {
			e := fmt.Sprintf("Table name length cannot be greater than  %v", maxLen)
			d.setErrorList(INVALIDTABLENAME, e, t.TableName)

		}

		// Use simple regex until more sepcific requirements
		// and security review.
		// Must be modified to support specific target databases
		// Looked at specifications for Cockroachdb
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]*$`, t.TableName)

		if !matched {
			d.setErrorList(INVALIDTABLENAME, "Bad table name: ["+t.TableName+"]", t.TableName)

		}
	}

	// If a parent is specified make sure it exists
	if t.ParentTable != "" {
		_, e := d.tableByName(t.ParentTable)
		if e != nil {
			d.setErrorList(NOPARENT, "Parent table not found: ["+t.ParentTable+"]", t.TableName)

		}
	} else {
		//do checks only related to parent tables
		// Make sure it is a valid type
		isValidType := false
		for _, m := range validTypes {
			if strings.ToUpper(t.TableType) == m {
				isValidType = true
				break
			}
		}

		if !isValidType {
			d.setErrorList(INVALIDTABLETYPE, "Bad table type: ["+t.TableType+"]", t.TableName)

		}

	}
	return ErrList
}

func (d *tplDef) validateTableColumns(t Tables) *tblDefError {
	var validColTypes = []string{
		"string",
		"number",
		"integer",
		"boolean",
		"time",
		"null",
		"uuid",
	}
	var convName string

	// validate:
	//  - Name *
	//  - Modifiers
	//  - MappedName
	//  - Constraints
	//  - Type *
	//

	for _, v := range t.Columns {
		//Check the column name
		if v.Name == "" {
			d.setErrorList(INVALIDCOLUMNNAME, "Missing column name", t.TableName)
		} else {

			matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]*$`, v.Name)

			if !matched {
				d.setErrorList(INVALIDCOLUMNNAME, "Bad column name: ["+v.Name+"]", t.TableName)

			}
		}
		//Check the column types
		convName = strings.ToLower(v.Type)

		if !isStringInList(validColTypes, convName) {
			d.setErrorList(INVALIDCOLUMNTYPE, "Invalid column type: ["+v.Type+"]", t.TableName)

		}

		//Check the mapped Name
		//This wouldn't be an error if it is the functionality
		//required.
		if v.MappedName == "" {
			v.MappedName = strings.ToLower(v.Name)
			//d.setErrorList(NOMAPPEDNAME, "Column : ["+v.MappedName+"] had no mapped name.", t.TableName)
		}

	}
	return ErrList
}

// isStringInList
// Should move to a generic function list
//
func isStringInList(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

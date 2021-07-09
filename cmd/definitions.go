// Package cmd from gobra
//   types and methods for blueprint definitions file
package cmd

import (
	"errors"
	"fmt"
	"os"

	//	"reflect"
	"regexp"
	"strings"
)

var StaticDefinitionFileVersion = "v0.9.0"

// Tables structure for user defined tables that need
// to be generated
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
	//   structure
	//
	TableType string `yaml:"table-type"`

	// ParentTable
	//   Nest this table under the named parent
	ParentTable string `yaml:"parent-tables"`

	// A list of table columns or object attributes
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

		// Constraints such as required
		// Valid swagger 2.0 validation
		//
		Constraints string `yaml:"constraints"`

		// TODO: need to map this to array/map for validation
		// Type such as boolean, string, float, or integer
		//   valid types: Stick to JSON for now
		//     string
		//     number
		//     integer
		//     boolean
		//     time
		//     null

		Type string `yaml:"type"`
	} `yaml:"columns"`
}

// Community files to be included
//   For example, CONTRIBUTING.md
type Community struct {
	CommunityFiles []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
		Src  string `yaml:"src"`
		Md5  string `yaml:"md5,omitempty"`
	} `yaml:"community-files"`
	Description string `yaml:"description"`
}

// Info defines information about the services and organization
type Info struct {
	APIVersion    string `yaml:"api-version"`
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Organization  string `yaml:"organization"`
	ReleaseStatus string `yaml:"release-status"`
	Version       string `yaml:"version"`
}

// Dependencies that this service requires
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
	Topics      []Topic       `yaml:"topics,omitempty"`
}

// Topic defines a Kafka topic and its partition and replication counts
type Topic struct {
	// Value is the topic name
	Value string `json:"value"`
	// Partitions is the number of partions
	Partitions int `json:"partitions"`
	// ReleaseStatus is the replication factor for this topic
	Replication int `json:"replication"`
}

// Maintainer contact information
type Maintainer struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
	Slack string `yaml:"slack"`
	Web   string `yaml:"web"`
}

// ProjectFiles blueprint files to be included
type ProjectFiles struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Src         string `yaml:"src"`
}

// Badge are links with graphics to be included in
// doc/service.html file.  These go to CI test results
type Badge struct {
	Enable bool   `yaml:"enable"`
	Link   string `yaml:"link"`
	Name   string `yaml:"name"`
}

// ConfigurationFile where the configuration can be found
type ConfigurationFile struct {
	ArtifactsDir string `yaml:"artifacts-dir"`
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Src          string `yaml:"src"`
}

// KubeConfig
type KubeConfig struct {
	// Namespace to use when constructing URLs
	Namespace string `yaml:"namespace"`

	// Liveness endpoint name for k8s checks
	Liveness string `yaml:"liveness"`

	// Readiness endpoint name for k8s checks
	Readiness string `yaml:"readiness"`

	// Metrics endpoint name for k8s checks
	Metrics string `yaml:"metrics"`

	// Management endpoint
	Management string `yaml:"management"`

	// Explain endpoint
	Explain string `yaml:"explain"`
}

// Integrations CI/CD tools
type Integrations struct {
	Badges  []string `yaml:"shields,omitempty"`
	Name    string   `yaml:"name"`
	Enabled bool     `yaml:"enable"`
	// TODO: Needs to be more generic
	//
	SonarCloudConfig struct {
		// A sonarcloud access token
		Login string `yaml:"login"`
		// Project key should be same as the name
		Key     string `yaml:"key"`
		Options struct {
			Badges   []string `yaml:"shields"`
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

// Query parameter
type queryParm struct {
	Name        string `json:"name"`
	DataType    string `json:"datatype"`
	Description string `json:"description"`
}

type httpMethod struct {
	Method string      `json:"method"`
	QP     []queryParm `json:"qp"`
}

// Endpoints
type endPoint struct {
	Name    string       `json:"name"`
	Methods []httpMethod `json:"methods"`
}

// Project information
type Project struct {
	TLD           string         `yaml:"top_level_domain"`
	Description   string         `yaml:"description"`
	Dependencies  Dependencies   `yaml:"dependencies"`
	License       string         `yaml:"license"`
	SchedulerName string         `yaml:"scheduler_name"`
	Maintainer    Maintainer     `yaml:"maintainer"`
	ProjectFiles  []ProjectFiles `yaml:"project-files"`
	Integrations  []Integrations `yaml:"integrations"`
	Kubernetes    KubeConfig     `yaml:"kubernetes"`
	Endpoints     []endPoint     `yaml:"endpoints"`
	Loggers       []Logger       `yaml:"loggers"`
	Blocks        []Block        `yaml:"blocks"`
}

type Go struct {
	DependencyManager string `yaml:"dependency-manager"`
}

// Core capabilities to enable / configure
type Core struct {
	Loggers []Logger `yaml:"loggers"`
}

type bpDef struct {
	DefinitionFileVersion string
	DefinitionFile        string    `yaml:"definitionFile"`
	TableList             []Tables  `yaml:"tables"`
	Community             Community `yaml:"community"`
	Info                  Info      `yaml:"info"`
	Project               Project   `yaml:"project"`
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
	INVALIDSERVICENAME
	UNKOWN
)

// tblDefError
// structure for returning table definition errors
//
type tblDefError struct {
	errorType    int
	errorMessage string
	tableName    string
}

// ErrList is a list of table definition error messages.
// Processing is done on all errors instead of exiting for
// every error.
var ErrList []tblDefError //Prior comment for exported format.

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

func (d *bpDef) setErrorList(msgNum int, msg string, tName string) {

	e := tblDefError{
		errorType:    msgNum,
		errorMessage: msg,
		tableName:    tName,
	}

	ErrList = append(ErrList, e)

}

// bpTableItem
type bpTableItem struct {
	// The name of this table
	Name string

	Root bool
	// Children: a list of table items containing
	//           child tables
	Children []*bpTableItem

	// IsList means this is a []Type talbe
	IsList bool
}

// devineOrder: Determine primary table and its
// children.  Generate an error if no primary is
// found or more than one primary is found
// TODO: The above logic needs to be specific to
//       the type of service build built
func (d *bpDef) devineOrder() bpTableItem {
	// ptName "" means this table does not have
	//   a parent
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
func (d *bpDef) walkOrder(item bpTableItem) {

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
//  add any children they may have recursively
func (d *bpDef) addChildren(parent *bpTableItem) {

	c := d.findTables(parent.Name)

	if len(c) == 0 {
		return
	}

	for _, v := range c {
		parent.Children = append(parent.Children, &v)
		d.addChildren(&v)
	}
	return

}

// findTables: Find primary parent table, or
// children for a given table
func (d *bpDef) findTables(parent string) []bpTableItem {
	rlist := []bpTableItem{}
	//	tlist := d.tables()

	for _, t := range d.TableList {
		if t.ParentTable == parent {
			c := make([]*bpTableItem, 0, 20)
			var isRoot = false
			if parent == "" {
				isRoot = true
			}
			var isList = false
			if strings.ToLower(t.TableType) == "list" {
				isList = true
			}
			newrec := bpTableItem{t.TableName, isRoot, c, isList}
			rlist = append(rlist, newrec)
		}
	}

	return rlist
}

// tables(): return a pointer(a copy?) to definitions Tables
func (d *bpDef) tables() []Tables {
	return d.TableList
}

// Search for Table by name
func (d *bpDef) tableByName(name string) (Tables, error) {
	e := Tables{}
	for _, v := range d.TableList {
		if v.TableName == name {
			return v, nil
		}
	}
	return e, errors.New("table not found")
}

func (d *bpDef) badges() []string {
	var badgelist []string
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

func (d *bpDef) findIntegration(name string) Integrations {
	for _, rec := range d.Project.Integrations {
		if strings.ToLower(rec.Name) == strings.ToLower(name) {
			return rec
		}
	}
	a := Integrations{}
	return a
}

func (d *bpDef) BadgesToString() string {
	badges := ""
	/*
		for _, b := range d.badges() {
			fmt.Println(b)
				if b.Enable == true {
					badges += b.Link + "\n"
				}
		}
	*/
	return badges
}

//Valid the table(s) definition, and other YAML defaults needed
//for anticipated execution

func (d *bpDef) Validate() (errCount int) {

	const badMicroserviceName = "yourMicroserviceName"
	const pavedroadSonarTestOrg = "acme-demo"

	// TODO(sgayle): This defined an empty error message and assigned it
	// to LastErr.  That caused and error to allways be returned

	//Doing YAML default test first
	//Blueprint default microservice should be changed
	//if sonar cloud testing under pavedroadSonarTestOrg is
	//needed.

	if (d.Info.Name == defMicroserviceName) && (d.Info.Organization == pavedroadSonarTestOrg) {
		d.setErrorList(INVALIDSERVICENAME, "Sonar cloud microservice name change expected.", "")
	}

	// Do all tables and report all potential errors
	for _, t := range d.tables() {
		// Metadata validation
		e := d.validateTableMetaData(t)
		if e > 0 {
			return e
		}

		// Column validation
		e = d.validateTableColumns(t)
		if e > 0 {
			return e
		}
	}

	return errCount
}

// validateTableMetaData
// name not blank, type is supported, parent table exists
// Table name only have allowed characters
// **All table names should be unique (not case sensitive)
// Table name length
func (d *bpDef) validateTableMetaData(t Tables) (errCount int) {

	var validTypes = []string{"JSONB", "OBJECT", "LIST"}
	const maxLen = 60

	// Make sure table name is set
	if t.TableName == "" {
		d.setErrorList(INVALIDTABLENAME, "Missing table name", "")
		errCount++
	} else {

		if len(t.TableName) > maxLen {
			e := fmt.Sprintf("Table name length cannot be greater than  %v", maxLen)
			d.setErrorList(INVALIDTABLENAME, e, t.TableName)
			errCount++

		}

		// Use simple regex until more specific requirements
		// and security review.
		// Must be modified to support specific target databases
		// Looked at specifications for Cockroachdb
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]*$`, t.TableName)

		if !matched {
			d.setErrorList(INVALIDTABLENAME, "Bad table name: ["+t.TableName+"]", t.TableName)
			errCount++
		}
	}

	// Make sure it is a valid type
	if t.ParentTable == "" {
		isValidType := false
		for _, m := range validTypes {
			if strings.ToUpper(t.TableType) == m {
				isValidType = true
				break
			}
		}

		if !isValidType {
			d.setErrorList(INVALIDTABLETYPE, "Bad table type: ["+t.TableType+"]", t.TableName)
			errCount += 1
		}
	}

	// If a parent is specified make sure it exists
	if t.ParentTable != "" {
		_, e := d.tableByName(t.ParentTable)
		if e != nil {
			d.setErrorList(NOPARENT, "Parent table not found: ["+t.ParentTable+"]", t.TableName)
			errCount += 1
		}
	}

	return errCount
}

func (d *bpDef) validateTableColumns(t Tables) (errCount int) {
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
			errCount += 1
		} else {

			matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]*$`, v.Name)

			if !matched {
				d.setErrorList(INVALIDCOLUMNNAME, "Bad column name: ["+v.Name+"]", t.TableName)
				errCount += 1
			}
		}

		//Check the column types
		m := mappedTypes{}
		convName = strings.ToLower(v.Type)

		if !m.validInputType(convName) {
			d.setErrorList(INVALIDCOLUMNTYPE, "Invalid column type: ["+v.Type+"]", t.TableName)
			errCount += 1
		}

		//Check the mapped Name
		//This wouldn't be an error if it is the functionality
		//required.
		if v.MappedName == "" {
			v.MappedName = strings.ToLower(v.Name)
			d.setErrorList(NOMAPPEDNAME, "Column : ["+v.MappedName+"] had no mapped name.", t.TableName)
			errCount += 1
		}

	}
	return errCount
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

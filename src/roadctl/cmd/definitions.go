// definitions structure

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Tables struct {
	// TableName is the name of the table to create
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

type Community struct {
	CommunityFiles []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
		Src  string `yaml:"src"`
		Md5  string `yaml:"md5,omitempty"`
	} `yaml:"community-files"`
	Description string `yaml:"description"`
}

type Info struct {
	APIVersion    string `yaml:"api-version"`
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Organization  string `yaml:"organization"`
	ReleaseStatus string `yaml:"release-status"`
	Version       string `yaml:"version"`
}
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
type Maintainer struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
	Slack string `yaml:"slack"`
	Web   string `yaml:"web"`
}
type ProjectFiles struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Src         string `yaml:"src"`
}

type Badges struct {
	Enable bool   `yaml:"enable"`
	Link   string `yaml:"link"`
	Name   string `yaml:"name"`
}

type ConfigurationFile struct {
	ArtifactsDir string `yaml:"artifacts-dir"`
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Src          string `yaml:"src"`
}

type Integrations struct {
	Badges []Badges `yaml:"badges,omitempty"`
	Name   string   `yaml:"name"`
	// TODO: Needs to be more generic
	//
	SonarCloudConfig struct {
		// A sonarcloud access token
		Login string `yaml:login`
		// Project key should be same a the name
		Key     string `yaml:key`
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
	Commuity  Community `yaml:"community"`
	Info      Info      `yaml:"info"`
	Project   Project   `yaml:"project"`
}

// Define constants for error types
//
const (
  INVALIDTYPE = iota
  INVALIDCONTENT
  BADTABLENAME
  NOPARENT
  BADTABLETYPE
  UNKOWN 
)

// tblDefError
// structure for returning table defintion errors
//
type tblDefError struct {
  errorType  int
  errorMessage string
}

// tlbDefError
// implements error.Error() interface
//
func (d *tblDefError) Error() string {

  e := fmt.Sprintf("Error number %v, %v\n", 
    d.errorType,
    d.errorMessage)

  return e
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

	// Get primary table and make sure this is only one
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
	} else {
		for _, v := range c {
			parent.Children = append(parent.Children, &v)
			d.addChildren(&v)
		}
	}

	return
}

// findTables: Given a parent, see if they have children
func (d *tplDef) findTables(parent string) []tplTableItem {
	rlist := []tplTableItem{}
	tlist := d.tables()

	for _, t := range tlist {
		if t.ParentTable == parent {
			c := make([]*tplTableItem, 0, 20)
			var isRoot bool = false
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

// tables(): return a pointer to definitions Tables
func (d *tplDef) tables() []Tables {
	return d.TableList
}

// table
//
func (d *tplDef) table(name string) (Tables, error) {
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

func (d *tplDef) Validate() []tblDefError {
  var errList  []tblDefError

  // validateTables
  for _, t := range d.tables() {
    // Metadata validation
    e := d.validateTableMetaData(t)
    if len(e) > 0  {
      for _, x := range e {
         errList = append(errList, x)
      }
    }

    // Column validation
    e = d.validateTableColumns(t)
    if len(e) > 0  {
      for _, x := range e {
         errList = append(errList, x)
      }
    }

  }
  return errList
}

// validateTableMetaData
// name not blank, type is supported, parent table exists
//
func (d *tplDef) validateTableMetaData(t Tables) []tblDefError {
  var eList []tblDefError
  var validTypes = []string{"JSONB"}

  // Make sure table name is set
  if t.TableName == "" {
    e := tblDefError{
      errorType: BADTABLENAME,
      errorMessage: "Bad or missing table name: [" + t.TableName + "]",
    }
    eList = append(eList, e)
  }

  // Make sure it is a valid type
  isValidType := false
  for _, m := range validTypes {
    if t.TableType == m {
      isValidType = true
      break;
    }
  }

  if !isValidType {
    e := tblDefError{
      errorType: BADTABLETYPE,
      errorMessage: "Bad table type: [" + t.TableType +"]",
    }
    eList = append(eList, e)
  }

  // If a parent is specified make sure it exists
  if t.ParentTable != "" {
    _, e := d.table(t.ParentTable)
    if e != nil {
      em := tblDefError{
        errorType: NOPARENT,
        errorMessage: "Parent table not found: [" + t.ParentTable +"]",
      }
      eList = append(eList, em)
    }
  }
  return eList
}

func (d *tplDef) validateTableColumns(t Tables) []tblDefError {
  var eList []tblDefError
  var validColTypes = []string{
    "string",
    "number",
    "integer",
    "boolean",
    "time",
    "uuid",
    "null",
  }

  // validate:
  //  - Name
  //  - Modifiers
  //  - MappedName
  //  - Constraints
  //  - Type
  // 
  for _,v := range t.Columns {
    if !isStringInList(validColTypes, v.Type){
    e := tblDefError{
      errorType: INVALIDTYPE,
      errorMessage: "Invalid column type: [" + v.Type +"]",
    }
    eList = append(eList, e)

    }
  }
  return eList
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

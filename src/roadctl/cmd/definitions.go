// definitions structure

package cmd

import (
	"fmt"
	"os"
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
		Modifiers string `yaml:"modifiers"`

		// MappedName to use when marshalling
		MappedName string `yaml:"mapped-name"`

		// Contraints such as required
		Contraints string `yaml:"contraints"`

		// TODO: need to map this to array/map for validation
		// Type such as boolen, string, float, or integer
		//   valid types: Stick to JSON for now
		//     string
		//     number
		//     integer
		//     boolen
		//     null

		Type string `yaml:"type"`
	} `yaml:"columns"`
}

type tplDef struct {
	TableList []Tables `yaml:"tables"`
	Project   struct {
		MaintainerEmail string `yaml:"maintainer-email"`
		Integrations    []struct {
			Config struct {
				Options struct {
					Coverage struct {
						Report string `yaml:"report"`
						Enable bool   `yaml:"enable"`
					} `yaml:"coverage"`
					Lint struct {
						Report string `yaml:"report"`
						Enable bool   `yaml:"enable"`
					} `yaml:"lint"`
					GoSec struct {
						Report string `yaml:"report"`
						Enable bool   `yaml:"enable"`
					} `yaml:"go-sec"`
				} `yaml:"options"`
				ConfigurationFile struct {
					Path         string `yaml:"path"`
					Name         string `yaml:"name"`
					ArtifactsDir string `yaml:"artifacts-dir"`
					Src          string `yaml:"src"`
				} `yaml:"configuration-file"`
			} `yaml:"config,omitempty"`
			Name        string `yaml:"name"`
			Path        string `yaml:"path,omitempty"`
			Description string `yaml:"description,omitempty"`
			Src         string `yaml:"src,omitempty"`
		} `yaml:"integrations"`
		Maintainer   string `yaml:"maintainer"`
		Dependencies []struct {
			DockerInfo struct {
				Image   string `yaml:"image"`
				Command string `yaml:"command"`
				Ports   []struct {
					Internal string `yaml:"internal"`
					External string `yaml:"external"`
				} `yaml:"ports"`
				Comments string `yaml:"comments"`
				Volumes  []struct {
					Path  string  `yaml:"path"`
					Mount float64 `yaml:"mount"`
				} `yaml:"volumes"`
			} `yaml:"docker-info"`
			Name string `yaml:"name"`
		} `yaml:"dependencies"`
		MaintainerSlack string `yaml:"maintainer-slack"`
		License         string `yaml:"license"`
		Description     string `yaml:"description"`
	} `yaml:"project"`
	Name          string `yaml:"name"`
	Organization  string `yaml:"organization"`
	APIVersion    string `yaml:"api-version"`
	Version       string `yaml:"version"`
	ID            string `yaml:"id"`
	ReleaseStatus string `yaml:"release-status"`
	Community     struct {
		CommunityFiles []struct {
			Path string `yaml:"path"`
			Name string `yaml:"name"`
			Src  string `yaml:"src"`
			Md5  string `yaml:"md5,omitempty"`
		} `yaml:"community-files"`
		Description string `yaml:"description"`
	} `yaml:"community"`
}

// tplTableItem
type tplTableItem struct {
  // The name of this table
	Name     string

  // Children: a list of table items containing
  //           child tables
	Children []*tplTableItem
}

// devineOrder: Determine primary table and its
// children.  Generate an error if no primary is
// found or more than one primary is found
// TODO: The above logic needs to be specific to
//       the type of service build built
func (d *tplDef) devineOrder() {
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
  fmt.Println("Dump tables")
  d.walkOrder(x[0])

  return
}

// walkOrder: Given a parent, print out all of its
//   children
func (d *tplDef)walkOrder(item tplTableItem) {

    if len(item.Children) > 0 {
      for _, v := range item.Children {
        fmt.Printf("Parent: %v Child: %v\n",item.Name, v.Name)
        d.walkOrder(*v)
      }
      } else {
        fmt.Printf("Parent: %v Child: no children\n",item.Name)
      }
  return
}

// addChildren: Add children to a parent, then 
//  add any children they may have recursivley
func (d *tplDef)addChildren(parent *tplTableItem) {

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
			newrec := tplTableItem{t.TableName, c}
			rlist = append(rlist, newrec)
		}
	}

	return rlist
}

// talbes(): return a pointer to definitions Tables
func (d *tplDef) tables() []Tables {
	return d.TableList
}

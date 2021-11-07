// cmd
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var t Tables

// SQLToJSONBlock produce test data for a SQL table
var SQLToJSONBlock = Block{
	APIVersion: "v1beta",
	Kind:       "FunctionBlock",
	ID:         "http://io.pavedroad/test/sql/jsonGenerator",
	Family:     "pavedroad/test/data/generation",
	Metadata: Metadata{
		Labels: []string{
			"pavedroad",
			"test",
			"sql",
			"table",
			"generator",
			"outputJSON",
		},
		Tags: []string{
			"pavedroad",
			"test",
			"sql",
			"table",
			"generator",
			"outputJSON",
		},
		Information: BlockInformation{
			Description: "Generate test data based on an SQL table structure",
			Title:       "Generate test data based on an SQL table structure",
			Contact: Contact{
				Author:       "John Scharber",
				Organization: "PavedRoad",
				Email:        "support@pavedroad.io",
				Website:      "www.pavedroad.io",
				Support:      "pavedroad-io.slack.com",
			},
		},
	},
	UsageRights: UsageRights{
		TermsOfService: "As is",
		Licenses:       "Apache 2",
		AccessToken:    "",
	},
	Language:      "go",
	HomeDirectory: "dev",
	Environment:   "dev",
	FunctionMappings: []FunctionItem{
		{
			Function:           createPostData,
			OutputFileName:     "post.json",
			OutputType:         "OutputFile",
			ExecutePermissions: false,
			Description:        "JSON content for a new POST",
		},
		{
			Function:           createPutData,
			OutputFileName:     "put.json",
			OutputType:         "OutputFile",
			ExecutePermissions: false,
			Description:        "JSON content for a PUT updated",
		},
	},
}

// createPostData read a table for a given definitions file
// and create sample test data in JSON format
// example, json, err := createPostData(defs bpDef)
func createPostData(opt ...interface{}) (result []byte, err error) {
	var jsonString string
	defs, ok := opt[0].(bpDef)
	if !ok {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Cast to dpDef failed: %v", ok)}
		return nil, ne.WrappedError()
	}
	order := defs.devineOrder()
	//bpAddJSON(order, defs, &jsonString)
	bpAddJSON(order, defs, &jsonString)

	//Make it pretty
	var pj bytes.Buffer
	err = json.Indent(&pj, []byte(jsonString), "", "\t")
	if err != nil {
		log.Fatal("Failed to generate json data with ", jsonString)
	}
	return pj.Bytes(), nil
}

// createPut read a table for a given definitions file
// and creates sample test data in JSON format
// example, json, err := createPut(defs bpDef, postData []byte)
func createPutData(opt ...interface{}) (result []byte, err error) {

	if len(opt) != 2 {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Usage: (bpDef, JSON( are required")}
		return nil, ne.WrappedError()
	}
	def, ok := opt[0].(bpDef)
	if !ok {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Cast to bpDef for JSON failed: %v", ok)}
		return nil, ne.WrappedError()
	}

	jsonData, ok := opt[1].([]byte)
	if !ok {
		ne := bpError{Type: ErrorGeneric,
			Err: fmt.Errorf("Cast to []byte for JSON failed: %v", ok)}
		return nil, ne.WrappedError()
	}

	var post map[string]interface{}

	err = json.Unmarshal(jsonData, &post)

	if putData, err := createPostData(opt[0]); err != nil {
		return nil, err
	} else {
		item := def.devineOrder()

		var put map[string]interface{}
		err = json.Unmarshal(putData, &put)
		if err != nil {
			log.Println(err)
		}

		id := strings.ToLower(item.Name) + "uuid"
		// TODO: add test for id before just assumming it is there
		b := strings.Replace(string(putData),
			put[id].(string),
			post[id].(string), 1)

		return []byte(b), nil
	}

	return nil, nil
}

// Package sonarcloud API
package fossa

// Constants used in URI and query parameters
const (
	// Default api server for SonarCloud
	DefaultScheme = "https://"
	// Default api server for SonarCloud
	DefaultHost = "app.fossa.com"

	// Project URI
	Project = "/projects/custom%%2B9819%%2F%s"

	// Project URI
	ProjectAPI = "/api/projects/custom%%2B9819%%2F%s.svg"
	//	ProjectAPI = "/api/projects/custom+9819/%s.svg"

	// Query parameter strings
	Type = "?type=%s"
	Ref  = "?ref=badge_%s"
)

// Valid metric types
const (
	Badge = iota
	Shield
	Small
	Large
)

// MetricName maps integer constants to expected FOSSA string name
var MetricName = map[int]string{
	0: "badge",
	1: "shield",
	2: "small",
	3: "large",
}

// Valid output types
const (
	MarkDown = iota
	HTML
	Link
)

var URL = DefaultScheme + DefaultHost + Project + Type
var APIURL = DefaultScheme + DefaultHost + ProjectAPI + Type
var RefURL = DefaultScheme + DefaultHost + Project + Ref
var RefAPIURL = DefaultScheme + DefaultHost + ProjectAPI + Ref

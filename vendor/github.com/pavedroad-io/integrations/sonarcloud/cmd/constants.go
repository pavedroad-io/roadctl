// Package sonarcloud API
package sonarcloud

// Constants used in URI and query parameters
const (
	// Default api server for SonarCloud
	DefaultScheme = "https://"
	// Default api server for SonarCloud
	DefaultHost = "sonarcloud.io"
	// Default api version
	DefaultAPI = "/api"

	// ProjectSearch URI
	ProjectSearch = DefaultAPI + "/projects/search"

	// ProjectCreate URI
	ProjectCreate = DefaultAPI + "/projects/create"

	// ProjectDelete URI
	ProjectDelete = DefaultAPI + "/projects/delete"

	// TokenSearch URI
	TokenSearch = DefaultAPI + "/user_tokens/search"

	// TokenCreate URI
	TokenCreate = DefaultAPI + "/user_tokens/generate"

	// TokenRevoke URI
	TokenRevoke = DefaultAPI + "/user_tokens/revoke"

	// BadgeMetric URI
	BadgeMetric = DefaultAPI + "/project_badges/measure"

	// QualityGate URI
	QualityGate = DefaultAPI + "/project_badges/quality_gate"

	// Query parameter strings
	Branch       = "branch=%s"
	Login        = "login=%s"
	Metric       = "metric=%s"
	Name         = "name=%s"
	Organization = "organization=%s"
	Project      = "project=%s"
	Projects     = "projects=%s"

	// KeyPrefix to append to SonarCloud Key to ensure uniqueness
	KeyPrefix = "PavedRoad_"

	// Testing constants
	projectKey       = "test123"
	projectName      = "Test project 123"
	orgname          = "acme-demo"
	visibility       = "public"
	tokenName        = "userTestToken123"
	badServerAddress = "localhost:3000"
)

// Valid metric types
const (
	Bugs = iota
	CodeSmells
	Coverage
	DuplicatedLinesDensity
	Ncloc
	SqaleRating
	AlertStatus
	ReliabilityRating
	SecurityRating
	SqaleIndex
	Vulnerabilities
)

// MetricName maps integer constants to expected Sonar string name
var MetricName = map[int]string{
	0:  "bugs",
	1:  "code_smells",
	2:  "coverage",
	3:  "duplicated_lines_density",
	4:  "ncloc",
	5:  "sqale_rating",
	6:  "alert_status",
	7:  "reliability_rating",
	8:  "security_rating",
	9:  "sqale_index",
	10: "vulnerabilities",
}

const (
	testErrorMsg      = "Expected err to be nil Got %v\n"
	testErrorMsgValue = "Expected err to be '%v' Got %v\n"
	testMarshalFail   = "Unmarshal failed Got %v\n"
	errorIs           = "Error is "
	contentType       = "Content-Type"
	contentLength     = "Content-Length"
	wwwForm           = "application/x-www-form-urlencoded"
)

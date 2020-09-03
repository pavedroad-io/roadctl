// Package sonarcloud API wrapper
//
// Provides a basic wrapper for the SonarCloud API
//   Support is limited to:
//		projects
//		tokens
//		metrics
//		quality gates
package sonarcloud

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ProjectSearchResponse for a GET / search on projects
type ProjectSearchResponse struct {
	// Paging object
	Paging PagingObject `json:"paging"`

	// Components list
	Components []ComponentsObject `json:"components"`
}

// PagingObject used for response that are limited by size
type PagingObject struct {

	// Index page number
	Index int `json:"pageIndex"`

	// Size elements on this page
	Size int `json:"pageSize"`

	// Total number of pages
	Total int `json:"total"`
}

// ComponentsObject structure
type ComponentsObject struct {
	// Organization name
	Organization string `json:"organization"`

	// Key for access this object
	Key string `json:"key"`

	// Name for display
	Name string `json:"name"`

	// Qualifier is type of component
	Qualifier string `json:"qualifier"`

	// LastAnalysisDate time of last CI
	LastAnalysisDate string `json:"lastAnalysisDate"`

	// Revision hash
	Revision string `json:"revision"`
}

// NewProject Used to create a new project
//
type NewProject struct {
	// Organization (required) is a valid SonarCloud organization
	Organization string `json:"organization"`

	// Name friendly name for display
	Name string `json:"name"`

	// Project is the SonarCloud Key
	Project string `json:"project"`

	// Visibility (optional) private or public
	Visibility string `json:"visibility"`
}

// NewProjectResponse includes wrapper "project"
// This sucks we can't use the same structure but they
// change field names
type NewProjectResponse struct {
	Project NewProjectResponseObject `json:"project"`
}

// NewProjectResponseObject Object include in response
//
type NewProjectResponseObject struct {
	// Key
	Key string `json:"key"`

	// Name friendly name for display
	Name string `json:"name"`

	// Qualifier is the SonarCloud component type
	Qualifier string `json:"qualifier"`

	// Visibility (optional) private or public
	Visibility string `json:"visibility"`
}

// SonarCloudClient type and methods used for accessing SonarCloud API
type SonarCloudClient struct {
	//   Client is an http.Client created when New() is called
	Client *http.Client

	//   Host is the default host, sonarcloud.io by default
	Host string

	//   APIversion is the api prefix to use in API calls /api
	APIVersion string

	//   Token used to authenticate
	Token string

	// connectino string
	URI string
}

// NewTokenResponse holds response from user_tokens/generate
type NewTokenResponse struct {
	// Login name token is for
	Login string `json:"login"`

	// Name of the token
	Name string `json:"name"`

	// The Token
	Token string `json:"token"`

	// CreatedAt date and time of creation
	CreatedAt string `json:"createdAt"`
}

// GetTokenResponse user_tokens/search returns a user and
// a list of their tokens
type GetTokenResponse struct {
	// Login name of user
	Login string `json:"login"`

	// List of tokens
	Tokens []GetTokenItem `json:"userTokens"`
}

// GetTokenItem items returned in a token search
type GetTokenItem struct {
	// Name of the token
	Name string `json:"name"`

	// CreatedAt date and time
	CreatedAt string `json:"createdAt"`

	// LastConnectionDate date and time token was last used
	// Only updated hourly
	LastConnectionDate string `json:"lastConnectionDate"`
}

// SonarCloudError
type sonarCloudError struct {
	errNumber int
	errMsg    string
}

func (e *sonarCloudError) Error() string {
	return fmt.Sprintf("Err: %v, %v\n", e.errNumber, e.errMsg)
}

// New create a new Sonarcloud client
//   expample New.(sondarcloudclient, token)
//   token is a valid sonarcloud user token
//	 if must have admin access
func (c *SonarCloudClient) New(token string, timeoutSeconds int) error {

	c.Client = &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	if c.Host == "" {
		c.Host = DefaultHost
	}

	if c.APIVersion == "" {
		c.APIVersion = DefaultAPI
	}

	if token == "" {
		a := sonarCloudError{errNumber: -1, errMsg: "Token is require"}
		return &a
		//		panic(a.Error())
	}

	c.Token = token

	// https://token@host
	c.URI = fmt.Sprintf("%s%s@%s", DefaultScheme, c.Token, c.Host)

	return nil
}

// GetProject read a project from sonar cloud
// TODO make this a vardic function taking a list of project names
func (c *SonarCloudClient) GetProject(org, name string) (*http.Response, error) {
	options := "?"
	options += fmt.Sprintf(Projects, name)
	options += fmt.Sprintf("&"+Organization, org)

	url := c.URI + ProjectSearch + options

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(resp, err)
	}

	return resp, nil
}

// CreateProject create a new project on SonarCloud
//   example: CreateProject.(Project)
//   Create a new SonarCloud project usinig p
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *SonarCloudClient) CreateProject(p NewProject) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", p.Name)

	// Project names for none default organization are global
	// Add a prefix to avoid naming conflicts
	// TODO: Make this configurable by the end user
	data.Set("project", KeyPrefix+p.Project)
	data.Set("organization", p.Organization)
	data.Set("visibility", p.Visibility)

	url := c.URI + ProjectCreate
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	req.Header.Add(contentType, wwwForm)
	req.Header.Add(contentLength, strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(rsp, err)
	}

	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	return rsp, nil
}

// DeleteProject delete a SonarCloud project
//   Example: DeleteProject.(projectKey)
//   Delete a new SonarCloud project usinig p
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *SonarCloudClient) DeleteProject(p string) (*http.Response, error) {

	// Project names for none default organization are global
	pk := KeyPrefix + p
	data := url.Values{}
	data.Set("project", pk)

	url := c.URI + ProjectDelete
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	req.Header.Add(contentType, wwwForm)
	req.Header.Add(contentLength, strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(rsp, err)
	}

	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	return rsp, nil
}

// CreateToken  create a new SonarCloud token
//   example:  CreateToken.(tn string)
//   Create a new SonarCloud token with the name tn
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *SonarCloudClient) CreateToken(tn string) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", tn)

	url := c.URI + TokenCreate
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	req.Header.Add(contentType, wwwForm)
	req.Header.Add(contentLength, strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(rsp, err)
	}

	// return errmsg in the body as the error
	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}
	return rsp, nil
}

// RevokeToken revoke the given token
//   example: RevokeToken(tn string)
//   Revoke a SonarCloud token with the name tn
func (c *SonarCloudClient) RevokeToken(tn string) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", tn)

	url := c.URI + TokenRevoke
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	req.Header.Add(contentType, wwwForm)
	req.Header.Add(contentLength, strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	// There was a problem with the connection
	if err != nil {
		return HandleHTTPClientError(rsp, err)
	}

	// There was a problem with the payload
	// return errmsg in the body as the error
	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	return rsp, nil
}

// GetTokens get a list of tokens
// example: GetTokens(login)
// Return a list of tokens for the current user
// or if login is specified use it
//
func (c *SonarCloudClient) GetTokens(name string) (*http.Response, error) {
	var url string
	if name != "" {
		options := "?"
		options += fmt.Sprintf(Login, name)
		url = c.URI + TokenSearch + options
	} else {
		url = c.URI + TokenSearch
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(resp, err)
	}

	return resp, nil
}

// NameToEnum return integer value of enum or -1 if not found
func (c *SonarCloudClient) NameToEnum(name string) (enum int) {
	for i, v := range MetricName {
		if v == strings.ToLower(name) {
			return i
		}
	}
	return -1
}

// GetMetric return a badge for a given metric
// example: GetMetric(metric, branch string)
// Return an SVG badge for inclussion in HTML
//
// 	metric (required) is one of the following constants:
//  	Bugs
//		CodeSmells
//		Coverage
//		DuplicatedLinesDensity
//		Ncloc
//		SqaleRating
//		AlertStatus
//		ReliabilityRating
//		SecurityRating
//		SqaleIndex
//
//  project  (required) project to produce bade for
//  branch (optional) a long living branch
//
func (c *SonarCloudClient) GetMetric(metric int, project, branch string) (*http.Response, error) {
	var url string
	options := "?"
	options += fmt.Sprintf(Metric, MetricName[metric])
	options += fmt.Sprintf("&"+Project, project)
	if branch != "" {
		options += fmt.Sprintf("&"+Branch, branch)
	}
	url = c.URI + BadgeMetric + options

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(resp, err)
	}

	return resp, nil
}

// GetQualityGate return badge for a SonarCloud quality gate
// example: GetQualityGate(project string) (*http.Response, error)
// Return an SVG badge for inclusion in HTML
// 	project (required) is a valid project name
//
func (c *SonarCloudClient) GetQualityGate(project string) (*http.Response, error) {
	var url string
	options := "?"
	options += fmt.Sprintf(Project, project)
	url = c.URI + QualityGate + options

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HandleHTTPClientError(nil, err)
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return HandleHTTPClientError(resp, err)
	}

	return resp, nil
}

// HandleHTTPClientError returns (*http.Response, error)
// Prints error message and returns response and error
func HandleHTTPClientError(rsp *http.Response, err error) (*http.Response, error) {
	fmt.Println(errorIs, err)
	return rsp, err
}

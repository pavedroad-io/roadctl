![PavedRoad, Inc.](assets/images/pavedroad_black_230x36.png)

[![Build Status](https://travis-ci.org/pavedroad-io/roadctl.svg?branch=travisSetup2)](https://travis-ci.org/pavedroad-io/roadctl)[![Go Report Card](https://goreportcard.com/badge/github.com/pavedroad-io/roadctl)](https://goreportcard.com/report/github.com/pavedroad-io/roadctl)[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=alert_status)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=ncloc)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=sqale_index)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=security_rating)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_roadctl&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=pavedroad-io_roadctl)[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B9819%2Fgit%40github.com%3Apavedroad-io%2Froadctl.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B9819%2Fgit%40github.com%3Apavedroad-io%2Froadctl.git?ref=badge_shield)

# Roadctl
Roadctl is a command line interface for:

- Creating microservices, CRDs, and serverless functions from low-code blueprints
- Managing a preconfigured CI/CD pipeline
- Controlling deployment options

This overview covers roadctl syntax, describes the command operations, and provides common examples. For details about each command, including all the supported flags and subcommands, see the roadctl reference documentation. For installation instructions see [installing roadctl](http://www.pavedroad.io/roadctl/install.md).

## Syntax
Use the following syntax to run roadctl commands from your terminal window:

`$ roadctl [command] [TYPE] [NAME] [flags]`

Where command, TYPE, NAME, and flags are:

- command: Specifies the operation that you want to perform on one or more resources, for example create, get, describe, delete
- TYPE: Specifies the resource type. Resource types are case insensitive and you can specify the singular, or plural forms. For example, the following commands produce the same output:

```
    $ roadctl get blueprint blueprint1
    $ roadctl get blueprints blueprint1
```
- NAME: Specifies the name of the resource. Names are case-sensitive. If the name is omitted, details for all resources are displayed, for example roadctl get blueprints
- flags: Specify roadctl global options or command specific options

## Resource Types
The following table includes a list of all the supported resource types:

| Resource Type | Description |
|:--------------|:------------|
|builders| Manage build pipelines to produce compiled results such as code and CSS|
|environments| Manages target environments to deploy into, for example dev, test, staging|
|packagers|Create images/containers along with manifest for docker, docker-compose, and Kubernetes|
|taggers|Support tagging artifacts, images, and releases|
|tests|Manage unit, function, benchmarks, and container tests|
|blueprints|Allow applications to be built from predefined blueprints such as API gateways or data managers|
|integrations|Allows you to tailor the preconfigured integrations|
|artifacts|Manage development artifacts such as logs and code coverage|
|providers|Manage local and cloud providers you want to deploy to|
|deployments|Manage deployment strategies for various environments|

## Output Options
By default, roadctl outputs the results of a command as text.  However,
you can control that by using the _--format_ option:

    --format text
    --format json
    --format yaml

## Requirements Before You Begin

### FOSSA Requirements
Support for FOSSA is integrated in the generated Makefile.
Please set up your FOSSA account and add the following
variable to your environment before executing make:

`$ export FOSSA_API_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXX"`

NOTE: You can disable FOSSA in the integrations section of your
definitions file you don't need it.

### SonarCloud Requirements
PavedRoad utilizes SonarCloud to champion quality code in this project.
Please set up your SonarCloud account and the following
environment variable to take full advantage of the CLI:

`$ export SONARCLOUD_TOKEN="#########"`

NOTE: You can disable SonarCloud in the integrations section of your
definitions file you don't need it.

## Examples: Common Operations

### Initialize a Local Template Repository
The following command populates available blueprints on your local hard
drive.  By default, they are placed $HOME/.pavedroad.d/blueprints.

`$ roadctl init`

###
### Print a List of Available Templates
The roadctl command can print all available blueprints or verify an individual blueprints name.
The output includes the blueprints type, name and its release status.

`$ roadctl get blueprints`

**Or**

`$ roadctl get blueprints name`

Example output for the above command with no name specified:

```
Template Type   Name                 Release Status
crd             kubebuilder          incubation
microservices   workerPool           ga
microservices   ux                   ga
microservices   service              ga
microservices   datamgr              ga
microservices   gateway              ga
microservices   subscriber           experimental
serverless      go-knative           ga
serverless      go-open-faas         ga
```
The meaning for each type of release status is as follows:

| Release Status | Meaning |
|:---------------|:------- |
| ga             | For general availability|
| incubation     | For blueprints working towards ga|
| experimental   | Not stable or work in progress or simple examples|

### Create and Edit a Template Definitions File
The blueprints definitions file allow you to tailor your application to your requirements, such as:

- Define fields and structures to create
- Specify community files
- Tailor the initial integrations that get included
- Set organizational and project information like: license, company name, or project description

Each blueprints comes with a default definitions file you can use as a beginning point.

Use the describe command to create your definitions file:

`$ roadctl describe blueprints datamgr > myNewService.yaml`

Partial example output:

```
tables:
- columns:
  - constraints: ""
    mapped-name: id
    modifiers: ""
    name: id
    type: string
  - constraints: ""
    mapped-name: title
    modifiers: ""
    name: title
    type: string
  - constraints: ""
    mapped-name: updated
    modifiers: ""
    name: updated
    type: time
  - constraints: ""
    mapped-name: created
    modifiers: ""
    name: created
    type: time
  parent-tables: ""
  table-name: users
  table-type: jsonb
- columns:
  - constraints: ""
    mapped-name: id
    modifiers: ""
    name: id
    type: string
  parent-tables: users
  table-name: metadata
- columns:
  - constraints: ""
    mapped-name: key
    modifiers: ""
    name: key
    type: string
  parent-tables: metadata
  table-name: test
```

Then edit it using vi:

`$ vi myNewService.yaml`

### To See Valid Contents of a Template
Use the explain command to learn the valid syntax for the named blueprints is:

`$ roadctl explain blueprints datamgr`

Partial example output:

```
Name: blueprints

DESCRIPTION:
Templates provide a low-code environment for serverless, CRD, and microservices.
The roadctl CLI uses the blueprints skaffold combined code generation to create your
application, CI, and test framework.

FIELDS:
name <string>
     A user friendly name for this blueprints
api-version <string>
     API version used to generate it
version <string>
     Object data model version
id <string>
     UUID that uniquely identified a combination of api-verion + version
     This UUID is immutable for the above combination
```

### Generate an Application
Will build your source code, test cases, integrations, documentation, and CI pipeline

`$ roadctl create blueprints datamgr -f myNewService.yaml`

To compile and invoke the CI/CD pipeline enter:

`$ make`

### Build Defaults
The defaults are to execute lint, go sec, go vet, go test,
and the FOSSA and SonarCloud scanners.

- The artifacts hold the results for each command

The build generates the following components:

-  Go source code
-  Dependency graph
-  Swagger API specification in docs/api.json
-  HTML service documentation in docs/myServiceName.html
-  HTML API documentation in docs/api.html
-  Dependency management with dep and insure all includes are present
-  Dockerfile
-  docker-compose files
-  Kubernetes manifests
-  Skaffold configuration file
-  Go test harness

The deployment options are:

- Push image to local microk8s instance
- Push service and run service to local microk8s instance

### Make Options

Make can be executed with the following options:

| Make Options | Meaning |
|:-------------|:------- |
|`$ make`|Check and compile the code|
|`$ make compile`|Just compile the code|
|`$ make check`| Just execute the test suite|
|`$ make deploy`| Package and deploy code to local k8s cluster|
|`$ make fmt`| Rewrite code in go formatting|
|`$ make simplify`| Format and simplify the code use|

## Initializing the Template Repository

A GitHub repository stores the PavedRoad blueprints.

As of version 0.6, the git clone command is the default method for creating
blueprints repositories.  The default is to checkout the "release" branch:

`$ roadctl init`

To use a different branch, use the _--branch_ option:

`$ roadctl init --branch <branch-name>`

### Templates Location

The default location for the blueprints is \$HOME/.pavedroad.d/blueprints.

This can be changed by setting the PR_TEMPLATE_DIR environment variable, or via the
roadctl command line with the _--blueprints flag.

## GitHub Authentication

For backward compatibility, use the _--api_ option with authentication.

The GitHub API enforces rate limits that may affect your ability to download blueprints.  Authenticated users have significantly higher rate limits.  You can provide GitHub authentication using HTTP basic authentication or an OAUTH2 access token.

### From the Command Line
`$ roadctl init --api --password XXXXXXX --user YYYYYYY`

**Or**

`$ roadctl init --token #######`

### Using Environment Variables
`$ export GH_ACCESS_TOKEN="#########"`

**Or**

```
$ export GH_USER_NAME="#########"
$ export GH_USER_PASSWORD="#########"
```

### Or a Combination
```
$ roadctl init --api --user YYYYYYY
$ export GH_USER_PASSWORD="#########"
```

Package and deploy code to local k8s cluster:

`$ make deploy`

Rewrite code in go formatting:

`$ make fmt`

Or to format and simplify the code use:

`$ make simplify`

## Project Status

The project is an early preview. We realize that it's going to take a village to arrive at the vision of a multi-cloud control plane, and we wanted to open this up early to get your help and feedback. Please see the [Roadmap](/ROADMAP.md) for details on what we are planning for future releases.

## Official Releases

Official releases of PavedRoad can be found here:
[Official Releases](https://github.com/pavedroad-io/pavedroad/releases).

Please note that it is **strongly recommended** that you use the official releases
of PavedRoad, as unreleased versions from the master branch are subject to
changes and incompatibilities that will not be supported in the official releases.
Builds from the master branch can have functionality changed and even removed
at any time without compatibility support and without prior notice.

## Links to More Information

### Community Meeting
This project meets on a regular basis: [TBD Community Meeting](https://zoom.us/j/7886774843).
### Getting Help
For contact information or to report a bug see [Support](/SUPPORT.md).
### How to Contribute
For guidelines on contributions see [Contributing](/CONTRIBUTING.md).
### Code of Conduct
This project follows this [Code of Conduct](/CODE_OF_CONDUCT.md).
### License
This project is licensed under the following [License](/LICENSE).

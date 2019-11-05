<p align="center"><img src="https://github.com/pavedroad-io/kevlar-repo/blob/master/assets/images/banner.png" alt="PavedRoad.io"></p>

# Roadctl
Roadclt is a command-line interface for:

-	Creating microservices, CRDs, and serverless functions from low-code templates
-	Managing a pre-configured CI/CD pipeline
-	Controlling deployment options

This overview covers roadctl syntax, describes the command operations, and provides common examples. For details about each command, including all the supported flags and subcommands, see the roadctl reference documentation. For installation instructions see [installing roadctl](http://www.pavedroad.io/roadctl/install.md).

## Syntax
Use the following syntax to run roadctl commands from your terminal window:

$ roadctl [command] [TYPE] [NAME] [flags]

Where command, TYPE, NAME, and flags are:

-	command: Specifies the operation that you want to perform on one or more resources, for example create, get, describe, delete.
-	TYPE: Specifies the resource type. Resource types are case-insensitive and you can specify the singular, or plural forms. For example, the following commands produce the same output:

     `$ roadctl get template template1`

     `$ roadctl get templates template1`

-	NAME: Specifies the name of the resource. Names are case-sensitive. If the name is omitted, details for all resources are displayed, for example roadctl get templates.

## Resource types
The following table includes a list of all the supported resource types
### builders
Manage build pipelines to produce compiled results such; code and css
### packagers
Create images/containers along with manifest for docker, docker-compose, and kubernetes
### taggers
Support tagging artifacts, images, and releases
### tests
Manage unit, function, benchmarks, and container tests
### templates
Allow applications to be built from pre-defined templates such as API gateways or data managers

### integrations
Allows you to tailor the preconfigured integrations

### artifacts
Manage development artifacts such as logs and code coverage.
### providers
Manage local and cloud providers you want to deploy to
### deployments
Manage deployment strategies for various environments.

## Output options
By default, roadclt outputs the results of a command as text.  However, 
you can control that by using the --format option:

--format text

--format json

--format yaml

## Examples: Common operations

### Initialize a local template repository
The following command populates available templates on your local hard drive.  By default, they are placed in your current directory in “.templates.”  The prefix “.” Is necessary to tell compilers to ignore its contents.

`$ roadctl get templates --init`

### 
### Print a list of available templates
The output includes the template name and its release status.  Release status is one of the following:

```
```

| Release | Meaning |
| :---------- | :-------------------------------- |
| ga         | For general availability|
| incubation | For templates working towards ga|
| experimental | Not stable or work in progress or simple examples|

`$ roadctl get templates`

Or

`$ roadctl get templates name`

Example output

```
$ roadctl get templates

Template Type   Name                 Release Status
crd             kubebuilder          incubation
microservices   datamgr              ga
microservices   service              ga
microservices   ux                   ga
microservices   gateway              ga
microservices   subscriber           experimental
microservices   ml                   incubation
serverless      go-open-faas         ga
serverless      go-knative           ga
```

### Create and edit a template definitions file
Template definitions file allow you to tailor your application to your requirements, such as:

-	Define fields and structures to create
-	Specify community files
-	Tailor the initial integrations that get included
-	Set organizational and project information like; license, company name, or project description

Each template comes with a default definitions file you can use as a beginning point.

Use the describe command to create your definitions file

`$ roadctl describe templates datamgr > myNewService.yaml`

Then edit it using vi

`$ vi myNewService.yaml`

Example

```
  tables:
  - columns:
    - {constraints: '', mapped-name: id, modifiers: '', name: id, type: string}
    - {constraints: '', mapped-name: updated, modifiers: '', name: updated, type: string}
    - {constraints: '', mapped-name: created, modifiers: '', name: created, type: string}
    parent-tables: ''
    table-name: users
    table-type: jsonb
```

### To see valid contents for a template
Use the explain command to learn the valid syntax for the named templates is.

```
$ roadctl explain templates datamgr

Example
Name: templates

DESCRIPTION:
Templates provide a low-code environment for serverless, crd, and Microservices.
The roadctl CLI uses the template scaffold combined code generation to create
your application, CI, and test framework.
  
  FIELDS:
  name <string>
       A user friendly name for this template
  api-version <string>
       API version used to generate it
  version <string>
       Object data model version
  id <string>
```

### Generate your application
Will build your source code, test cases, integrations, documentation, and CI pipeline

`$ roadctl create templates --template datamgr --definition my myNewService.yaml`

To compile and invoke the CI/CD pipeline enter:

`$ make`

### Build defaults
Execute lint, go sec, go test, and sonar scanner

-	The artifacts hold the results for each command

Generate the the following components:

-  Go source code
-  Dependency graph
-  Swagger API specification in docs/api.json
-  HTLM service documentation in docs/myServiceName.html
-  HTML API documentation in docs/api.html
-  Dependency management with dep and insure all includes are present
-  Dockerfile
-  docker-compose files
-  Kubernetes manifests
-  Skaffold configuration file
-  go test harness

Deployment options

- Push image to local microk8s instance
- Push service and run service to local microk8s instance

### Make options
Check and compile the code

`$ make`

Just compile the code

`$ make compile`

Just execut the test suite

`$ make check`

Package and deploy code to local k8s cluster

`$ make deploy`

Rewrite code in go formatting

`$ make fmt`

Or to format and simplify the code use

`$ make simplify`

## GitHub Authentication
A GitHub repository stores PavedRoad templates.  The GitHub API enforces rate limits that may affect your ability to download templates.  Authenticated users have significantly higher rate limits.  You can provide GitHub authentication using HTTP basic authentication or an OAUTH2 access token.

### From the command line
`$ roadctl get templates --init --password XXXXXXX --user YYYYYYY`

**or**

`$ roadctl get templates --token #######`

### Using environment variables
`$ export ACCESS-TOKEN="#########"`

**or**

`$ export USER-NAME="#########"`
`$ export USER-PASSWORD="#########"`

### Or a combination
`$ roadctl get templates --init  --user YYYYYYY`

`$ export USER-PASSWORD="#########"`

## GitHub Authentication
A GitHub repository stores PavedRoad templates.  The GitHub API enforces rate limits that may affect your ability to download templates.  Authenticated users have significantly higher rate limits.  You can provide GitHub authentication using HTTP basic authentication or an OAUTH2 access token.

### From the command line
roadctl get templates --init --password XXXXXXX --user YYYYYYY
**or**
roadctl get templates --token #######

### Using environment variables
export ACCESS-TOKEN="#########"
**or**
export USER-NAME="#########"
export USER-PASSWORD="#########"

### Or a combination
roadctl get templates --init  --user YYYYYYY
export USER-PASSWORD="#########"

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

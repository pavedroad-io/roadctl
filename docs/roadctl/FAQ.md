# FAQ

## not within known GOPATH/src error

Go projects have a strict directory structure.  By default, Go creates a $HOME/go/src directory for you to place your packages and programs.  Your GOPATH will point here unless you change it.
```bash
$ go env GOPATH
GOPATH="/home/jscharber/go"
```
If you want to place your project in a different directory, you must include a parent src directory and change your GOPATH.

```bash
mkdir -p $HOME/eng/src/myproject
cd $HOME/eng/src/myproject
export $GOPAPTH=(cd ../../;pwd)
```
Omitting the "src" directory or not setting the GOPATH results in the "not within known GOPATH/src" error message.

## Starting or Stopping microk8s prompts for a password

This is the expected behavior.  Under the hood, microk8s calls sudo when starting or stopping.

## GitHub rate limiting error

A GitHub repository stores PavedRoad templates. The GitHub API enforces rate limits that may affect your ability to download templates. Authenticated users have significantly higher rate limits. You can provide GitHub authentication using HTTP basic authentication or an OAUTH2 access token.

Providing GitHub credentials will avoid these errors.  See [https://github.com/pavedroad-io/roadctl](https://github.com/pavedroad-io/roadctl) for more information.

### Sample error
```bash
$ roadctl get templates datamgr --init

Initializing template repository
GET https://api.github.com/repos/pavedroad-io/templates/contents/microservices/ga/datamgr/manifests/kubernetes/template-service.yaml: 403 API rate limit exceeded for 208.96.177.111

GET https://api.github.com/repos/pavedroad-io/templates/contents/microservices/ga/datamgr/sonarcloud.sh: 403 API rate limit exceeded for 208.96.177.111.

GET https://api.github.com/repos/pavedroad-io/templates/contents/microservices/ga/datamgr/templateApp.go: 403 API rate limit exceeded for 208.96.177.111.

GET https://api.github.com/repos/pavedroad-io/templates/contents/microservices/ga/datamgr/templateDoc.go: 403 API rate limit exceeded for 208.96.177.111.
```

### Solution
```bash
$ roadctl get templates --init --password XXXXXXX --user YYYYYYY
```
or

```bash
$ export GITHUB_ACCESS_TOKEN="#########"
```

# Go Modules

As of 0.5.1alpha, we have migrated to Go modules and
vendoring.  This change requires making changes to
your build environment.  Also, the Makefile is changed
to support modules, and the Gopkg.lock and Gopkg.toml
removed.

## Steps to migrate

### Check your directory structure
Your code must be in the following directory structure

github.com/pavedroad-io/roadctl

This is required for modules and Travis CI

### Initialize go modules

```bash
go mod init github.com/pavedroad-io/roadctl
```
### Rebuild your vendor directory

```bash
rm -R vendor
go mod vendor
```

## Makefile changes
`dep ensure` has been replaced with `go mod vendor`
`go build` now includes the `mod=vendor`
This eliminates go downloading dependencies

## New environment variable

`GO111MODULE=on` is required 

Without this option, Go attempts to use dep.





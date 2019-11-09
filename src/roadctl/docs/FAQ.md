# FAQ

## not within known GOPATH/src error

Go projects have a strict directory structure.  By default, Go creates a $HOME/go/src directory for you to place your packages and programs.  Your GOPATH will point here unless you change it.
```bash
$ go env | grep GOPATH
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

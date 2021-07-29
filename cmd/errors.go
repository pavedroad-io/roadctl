package cmd

import (
	"fmt"
	"runtime"
)

// Constants for error numbers
const (
	ErrorGeneric int = iota
	ErrorNotFound
)

// errorMessage map number to messages
var errorMessage = map[int]string{
	ErrorGeneric:  "Generic(%v:%v) Error(%v)\n",
	ErrorNotFound: "File not found(%v:%v) Error(%v)\n",
}

// bpError standardized error type for roadctl
type bpError struct {

	// Type categorizes types and messages used to report
	Type int

	// Error a wrapped error message
	Err error
}

// WrappedError returns an error
// %w in new in go 1.13
func (e *bpError) WrappedError() error {
	_, f, l, _ := runtime.Caller(0)
	err := fmt.Errorf("%v: %v Error(%w)\n", f, l, e.Err)
	return err
}

// Error returns an error message as a string
// old error handling
func (e *bpError) Error() string {
	_, f, l, _ := runtime.Caller(0)
	msg := fmt.Sprintf(errorMessage[e.Type], f, l, e.Err.Error())
	return msg
}

// New just in case we decide to now export Type and Err
func (e *bpError) New(t int, err error) {
	e.Type = t
	e.Err = err
}

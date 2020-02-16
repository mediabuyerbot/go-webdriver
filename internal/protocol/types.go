package protocol

import (
	"fmt"
)

type StackFrame struct {
	FileName   string
	ClassName  string
	MethodName string
	LineNumber int
}

type CommandError struct {
	StatusCode int
	ErrorType  string
	Message    string
	Screen     string
	Class      string
	StackTrace []StackFrame
}

func (e CommandError) Error() string {
	m := e.ErrorType
	if m != "" {
		m += ": "
	}
	if e.StatusCode == -1 {
		m += "status code not specified"
	} else if str, found := statusCode[e.StatusCode]; found {
		m += str + ": " + e.Message
	} else {
		m += fmt.Sprintf("unknown status code (%d): %s", e.StatusCode, e.Message)
	}
	return m
}

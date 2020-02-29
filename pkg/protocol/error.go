package protocol

import (
	"encoding/json"
	"fmt"
)

type StackFrame struct {
	FileName   string
	ClassName  string
	MethodName string
	LineNumber int
}

// Error represents a WebDriver protocol error.
type Error struct {
	Code       string                 `json:"error"`
	Message    string                 `json:"message"`
	StackTrace string                 `json:"stacktrace"`
	Data       map[string]interface{} `json:"data"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", e.Message, e.Code)
}

func parseError(respStatusCode int, resp *Response) error {
	cmdErr := new(Error)
	// if error not JSON
	if err := json.Unmarshal(resp.Value, cmdErr); err != nil {
		cmdErr.Code = httpStatusCode(respStatusCode)
		cmdErr.Message = string(resp.Value)
		return cmdErr
	}

	var (
		isNullValue                     = string(resp.Value) == "null"
		isErrorWithStatusAndPayloadNull = resp.Status > 0 && isNullValue
		isErrorWithStatusAndPayload     = resp.Status > 0
		isErrorWithoutAndPayloadNull    = resp.Status == 0 && isNullValue
	)

	switch {
	// {"status": 19, "value": null}
	case isErrorWithStatusAndPayloadNull:
		cmdErr.Code = httpStatusCode(respStatusCode)
		statusMsg, ok := statusCode[resp.Status]
		if ok {
			cmdErr.Message = statusMsg
		}

	// {"status": 19, "value": {"error":"error code", "message": "error message"}}
	case isErrorWithStatusAndPayload:
		if len(cmdErr.Code) == 0 {
			cmdErr.Code = httpStatusCode(respStatusCode)
		}
		if len(cmdErr.Message) == 0 {
			cmdErr.Message = cmdErr.Code
			sm, ok := statusCode[resp.Status]
			if ok {
				cmdErr.Message = sm
			}
		}

	// {"value": null} httpStatus >= 400
	case isErrorWithoutAndPayloadNull:
		cmdErr.Code = httpStatusCode(respStatusCode)
		cmdErr.Message = cmdErr.Code

	default:
		if len(cmdErr.Code) == 0 {
			cmdErr.Code = httpStatusCode(respStatusCode)
		}
	}
	return cmdErr
}

package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidResponse = errors.New("protocol: invalid response")

// Error represents a WebDriver protocol error.
type Error struct {
	Code          string                 `json:"error"`
	Message       string                 `json:"message"`
	RawStacktrace string                 `json:"stacktrace"`
	Data          map[string]interface{} `json:"data"`

	Stacktrace []string
}

func (e Error) Error() string {
	return fmt.Sprintf("ErrorCode:%s, Message:%s",
		e.Code,
		e.Message,
	)
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
		isNullValue                        = string(resp.Value) == "null"
		isErrorWithStatusAndPayloadNull    = resp.Status > 0 && isNullValue
		isErrorWithStatusAndPayload        = resp.Status > 0
		isErrorWithoutStatusAndPayloadNull = resp.Status == 0 && isNullValue
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
	case isErrorWithoutStatusAndPayloadNull:
		cmdErr.Code = httpStatusCode(respStatusCode)
		cmdErr.Message = cmdErr.Code

	default:
		if len(cmdErr.Code) == 0 {
			cmdErr.Code = httpStatusCode(respStatusCode)
		}
	}
	return cmdErr
}

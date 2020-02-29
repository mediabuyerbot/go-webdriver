package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_ParseError(t *testing.T) {

	tests := []struct {
		name           string
		r              *Response
		e              *Error
		httpStatusCode int
	}{
		{
			name: "chromeBrowser with status=44, httpResponseCode=200, value=JSON, overrideError=true",
			r: &Response{
				SessionID: "123",
				Status:    19,
				Value:     []byte(`{"error":"error code", "message": "error message"}`),
			},
			e: &Error{
				Code:          "error code",
				Message:       "error message",
				RawStacktrace: "",
				Data:          nil,
			},
			httpStatusCode: 200,
		},

		{
			name: "chromeBrowser with status=19, httpResponseCode=200, value=JSON, overrideMessage=true",
			r: &Response{
				SessionID: "123",
				Status:    19,
				Value:     []byte(`{"message":"error message", "data": {"x":"1"}}`),
			},
			e: &Error{
				Code:          httpStatusCode(200),
				Message:       "error message",
				RawStacktrace: "",
				Data: map[string]interface{}{
					"x": "1",
				},
			},
			httpStatusCode: 200,
		},

		{
			name: "chromeBrowser with status=19, httpResponseCode=200, value=JSON, overrideCode=true",
			r: &Response{
				SessionID: "123",
				Status:    19,
				Value:     []byte(`{"error":"error", "data": {"x":"1"}}`),
			},
			e: &Error{
				Code:          "error",
				Message:       StatusText(XPathLookupErrorStatusCode),
				RawStacktrace: "",
				Data: map[string]interface{}{
					"x": "1",
				},
			},
			httpStatusCode: 200,
		},

		{
			name: "chromeBrowser httpResponseCode=400, value={}",
			r: &Response{
				SessionID: "123",
				Value:     []byte(`{}`),
			},
			e: &Error{
				Code:          httpStatusCode(400),
				RawStacktrace: "",
				Data:          nil,
			},
			httpStatusCode: 200,
		},

		{
			name: "chromeBrowser httpResponseCode=400, value=JSON",
			r: &Response{
				SessionID: "123",
				Status:    0,
				Value:     []byte(`{"error":"error", "message":"error", "data": {"x":"1"}}`),
			},
			e: &Error{
				Code:          "error",
				Message:       "error",
				RawStacktrace: "",
				Data: map[string]interface{}{
					"x": "1",
				},
			},
			httpStatusCode: 400,
		},

		{
			name: "chromeBrowser httpResponseCode=400, value=NULL",
			r: &Response{
				SessionID: "123",
				Status:    0,
				Value:     []byte(`null`),
			},
			e: &Error{
				Code:          httpStatusCode(400),
				Message:       httpStatusCode(400),
				RawStacktrace: "",
				Data:          nil,
			},
			httpStatusCode: 400,
		},

		{
			name: "chromeBrowser httpResponseCode=400, value=error text",
			r: &Response{
				SessionID: "123",
				Status:    0,
				Value:     []byte(`error msg`),
			},
			e: &Error{
				Code:          httpStatusCode(400),
				Message:       "error msg",
				RawStacktrace: "",
				Data:          nil,
			},
			httpStatusCode: 400,
		},

		{
			name: "chromeBrowser with status=19, httpResponseCode=200, value=null",
			r: &Response{
				SessionID: "123",
				Status:    TimeoutStatusCode,
				Value:     []byte(`null`),
			},
			e: &Error{
				Code:          httpStatusCode(200),
				Message:       StatusText(TimeoutStatusCode),
				RawStacktrace: "",
				Data:          nil,
			},
			httpStatusCode: 200,
		},
	}

	var err error
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			err = parseError(c.httpStatusCode, c.r)
			assert.Error(t, err)
			cmdErr, ok := err.(*Error)
			assert.True(t, ok)
			assert.Equal(t, cmdErr.Code, c.e.Code, "error code")
			assert.Equal(t, cmdErr.Message, c.e.Message, "error message")
			assert.Equal(t, cmdErr.Data["x"], c.e.Data["x"])
			assert.Equal(t, cmdErr.Error(), c.e.Error())
		})
	}
}

package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
)

type Client interface {
	Delete(ctx context.Context, path string) (resp *Response, err error)
	Get(ctx context.Context, path string) (resp *Response, err error)
	Post(ctx context.Context, path string, params interface{}) (resp *Response, err error)
}

type params map[string]interface{}

type httpClient struct {
	client  httpclient.Client
	Headers http.Header
}

func NewClient(client httpclient.Client) Client {
	headers := make(http.Header)
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	headers.Add("Accept-charset", "utf-8")
	headers.Add("Cache-Control", "no-cache")
	return &httpClient{
		client:  client,
		Headers: headers,
	}
}

func (c *httpClient) handleResponse(r *http.Response) (resp *Response, err error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &resp); err != nil {
		return nil, err
	}
	if r.StatusCode >= 400 || resp.Status != 0 {
		return nil, parseError(r.StatusCode, resp)
	}
	resp.SessionID = strings.Trim(resp.SessionID, "{}\"")
	return resp, nil
}

func (c *httpClient) Delete(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Delete(ctx, path, c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

func (c *httpClient) Get(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Get(ctx, path, c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

func (c *httpClient) Post(ctx context.Context, path string, params interface{}) (resp *Response, err error) {
	if params == nil {
		params = map[string]interface{}{}
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.client.Post(ctx, path, bytes.NewReader(payload), c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

type StackFrame struct {
	FileName   string
	ClassName  string
	MethodName string
	LineNumber int
}

type Error struct {
	Code       string      `json:"error"`
	Message    string      `json:"message"`
	StackTrace string      `json:"stacktrace"`
	Data       interface{} `json:"data"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", e.Message, e.Code)
}

type Response struct {
	SessionID string          `json:"sessionId"`
	Status    int             `json:"status"`
	Value     json.RawMessage `json:"value"`
}

func parseError(respStatusCode int, resp *Response) error {
	var (
		sc  = httpStatusCode(respStatusCode)
		msg string
	)
	// if status exists
	if resp.Status > 0 {
		sm, ok := statusCode[resp.Status]
		if ok {
			sc = strconv.Itoa(resp.Status)
			msg = sm
		}
	}
	cmdErr := &Error{
		Code:    sc,
		Message: msg,
	}
	if err := json.Unmarshal(resp.Value, cmdErr); err != nil {
		cmdErr.Message = string(resp.Value)
	}
	if len(cmdErr.Message) == 0 {
		cmdErr.Message = cmdErr.Code
	}
	return cmdErr
}

func httpStatusCode(code int) (s string) {
	switch code {
	case 200, 400:
		s = StatusMissingCommandParameters
	case 404:
		s = StatusUnknownCommand
	case 405:
		s = StatusInvalidCommandMethod
	case 500:
		s = StatusFailedCommand
	case 501:
		s = StatusUnimplementedCommand
	}
	return s
}

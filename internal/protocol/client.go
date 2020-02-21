package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
)

type Client interface {
	Delete(ctx context.Context, path string) (resp *Response, err error)
	Get(ctx context.Context, path string) (resp *Response, err error)
	Post(ctx context.Context, path string, params interface{}) (resp *Response, err error)
}

type params map[string]interface{}

type httpClient struct {
	client httpclient.Client
}

func NewClient(client httpclient.Client) Client {
	return &httpClient{
		client: client,
	}
}

func (c *httpClient) defaultHeader() http.Header {
	header := make(http.Header)
	header.Add("Content-Type", "application/json;charset=utf-8")
	header.Add("Accept", "application/json")
	header.Add("Accept-charset", "utf-8")
	return header
}

func (c *httpClient) handleResponse(r *http.Response) (resp *Response, err error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	resp = &Response{}
	err = json.Unmarshal(buf, resp)
	if err != nil && r.StatusCode == 200 {
		return nil, errors.New("error: response must be a JSON object")
	}
	if r.StatusCode >= 400 || resp.Status != 0 {
		return nil, parseError(r.StatusCode, *resp)
	}
	resp.SessionID = bytes.Trim(resp.SessionID, "{}\"")
	return resp, nil
}

func (c *httpClient) Delete(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Delete(ctx, path, c.defaultHeader())
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

func (c *httpClient) Get(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Get(ctx, path, c.defaultHeader())
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
	httpResp, err := c.client.Post(ctx, path, bytes.NewReader(payload), c.defaultHeader())
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

type Response struct {
	SessionID []byte `json:"sessionId"`
	Status    int    `json:"status"`
	Value     []byte `json:"value"`
}

func parseError(code int, resp Response) error {
	var responseCodeError string
	switch code {
	// workaround: chromedriver could returns 200 code on errors
	case 200:
	case 400:
		responseCodeError = "400: Missing Command Parameters"
	case 404:
		responseCodeError = "404: Unknown command/Resource Not Found"
	case 405:
		responseCodeError = "405: Invalid Command Method"
	case 500:
		responseCodeError = "500: Failed Command"
	case 501:
		responseCodeError = "501: Unimplemented Command"
	default:
		responseCodeError = "Unknown error"
	}
	if resp.Status == 0 {
		return &CommandError{StatusCode: -1, ErrorType: responseCodeError}
	}
	commandError := &CommandError{StatusCode: resp.Status, ErrorType: responseCodeError}
	err := json.Unmarshal(resp.Value, commandError)
	if err != nil {
		// workaround: firefox could returns a string instead of a JSON object on errors
		commandError.Message = string(resp.Value)
	}
	return commandError
}

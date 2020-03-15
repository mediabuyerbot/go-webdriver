package w3c

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mediabuyerbot/httpclient"
)

type Doer interface {
	Do(ctx context.Context, method string, path string, p Params) (resp *Response, err error)
}

var (
	ErrResponseWithoutBody = errors.New("w3c: response without body")
	ErrEmptyResponse       = errors.New("w3c: empty response")
)

// Params represents the key-value pairs.
type Params map[string]interface{}

// Set sets the params entries associated with key to the
// single element value. It replaces any existing values
// associated with key.
func (p Params) Set(key string, value interface{}) {
	p[key] = value
}

func (p Params) IsEmpty() bool {
	return len(p) == 0
}

type transport struct {
	client  httpclient.Client
	Headers http.Header
}

func WithClient(client httpclient.Client) Doer {
	headers := make(http.Header)
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	headers.Add("Accept-charset", "utf-8")
	headers.Add("Cache-Control", "no-cache")
	return &transport{
		client:  client,
		Headers: headers,
	}
}

func (c *transport) handleResponse(r *http.Response) (resp *Response, err error) {
	if r.Body == nil {
		return nil, ErrResponseWithoutBody
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, ErrEmptyResponse
	}
	if err := json.Unmarshal(buf, &resp); err != nil {
		return nil, fmt.Errorf("w3c: %v, %s", err, string(buf))
	}
	if r.StatusCode >= 400 || resp.Status != 0 {
		return nil, parseError(r.StatusCode, resp)
	}

	resp.SessionID = strings.Trim(resp.SessionID, "{}\"")
	return resp, nil
}

func (c *transport) Do(ctx context.Context, method string, path string, p Params) (resp *Response, err error) {
	switch method {
	case http.MethodGet:
		resp, err = c.getRequest(ctx, path)
	case http.MethodDelete:
		resp, err = c.deleteRequest(ctx, path)
	case http.MethodPost:
		resp, err = c.postRequest(ctx, path, p)
	default:
		err = errors.New("protocol: unknown HTTP method")
	}
	return resp, err
}

func (c *transport) deleteRequest(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Delete(ctx, path, c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

func (c *transport) getRequest(ctx context.Context, path string) (resp *Response, err error) {
	httpResp, err := c.client.Get(ctx, path, c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

func (c *transport) postRequest(ctx context.Context, path string, p Params) (resp *Response, err error) {
	if p == nil {
		p = Params{}
	}
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.client.Post(ctx, path, bytes.NewReader(payload), c.Headers)
	if err != nil {
		return nil, err
	}
	return c.handleResponse(httpResp)
}

type Response struct {
	SessionID string          `json:"sessionId"`
	Status    int             `json:"status"`
	Value     json.RawMessage `json:"value"`
}

func (r Response) Success() bool {
	return string(r.Value) == "null"
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

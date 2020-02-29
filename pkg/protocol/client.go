package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	log.Println("RESP", string(buf))
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

package httpclient

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gojek/valkyrie"

	"github.com/pkg/errors"
)

type params map[string]interface{}

const (
	defaultRetryCount  = 0
	defaultHTTPTimeout = 30 * time.Second
)

// Doer interface has the method required to use a type as custom http client.
// The net/*http.Client type satisfies this interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client Is a generic HTTP client interface
type Client interface {
	Get(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Patch(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

// RequestHook allows a function to run before each retry. The HTTP
// request which will be made, and the retry number (0 for the initial
// request) are available to users.
type RequestHook func(*http.Request, int) *http.Request

// ResponseHook is like RequestHook, but allows running a function
// on each HTTP response. This function will be invoked at the end of
// every HTTP request executed, regardless of whether a subsequent retry
// needs to be performed or not. If the response body is read or closed
// from this method, this will affect the response returned from Do().
type ResponseHook func(*http.Request, *http.Response)

// CheckRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If CheckRetry returns false, the Client stops retrying
// and returns the response to the caller. If CheckRetry returns an error,
// that error value is returned in lieu of the error from the request. The
// Client will close any response body when retrying, but if the retry is
// aborted it is up to the CheckRetry callback to properly close any
// response body before returning.
type CheckRetry func(req *http.Request, resp *http.Response, err error) (bool, error)

// Backoff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type Backoff func(attemptNum int, resp *http.Response) time.Duration

// ErrorHandler is called if retries are expired, containing the last status
// from the http library. If not specified, default behavior for the library is
// to close the body and return an error indicating how many tries were
// attempted. If overriding this, be sure to close the body if needed.
type ErrorHandler func(resp *http.Response, err error, numTries int) (*http.Response, error)

type ErrorHook func(req *http.Request, err error, retry int)

// HttpClient is the http client implementation
type HttpClient struct {
	baseURL      *url.URL
	client       *http.Client
	retryCount   int
	requestHook  RequestHook
	responseHook ResponseHook
	errorHook    ErrorHook
	checkRetry   CheckRetry
	backoff      Backoff
	errorHandler ErrorHandler
}

// NewClient returns a new instance of http Client
func NewClient(baseURL string, opts ...Option) (Client, error) {
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	client := HttpClient{
		baseURL:    u,
		retryCount: defaultRetryCount,
		client: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client, nil
}

// Get makes a HTTP GET request to provided URL
func (c *HttpClient) Get(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *HttpClient) Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *HttpClient) Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *HttpClient) Patch(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *HttpClient) Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (c *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	req.Close = true

	var bodyReader *bytes.Reader
	var numTries int

	if req.Body != nil {
		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		req.Body = ioutil.NopCloser(bodyReader) // prevents closing the body between retries
	}

	multiErr := &valkyrie.MultiError{}
	for i := 0; i <= c.retryCount; i++ {
		numTries++
		if resp != nil {
			resp.Body.Close()
		}

		// request hook
		if c.requestHook != nil {
			req = c.requestHook(req, i)
		}

		resp, err = c.client.Do(req)
		if bodyReader != nil {
			// Reset the body reader after the request since at this point it's already read
			// Note that it's safe to ignore the error here since the 0,0 position is always valid
			_, _ = bodyReader.Seek(0, 0)
		}

		if err != nil && c.errorHook != nil {
			c.errorHook(req, err, i)
		}

		// response hook
		if err == nil && c.responseHook != nil {
			c.responseHook(req, resp)
		}

		// check retry
		if c.checkRetry != nil {
			checkOK, checkErr := c.checkRetry(req, resp, err)
			if !checkOK {
				if checkErr != nil {
					err = checkErr
				}
				c.client.CloseIdleConnections()
				return resp, err
			}
		}

		if c.backoff != nil {
			wait := c.backoff(i, resp)
			time.Sleep(wait)
			continue
		}

		if err == nil {
			break
		}
	}

	if c.errorHandler != nil {
		return c.errorHandler(resp, err, numTries+1)
	}

	if err != nil {
		return nil, err
	}

	return resp, multiErr.HasError()
}

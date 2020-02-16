package httpclient

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionsDefault(t *testing.T) {
	cli, err := NewClient("127.0.0.1")
	assert.NoError(t, err, "should not have failed to make a Client")

	httpcli, ok := cli.(*HttpClient)
	assert.Equal(t, true, ok)
	assert.Equal(t, httpcli.baseURL, "http://127.0.0.1")
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.backoff).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(defaultBackOffPolicy).Pointer()).Name(),
	)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.checkRetry).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(defaultCheckRetryPolicy).Pointer()).Name(),
	)
	assert.Equal(t, httpcli.retryCount, 0)
	assert.Nil(t, httpcli.errorHandler)
	assert.Nil(t, httpcli.requestHook)
	assert.Nil(t, httpcli.responseHook)
	assert.Nil(t, httpcli.errorHook)
}

func TestOptions(t *testing.T) {
	var (
		retry        = 5
		timeout      = time.Second
		errorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			return resp, nil
		}
		backOff = func(attemptNum int, resp *http.Response) time.Duration {
			return time.Second
		}
		reqHook = func(r *http.Request, n int) *http.Request {
			return r
		}
		respHook = func(*http.Request, *http.Response) {
			return
		}
		tr = &http.Transport{}
	)

	cli, err := NewClient("127.0.0.1",
		WithTimeout(timeout),
		WithErrorHandler(errorHandler),
		WithBackOff(backOff),
		WithRequestHook(reqHook),
		WithResponseHook(respHook),
		WithRetryCount(retry),
		WithTransport(tr),
	)
	assert.NoError(t, err, "should not have failed to make a Client")

	httpcli, ok := cli.(*HttpClient)
	assert.Equal(t, true, ok)

	assert.Equal(t, timeout, httpcli.client.Timeout)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.errorHandler).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(errorHandler).Pointer()).Name(),
	)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.backoff).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(backOff).Pointer()).Name(),
	)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.requestHook).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(reqHook).Pointer()).Name(),
	)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.responseHook).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(respHook).Pointer()).Name(),
	)
	assert.Equal(t, httpcli.retryCount, retry)
	assert.Equal(t, httpcli.client.Transport, tr)
}

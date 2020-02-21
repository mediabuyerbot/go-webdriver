package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

var defaultWithTimeoutPolicy = WithTimeout(10 * time.Millisecond)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		err     bool
		want    string
	}{
		{
			name:    "without http schema",
			baseURL: "google.com",
			err:     false,
			want:    "http://google.com",
		},
		{
			name:    "without http schema and bad url",
			baseURL: "googlecom",
			err:     false,
			want:    "http://googlecom",
		},
		{
			name:    "without http schema and ipv4 with port",
			baseURL: "127.0.0.1:9787",
			err:     false,
			want:    "http://127.0.0.1:9787",
		},
		{
			name:    "with https schema and ipv4 with port",
			baseURL: "https://127.0.0.1:9787",
			err:     false,
			want:    "https://127.0.0.1:9787",
		},
		{
			name:    "empty url",
			baseURL: "",
			err:     true,
			want:    "",
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			cli, err := NewClient(c.baseURL, WithRetryCount(2))
			if !c.err && err != nil {
				t.Fatal(err)
			}
			if c.err && err != nil {
				return
			}
			assert.Equal(t, cli.(*HttpClient).baseURL, c.want)
			assert.Equal(t, cli.(*HttpClient).retryCount, 2)
		})
	}
}

func TestHttpClient_DoSuccess(t *testing.T) {
	var payload = `{"want":"mediabuyerbot"}`
	cli, err := NewClient("127.0.0.1:8989", defaultWithTimeoutPolicy)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")
	response, err := cli.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, payload, string(body))
}

func TestHttpClient_DoFailWithReTray(t *testing.T) {
	var (
		payload   = `{"want":"mediabuyerbot"}`
		wantRetry = 3
		haveRetry = 0
	)

	cli, err := NewClient("127.0.0.1:8989",
		defaultWithTimeoutPolicy,
		WithRetryCount(2),
		WithBackOff(func(attemptNum int, resp *http.Response) time.Duration {
			return 3 * time.Millisecond
		}),
		WithRequestHook(func(request *http.Request, i int) *http.Request {
			haveRetry++
			return request
		}),
		WithErrorHandler(func(resp *http.Response, err error, numTries int) (response *http.Response, err2 error) {
			return resp, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")
	response, err := cli.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, payload, string(body))
	assert.Equal(t, wantRetry, haveRetry)
}

func TestHttpClient_DoRefusedConn(t *testing.T) {
	var (
		wantRetry = 3
		haveRetry = 0
	)

	cli, err := NewClient("127.0.0.1:8989",
		defaultWithTimeoutPolicy,
		WithRetryCount(2),
		WithBackOff(func(attemptNum int, resp *http.Response) time.Duration {
			return 3 * time.Millisecond
		}),
		WithRequestHook(func(request *http.Request, i int) *http.Request {
			haveRetry++
			return request
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://xyzqw.media.bot.nop", nil)
	require.NoError(t, err)
	response, err := cli.Do(req)
	require.Error(t, err)
	require.Nil(t, response)
	assert.Equal(t, wantRetry, haveRetry)
}

func TestHttpClient_Get(t *testing.T) {
	payload := `{ "response": "ok" }`
	path := "/path"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		assert.Equal(t, r.URL.Path, path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cli, err := NewClient(server.URL, defaultWithTimeoutPolicy)
	require.NoError(t, err, "should not have failed to make a Client")

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")
	resp, err := cli.Get(context.TODO(), path, headers)
	require.NoError(t, err, "should not have failed to make a GET request")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, resp))
}

func TestHttpClient_Delete(t *testing.T) {
	payload := `{ "response": "ok" }`
	path := "/path"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		assert.Equal(t, r.URL.Path, path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cli, err := NewClient(server.URL, defaultWithTimeoutPolicy)
	require.NoError(t, err, "should not have failed to make a Client")

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")
	resp, err := cli.Delete(context.TODO(), path, headers)
	require.NoError(t, err, "should not have failed to make a DELETE request")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, resp))
}

func TestHttpClient_Put(t *testing.T) {
	payload := `{ "response": "ok" }`
	path := "/path"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		assert.Equal(t, r.URL.Path, path)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")
		assert.Equal(t, payload, string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cli, err := NewClient(server.URL, defaultWithTimeoutPolicy)
	require.NoError(t, err, "should not have failed to make a Client")

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")
	resp, err := cli.Put(context.TODO(), path, bytes.NewBufferString(payload), headers)
	require.NoError(t, err, "should not have failed to make a PUT request")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, resp))
}

func TestHttpClient_Post(t *testing.T) {
	payload := `{ "response": "ok" }`
	path := "/path"
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))
		assert.Equal(t, r.URL.Path, path)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")
		assert.Equal(t, payload, string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cli, err := NewClient(server.URL, defaultWithTimeoutPolicy)
	require.NoError(t, err, "should not have failed to make a Client")

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")
	resp, err := cli.Post(context.TODO(), path, bytes.NewBufferString(payload), headers)
	require.NoError(t, err, "should not have failed to make a POST request")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, resp))
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer response.Body.Close()
	}
	respBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")
	return string(respBody)
}

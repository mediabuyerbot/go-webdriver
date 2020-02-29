package protocol

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpCli := httpclient.NewMockClient(ctrl)
	client := NewClient(httpCli).(*httpClient)
	assert.Equal(t, client.Headers.Get("Content-Type"), "application/json;charset=utf-8")
	assert.Equal(t, client.Headers.Get("Accept"), "application/json")
	assert.Equal(t, client.Headers.Get("Accept-charset"), "utf-8")
	assert.Equal(t, client.Headers.Get("Cache-Control"), "no-cache")
}

func TestHttpClient_Post(t *testing.T) {
	param := params{"id": "id"}
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.URL.Path, "/status")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json;charset=utf-8")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-charset"), "utf-8")
		assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		var p params
		err = json.Unmarshal(body, &p)
		assert.NoError(t, err)
		assert.Equal(t, p["id"], param["id"])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sessionId":"1234", "value": {}, "status": 0}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Post(context.TODO(), "/status", param)
	assert.Nil(t, err)
	assert.Equal(t, resp.SessionID, "1234")
	assert.Equal(t, resp.Status, 0)
}

func TestHttpClient_PostWithoutParam(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.URL.Path, "/status")
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		var p params
		err = json.Unmarshal(body, &p)
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sessionId":"1234", "value": null}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Post(context.TODO(), "/status", nil)
	assert.Nil(t, err)
	assert.Equal(t, resp.SessionID, "1234")
	assert.Equal(t, resp.Status, 0)
}

func TestHttpClient_HandleErrorWithoutStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"value": {
             "error": "unexpected alert open",
		     "message": "error",
		     "stacktrace": "stacktrace",
		     "data": {
			    "text": "Message from window.alert"
	    	 }
        }}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Post(context.TODO(), "/", nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

	cmdErr, ok := err.(*Error)
	assert.True(t, ok)

	assert.Equal(t, cmdErr.Code, "unexpected alert open")
	assert.Equal(t, cmdErr.Message, "error")
	assert.Equal(t, cmdErr.RawStacktrace, "stacktrace")
	assert.Equal(t, cmdErr.Data["text"], "Message from window.alert")
}

func TestHttpClient_HandleErrorBadResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(``))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Post(context.TODO(), "/", nil)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHttpClient_Delete(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, r.URL.Path, "/status")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json;charset=utf-8")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-charset"), "utf-8")
		assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sessionId":"", "value": {}, "status": 0}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Delete(context.TODO(), "/status")
	assert.Nil(t, err)
	assert.Equal(t, resp.SessionID, "")
	assert.Equal(t, resp.Status, 0)
}

func TestHttpClient_Get(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, r.URL.Path, "/status")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json;charset=utf-8")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-charset"), "utf-8")
		assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sessionId":"444", "value": {}}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	httpCli, err := httpclient.NewClient(srv.URL)
	assert.Nil(t, err)

	cli := NewClient(httpCli)
	resp, err := cli.Get(context.TODO(), "/status")
	assert.Nil(t, err)
	assert.Equal(t, resp.SessionID, "444")
	assert.Equal(t, resp.Status, 0)
}

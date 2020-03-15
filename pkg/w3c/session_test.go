package w3c

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mediabuyerbot/httpclient"
	"github.com/stretchr/testify/assert"
)

func makeBody(s string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(s)))
}

func newHttpClient(t *testing.T, ctrl *gomock.Controller) (httpclient.Client, *httpclient.MockDoer) {
	doer := httpclient.NewMockDoer(ctrl)
	cli, err := httpclient.New(
		httpclient.WithBaseURL("http://127.0.0.1:9515"),
		httpclient.WithDoer(doer),
	)
	assert.Nil(t, err)
	return cli, doer
}

func TestNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)

	client, transport := newHttpClient(t, ctrl)
	request := WithClient(client)
	sessOpts := NewMockBrowserOptions(ctrl)
	alwaysMatch := make([]Capabilities, 0)
	alwaysMatch = append(alwaysMatch, Capabilities{CapabilityBrowserName: "firefox"})
	sessOpts.EXPECT().FirstMatch().AnyTimes().Return(alwaysMatch)
	sessOpts.EXPECT().AlwaysMatch().AnyTimes().Return(Capabilities{
		CapabilityBrowserName: "firefox",
	})

	// returns success
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: 200,
		Body:       makeBody(newSessSuccessResponseFixture),
	}, nil).Do(func(req *http.Request) {
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.URL.String(), "http://127.0.0.1:9515/session")
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, body, []byte(`{"capabilities":{"alwaysMatch":{"browserName":"firefox"},"firstMatch":[{"browserName":"firefox"}]}}`))
	})

	sess, err := NewSession(request, sessOpts)
	assert.Nil(t, err)
	assert.Equal(t, sess.ID(), "4419604c-8c72-ea4c-8859-5b5de5098b2f")

	// returns error (HTTP method not allowed)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       makeBody("HTTP method not allowed"),
	}, nil)
	sess, err = NewSession(request, sessOpts)
	assert.Error(t, err)
	assert.Nil(t, sess)

	// returns error (HTTP bad request)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       makeBody(`{"value":{"error":"invalid argument","message":"invalid type: null, expected a sequence","stacktrace":""}}`),
	}, nil)
	sess, err = NewSession(request, sessOpts)
	assert.Error(t, err)
	assert.Nil(t, sess)
	cmdErr, ok := err.(*Error)
	assert.True(t, ok)
	assert.Equal(t, cmdErr.Code, "invalid argument")
	assert.Equal(t, cmdErr.Message, "invalid type: null, expected a sequence")

	// return error
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(nil, errors.New("some error"))
	sess, err = NewSession(request, sessOpts)
	assert.EqualError(t, err, "some error")
	assert.Nil(t, sess)

	// returns error (HTTP OK and empty session id)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"sessionId":""}}`),
	}, nil)
	sess, err = NewSession(request, sessOpts)
	assert.Equal(t, err, ErrUnknownSession)
	assert.Nil(t, sess)

	// returns error (HTTP OK bad JSON format)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"sessionId":"}}`),
	}, nil)
	sess, err = NewSession(request, sessOpts)
	assert.Error(t, err)
	assert.Nil(t, sess)
}

func TestSession_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	client, transport := newHttpClient(t, ctrl)
	request := WithClient(client)
	sessOpts := NewMockBrowserOptions(ctrl)
	alwaysMatch := make([]Capabilities, 0)
	alwaysMatch = append(alwaysMatch, Capabilities{CapabilityBrowserName: "firefox"})
	sessOpts.EXPECT().FirstMatch().AnyTimes().Return(alwaysMatch)
	sessOpts.EXPECT().AlwaysMatch().AnyTimes().Return(Capabilities{
		CapabilityBrowserName: "firefox",
	})
	ctx := context.Background()

	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"sessionId":"123"}}`),
	}, nil)
	sess, err := NewSession(request, sessOpts)
	assert.Nil(t, err)

	// returns success
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value": null}`),
	}, nil).Do(func(req *http.Request) {
		assert.Equal(t, http.MethodDelete, req.Method)
		assert.Equal(t, nil, req.Body)
		assert.Equal(t, "http://127.0.0.1:9515/session/123", req.URL.String())
	})
	err = sess.Delete(ctx)
	assert.Nil(t, err)

	// returns error (HTTP 404 invalid session id)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       makeBody(`{"value":{"error":"invalid session id","message":"Tried to run command without establishing a connection","stacktrace":""}}`),
	}, nil)
	err = sess.Delete(ctx)
	assert.Equal(t, "invalid session id", err.(*Error).Code)

	// returns error (HTTP 200 invalid response data)
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{}}`),
	}, nil)
	err = sess.Delete(ctx)
	assert.Equal(t, ErrInvalidResponse, err)
}

func TestSession_Capabilities(t *testing.T) {
	ctrl := gomock.NewController(t)

	client, transport := newHttpClient(t, ctrl)
	request := WithClient(client)
	sessOpts := NewMockBrowserOptions(ctrl)
	alwaysMatch := make([]Capabilities, 0)
	alwaysMatch = append(alwaysMatch, Capabilities{CapabilityBrowserName: "firefox"})
	sessOpts.EXPECT().FirstMatch().AnyTimes().Return(alwaysMatch)
	sessOpts.EXPECT().AlwaysMatch().AnyTimes().Return(Capabilities{
		CapabilityBrowserName: "firefox",
	})

	// returns success
	transport.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(newSessSuccessCapabilitiesResponseFixture),
	}, nil)
	sess, err := NewSession(request, sessOpts)
	assert.Nil(t, err)

	cap := sess.Capabilities()
	assert.True(t, GetAcceptInsecureCerts(cap))
	assert.Equal(t, "firefox", GetBrowserName(cap))
	assert.Equal(t, "73.0.1", GetBrowserVersion(cap))
	assert.True(t, cap.GetBool("moz:accessibilityChecks"))
	assert.Equal(t, 52708, cap.GetInt("moz:processID"))

	timeouts := GetTimeout(cap)
	assert.NotNil(t, timeouts)
	assert.Equal(t, uint(0), timeouts.Implicit)
	assert.Equal(t, uint(300000), timeouts.PageLoad)
	assert.Equal(t, uint(30000), timeouts.Script)
}

func TestSession_Status(t *testing.T) {
	ctrl := gomock.NewController(t)

	client, transport := newHttpClient(t, ctrl)
	request := WithClient(client)
	sessOpts := NewMockBrowserOptions(ctrl)
	alwaysMatch := make([]Capabilities, 0)
	alwaysMatch = append(alwaysMatch, Capabilities{CapabilityBrowserName: "firefox"})
	sessOpts.EXPECT().FirstMatch().AnyTimes().Return(alwaysMatch)
	sessOpts.EXPECT().AlwaysMatch().AnyTimes().Return(Capabilities{
		CapabilityBrowserName: "firefox",
	})
	ctx := context.Background()

	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"sessionId":"123"}}`),
	}, nil)
	sess, err := NewSession(request, sessOpts)
	assert.Nil(t, err)

	// returns success (w3c only)
	transport.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"message":"Session already started","ready":false}}`),
	}, nil)
	st, err := sess.Status(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "Session already started", st.Message)
	assert.Equal(t, false, st.Ready)

	// returns success (w3c + extends)
	transport.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`{"value":{"build":{"version":"79.0.3945.36 (3582db32b33893869b8c1339e8f4d9ed1816f143-refs/branch-heads/3945@{#614})"},"message":"ChromeDriver ready for new sessions.","os":{"arch":"x86_64","name":"Mac OS X","version":"10.15.2"},"ready":true}}`),
	}, nil)
	st, err = sess.Status(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "ChromeDriver ready for new sessions.", st.Message)
	assert.Equal(t, true, st.Ready)
	assert.Equal(t, "79.0.3945.36 (3582db32b33893869b8c1339e8f4d9ed1816f143-refs/branch-heads/3945@{#614})", st.Build.Version)
	assert.Equal(t, "x86_64", st.OS.Arch)
	assert.Equal(t, "Mac OS X", st.OS.Name)
	assert.Equal(t, "10.15.2", st.OS.Version)

	// returns error
	transport.EXPECT().Do(gomock.Any()).Return(nil, errors.New("some errors"))
	st, err = sess.Status(ctx)
	assert.Error(t, err)

	// returns error (invalid response)
	transport.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       makeBody(`[`),
	}, nil)
	st, err = sess.Status(ctx)
	assert.Error(t, err)
}

// fixtures
var (
	newSessSuccessResponseFixture             = `{"value":{"sessionId":"4419604c-8c72-ea4c-8859-5b5de5098b2f","capabilities":{"acceptInsecureCerts":false,"browserName":"firefox","browserVersion":"73.0.1","moz:accessibilityChecks":false,"moz:buildID":"20200217142647","moz:geckodriverVersion":"0.26.0","moz:headless":false,"moz:processID":45423,"moz:profile":"/var/folders/kw/zznxy_m949v5yjwx4qp4wg540000gn/T/rust_mozprofileToXug6","moz:shutdownTimeout":60000,"moz:useNonSpecCompliantPointerOrigin":false,"moz:webdriverClick":true,"pageLoadStrategy":"normal","platformName":"mac","platformVersion":"19.2.0","rotatable":false,"setWindowRect":true,"strictFileInteractability":false,"timeouts":{"implicit":0,"pageLoad":300000,"script":30000},"unhandledPromptBehavior":"dismiss and notify"}}}`
	newSessSuccessCapabilitiesResponseFixture = `{"value":{"sessionId":"fe9a1972-cbea-d14d-838d-020e3b152611","capabilities":{"acceptInsecureCerts":true,"browserName":"firefox","browserVersion":"73.0.1","moz:accessibilityChecks":true,"moz:buildID":"20200217142647","moz:geckodriverVersion":"0.26.0","moz:headless":true,"moz:processID":52708,"moz:profile":"/var/folders/kw/zznxy_m949v5yjwx4qp4wg540000gn/T/rust_mozprofile7I4zZS","moz:shutdownTimeout":60000,"moz:useNonSpecCompliantPointerOrigin":true,"moz:webdriverClick":true,"pageLoadStrategy":"normal","platformName":"mac","platformVersion":"19.2.0","rotatable":true,"setWindowRect":true,"strictFileInteractability":true,"timeouts":{"implicit":0,"pageLoad":300000,"script":30000},"unhandledPromptBehavior":"dismiss and notify"}}}`
)

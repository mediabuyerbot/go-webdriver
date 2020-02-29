package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func testNewSession(cli Client, t *testing.T) Session {
	mockCli, ok := cli.(*MockClient)
	assert.True(t, ok)
	mockCli.EXPECT().Post(context.TODO(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "",
			Status:    0,
			Value:     json.RawMessage(`{"sessionId":"123", "capabilities": {"browserName": "chrome"}}`),
		}, nil).Do(func(_ context.Context, p string, a map[string]interface{}) {
		assert.Equal(t, p, "/session")
		assert.Equal(t, a["desiredCapabilities"].(map[string]string)["platform"], "linux")
	})
	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	sess, err := NewSession(cli, desired, required)
	assert.Nil(t, err)
	return sess
}

func TestNewSession_WithoutW3CCompatibility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "123",
			Status:    0,
			Value:     json.RawMessage(`{"browserName": "chrome"}`),
		}, nil)

	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	sess, err := NewSession(cli, desired, required)

	assert.Nil(t, err)
	assert.Equal(t, sess.ID(), "123")
	assert.Equal(t, sess.Capabilities().BrowserName(), "chrome")
}

func TestNewSession_WitW3CCompatibility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "",
			Status:    0,
			Value:     json.RawMessage(`{"sessionId":"123", "capabilities": {"browserName": "chrome"}}`),
		}, nil)

	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	sess, err := NewSession(cli, desired, required)

	assert.Nil(t, err)
	assert.Equal(t, sess.ID(), "123")
	assert.Equal(t, sess.Capabilities().BrowserName(), "chrome")
}

func TestNewSession_ErrorWitW3CCompatibility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "",
			Status:    0,
			Value:     json.RawMessage(`{"sessionId":"", "capabilities": {"browserName": "chrome"}}`),
		}, nil)

	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	_, err := NewSession(cli, desired, required)
	assert.Error(t, err)
}

func TestNewSession_ErrorWitW3CCompatibility2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "",
			Status:    0,
			Value:     json.RawMessage(`{"capabilities": {"browserName": "chrome"}}`),
		}, nil)

	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	_, err := NewSession(cli, desired, required)
	assert.Error(t, err)
}

func TestNewSession_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	errCode := strconv.Itoa(SessionNotCreatedExceptionStatusCode)

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(nil, &Error{
			Code: errCode,
		})

	sess, err := NewSession(cli, nil, nil)
	assert.Nil(t, sess)
	assert.Error(t, err)

	cmdErr, ok := err.(*Error)
	assert.True(t, ok)
	assert.Equal(t, cmdErr.Code, errCode)
}

func TestNewSession_WithInvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "",
			Status:    0,
			Value:     json.RawMessage(`unknown command: response error`),
		}, nil)

	sess, err := NewSession(cli, nil, nil)
	assert.Nil(t, sess)
	assert.Error(t, err)
}

func TestSession_Capabilities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	sess := testNewSession(cli, t)
	assert.Equal(t, sess.Capabilities().BrowserName(), "chrome")
	assert.Equal(t, sess.ID(), "123")
}

func TestSession_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	ctx := context.Background()

	cli.EXPECT().Delete(context.TODO(), "/session/123").Times(1).Return(&Response{
		SessionID: "",
		Status:    0,
		Value:     json.RawMessage(`{"sessionId":"123", "status":0, "value":null}`),
	}, nil)

	sess := testNewSession(cli, t)
	err := sess.Delete(ctx)
	assert.Nil(t, err)

	cli.EXPECT().Delete(context.TODO(), "/session/123").Times(1).Return(&Response{
		SessionID: "",
		Status:    0,
		Value:     json.RawMessage(`{"sessionId":"123", "status":0, "value":null}`),
	}, errors.New("some error"))
	err = sess.Delete(ctx)
	assert.Error(t, err)

	cli.EXPECT().Delete(context.TODO(), "/session/123").Times(1).Return(&Response{
		SessionID: "",
		Status:    0,
		Value:     json.RawMessage(`{"sessionId":"123, "status":0, "value":null}`),
	}, errors.New("some error"))
	err = sess.Delete(ctx)
	assert.Error(t, err)
}

func TestSession_Status(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	ctx := context.Background()

	// chrome
	cli.EXPECT().Get(ctx, "/status").Times(1).Return(&Response{
		Value: []byte(`{"build":{"version":"79.0.3945.36 (3582db32b33893869b8c1339e8f4d9ed1816f143-refs/branch-heads/3945@{#614})"},"message":"ChromeDriver ready for new sessions.","os":{"arch":"x86_64","name":"Mac OS X","version":"1"},"ready":true}`)}, nil)

	sess := testNewSession(cli, t)

	st, err := sess.Status(ctx)
	assert.Nil(t, err)
	assert.True(t, st.HasExtensionInfo())
	assert.Equal(t, st.Build.Version, "79.0.3945.36 (3582db32b33893869b8c1339e8f4d9ed1816f143-refs/branch-heads/3945@{#614})")
	assert.True(t, st.Ready)

	// firefox
	cli.EXPECT().Get(ctx, "/status").Times(1).Return(&Response{
		Value: []byte(`{"message":"","ready":true}`)}, nil)

	st, err = sess.Status(ctx)
	assert.Nil(t, err)
	assert.False(t, st.HasExtensionInfo())
	assert.True(t, st.Ready)

	cli.EXPECT().Get(ctx, "/status").Times(1).Return(&Response{
		Value: []byte(`{"message":"","ready":true}`)}, errors.New("some error"))
	st, err = sess.Status(ctx)
	assert.NotNil(t, err)

	cli.EXPECT().Get(ctx, "/status").Times(1).Return(&Response{
		Value: []byte(`{"message":","ready":true}`)}, nil)
	st, err = sess.Status(ctx)
	assert.NotNil(t, err)

}

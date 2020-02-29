package protocol

import (
	"context"
	"encoding/json"
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
	assert.Equal(t, sess.ID(), SessionID("123"))
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
	assert.Equal(t, sess.ID(), SessionID("123"))
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
	assert.Equal(t, sess.ID(), SessionID("123"))
}

func TestSession_DeleteSuccess(t *testing.T) {
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
}

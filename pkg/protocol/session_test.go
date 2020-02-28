package protocol

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestNewSessionWithoutW3CCompatibility(t *testing.T) {
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

func TestNewSessionWitW3CCompatibility(t *testing.T) {
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

func TestNewSessionErrorWitW3CCompatibility(t *testing.T) {
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

func TestNewSessionErrorWitW3CCompatibility2(t *testing.T) {
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

func TestNewSessionWithError(t *testing.T) {
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

func TestNewSessionWithInvalidResponse(t *testing.T) {
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

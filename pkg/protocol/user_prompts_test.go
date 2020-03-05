package protocol

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var alertErr = &Error{Code: "error"}

func newAlert(t *testing.T, sessID string) (Alert, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewAlert(cli, sessID)
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestAlert_Accept(t *testing.T) {
	alert, cli, done := newAlert(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/accept", nil).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := alert.Accept(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/accept", nil).Return(nil, alertErr)
	err = alert.Accept(ctx)
	assert.Equal(t, err, alertErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/accept", nil).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = alert.Accept(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestAlert_Dismiss(t *testing.T) {
	alert, cli, done := newAlert(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/dismiss", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := alert.Dismiss(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/dismiss", gomock.Any()).Times(1).Return(nil, alertErr)
	err = alert.Dismiss(ctx)
	assert.Equal(t, err, alertErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/dismiss", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = alert.Dismiss(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestAlert_Text(t *testing.T) {
	alert, cli, done := newAlert(t, "123")
	defer done()

	ctx := context.TODO()
	wantMessage := "some message"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/alert/text", nil).Return(
		&Response{
			Value: []byte(wantMessage),
		}, nil)
	haveMessage, err := alert.Text(ctx)
	assert.Nil(t, err)
	assert.Equal(t, wantMessage, haveMessage)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/alert/text", nil).Return(nil, alertErr)
	haveMessage, err = alert.Text(ctx)
	assert.Equal(t, err, alertErr)
	assert.Empty(t, haveMessage)
}

func TestAlert_SetText(t *testing.T) {
	alert, cli, done := newAlert(t, "123")
	defer done()

	ctx := context.TODO()
	wantText := "some text"
	param := Params{"text": wantText}

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/text", param).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method, path string, p Params) {
		assert.Equal(t, p["text"], wantText)
	})
	err := alert.SetText(ctx, wantText)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/text", param).Return(nil, alertErr)
	err = alert.SetText(ctx, wantText)
	assert.Equal(t, err, alertErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/alert/text", param).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = alert.SetText(ctx, wantText)
	assert.Equal(t, err, ErrInvalidResponse)
}

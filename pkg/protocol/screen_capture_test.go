package protocol

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var screenCaptureErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newScreenCapture(t *testing.T, sessID string) (ScreenCapture, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewScreenCapture(cli, sessID)
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestScreenCapture_Take(t *testing.T) {
	sc, cli, done := newScreenCapture(t, "123")
	defer done()

	ctx := context.TODO()
	want := []byte{'m', 'm', 'a', 'd', 'f', 'o', 'x'}

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/screenshot", nil).Times(1).Return(
		&Response{
			Value: []byte(base64.StdEncoding.EncodeToString(want)),
		}, nil)
	have, err := sc.Take(ctx)
	assert.Nil(t, err)
	assert.Equal(t, want, have)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/screenshot", nil).Times(1).Return(nil, screenCaptureErr)
	have, err = sc.Take(ctx)
	assert.Equal(t, err, screenCaptureErr)
	assert.Empty(t, have)
}

func TestScreenCapture_TakeElement(t *testing.T) {
	sc, cli, done := newScreenCapture(t, "123")
	defer done()

	ctx := context.TODO()
	want := []byte{'m', 'm', 'a', 'd', 'f', 'o', 'x'}

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/123/screenshot", nil).Times(1).Return(
		&Response{
			Value: []byte(base64.StdEncoding.EncodeToString(want)),
		}, nil)
	have, err := sc.TakeElement(ctx, "123")
	assert.Nil(t, err)
	assert.Equal(t, want, have)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/123/screenshot", nil).Times(1).Return(nil, screenCaptureErr)
	have, err = sc.TakeElement(ctx, "123")
	assert.Equal(t, err, screenCaptureErr)
	assert.Empty(t, have)
}

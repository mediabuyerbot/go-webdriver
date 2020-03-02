package protocol

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var timeoutsErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newTimeouts(t *testing.T, sessID string) (Timeouts, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewTimeouts(cli, "123")
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestTimeouts_GetTimeouts(t *testing.T) {
	timeouts, cli, done := newTimeouts(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 10000, "script": 10000, "Implicit": 10000}`),
		}, nil)

	ti, err := timeouts.GetTimeouts(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ti.Implicit, DefaultTimeoutMs)
	assert.Equal(t, ti.Script, DefaultTimeoutMs)
	assert.Equal(t, ti.PageLoad, DefaultTimeoutMs)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 10000, "script": 10000, "Implicit": 10000}`),
		}, timeoutsErr)
	_, err = timeouts.GetTimeouts(ctx)
	assert.Error(t, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{`),
		}, nil)
	_, err = timeouts.GetTimeouts(ctx)
	assert.Error(t, err)
}

func TestTimeouts_SetTimeouts(t *testing.T) {
	timeouts, cli, done := newTimeouts(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(PageLoadTimeout)], DefaultTimeoutMs)
	})
	err := timeouts.SetTimeouts(ctx, PageLoadTimeout, DefaultTimeoutMs)
	assert.Nil(t, err)

	// returns error (validation)
	err = timeouts.SetTimeouts(ctx, Timeout("unknown"), DefaultTimeoutMs)
	assert.Equal(t, err, ErrTimeoutConfiguration)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, timeoutsErr)

	err = timeouts.SetTimeouts(ctx, PageLoadTimeout, DefaultTimeoutMs)
	assert.Error(t, err)

	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(ImplicitTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetImplicitTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)

	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(PageLoadTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetPageLoadTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)

	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(ScriptTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetScriptTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)
}

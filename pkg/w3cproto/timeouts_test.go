package w3cproto

import (
	"context"
	"net/http"
	"testing"
	"time"

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
	cx := NewTimeouts(cli, sessID)
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestTimeouts_Timeout(t *testing.T) {
	t1 := Timeout{
		Script:   1000,
		PageLoad: 1000,
		Implicit: 1000,
	}
	onesec := float64(1)
	assert.Equal(t, t1.GetScript().Seconds(), onesec)
	assert.Equal(t, t1.GetImplicit().Seconds(), onesec)
	assert.Equal(t, t1.GetPageLoad().Seconds(), onesec)
}

func TestTimeouts_GetTimeouts(t *testing.T) {
	timeouts, cli, done := newTimeouts(t, "123")
	defer done()

	ctx := context.TODO()
	wantTimeout := time.Second

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 1000, "script": 1000, "Implicit": 1000}`),
		}, nil)

	ti, err := timeouts.Get(ctx)

	assert.Nil(t, err)
	assert.Equal(t, ti.GetImplicit().Milliseconds(), wantTimeout.Milliseconds())
	assert.Equal(t, ti.GetScript().Milliseconds(), wantTimeout.Milliseconds())
	assert.Equal(t, ti.GetPageLoad().Milliseconds(), wantTimeout.Milliseconds())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 1000, "script": 1000, "Implicit": 1000}`),
		}, timeoutsErr)
	_, err = timeouts.Get(ctx)
	assert.Error(t, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/timeouts", nil).Times(1).Return(
		&Response{
			Value: []byte(`{`),
		}, nil)
	_, err = timeouts.Get(ctx)
	assert.Error(t, err)
}

func TestTimeouts_SetTimeouts(t *testing.T) {
	timeouts, cli, done := newTimeouts(t, "123")
	defer done()

	ctx := context.TODO()
	wantTimeout := time.Second

	// returns success (pageLoad)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[pageLoadTimeout], wantTimeout.Milliseconds())
	})
	err := timeouts.SetPageLoad(ctx, wantTimeout)
	assert.Nil(t, err)

	// returns success (script)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[implicitTimeout], wantTimeout.Milliseconds())
	})
	err = timeouts.SetImplicit(ctx, wantTimeout)
	assert.Nil(t, err)

	// returns success (implicit)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params[scriptTimeout], wantTimeout.Milliseconds())
	})
	err = timeouts.SetScript(ctx, wantTimeout)
	assert.Nil(t, err)
}

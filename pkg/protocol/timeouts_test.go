package protocol

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestTimeouts_GetTimeouts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)
	timeouts := NewTimeouts(cli, "123")

	cli.EXPECT().Get(ctx, "/session/123/timeouts").Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 10000, "script": 10000, "Implicit": 10000}`),
		}, nil)

	ti, err := timeouts.GetTimeouts(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ti.Implicit, DefaultTimeoutMs)
	assert.Equal(t, ti.Script, DefaultTimeoutMs)
	assert.Equal(t, ti.PageLoad, DefaultTimeoutMs)

	cli.EXPECT().Get(ctx, "/session/123/timeouts").Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad": 10000, "script": 10000, "Implicit": 10000}`),
		}, errors.New("some error"))
	_, err = timeouts.GetTimeouts(ctx)
	assert.Error(t, err)

	cli.EXPECT().Get(ctx, "/session/123/timeouts").Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad: 10000, "script": 10000, "Implicit": 10000}`),
		}, nil)
	_, err = timeouts.GetTimeouts(ctx)
	assert.Error(t, err)
}

func TestTimeouts_SetTimeouts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)
	timeouts := NewTimeouts(cli, "123")

	cli.EXPECT().Post(ctx, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad: 10000, "script": 10000, "Implicit": 10000}`),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(PageLoadTimeout)], DefaultTimeoutMs)
	})

	err := timeouts.SetTimeouts(ctx, PageLoadTimeout, DefaultTimeoutMs)
	assert.Nil(t, err)

	err = timeouts.SetTimeouts(ctx, Timeout("unknown"), DefaultTimeoutMs)
	assert.Equal(t, err, ErrTimeoutConfiguration)

	cli.EXPECT().Post(ctx, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, errors.New("some error"))

	err = timeouts.SetTimeouts(ctx, PageLoadTimeout, DefaultTimeoutMs)
	assert.Error(t, err)

	cli.EXPECT().Post(ctx, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad: 10000, "script": 10000, "Implicit": 10000}`),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(ImplicitTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetImplicitTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)

	cli.EXPECT().Post(ctx, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad: 10000, "script": 10000, "Implicit": 10000}`),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(PageLoadTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetPageLoadTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)

	cli.EXPECT().Post(ctx, "/session/123/timeouts", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{"pageLoad: 10000, "script": 10000, "Implicit": 10000}`),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Equal(t, params[string(ScriptTimeout)], DefaultTimeoutMs)
	})

	err = timeouts.SetScriptTimeout(ctx, DefaultTimeoutMs)
	assert.Nil(t, err)
}

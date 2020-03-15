package w3c

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var navigationErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newNavigation(t *testing.T, sessID string) (Navigation, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	n := NewNavigation(cli, "123")
	return n, cli, func() {
		ctrl.Finish()
	}
}

func TestNavigation_NavigateTo(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botURL := "https://mediabuyerbot.com"
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Equal(t, params["url"], botURL)
	})
	err := navigation.NavigateTo(ctx, botURL)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(`null`),
		}, navigationErr)
	err = navigation.NavigateTo(ctx, botURL)
	assert.Equal(t, err, navigationErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(`{}`),
		}, nil)
	err = navigation.NavigateTo(ctx, botURL)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_GetCurrentURL(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botURL := "https://mediabuyerbot.com"
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/url", nil).Times(1).Return(
		&Response{
			SessionID: "132",
			Status:    0,
			Value:     []byte(botURL),
		}, nil)
	curURL, err := navigation.GetCurrentURL(ctx)
	assert.Nil(t, err)
	assert.Equal(t, curURL, botURL)

	// returns errors (invalid response)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/url", nil).Times(1).Return(
		nil, ErrInvalidResponse)
	curURL, err = navigation.GetCurrentURL(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
	assert.Empty(t, curURL)
}

func TestNavigation_Back(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/back", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Back(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/back", gomock.Any()).Times(1).Return(nil, navigationErr)
	err = navigation.Back(ctx)
	assert.Equal(t, err, navigationErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/back", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Back(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_Forward(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/forward", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Forward(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/forward", gomock.Any()).Times(1).Return(nil, navigationErr)
	err = navigation.Forward(ctx)
	assert.Equal(t, err, navigationErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/forward", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Forward(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_Refresh(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/refresh", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Refresh(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/refresh", gomock.Any()).Times(1).Return(nil, navigationErr)
	err = navigation.Refresh(ctx)
	assert.Equal(t, err, navigationErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/refresh", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, method string, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Refresh(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_GetTitle(t *testing.T) {
	navigation, cli, done := newNavigation(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	wantTitle := "MediaBuyerBot"
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/title", nil).Times(1).Return(
		&Response{Value: []byte(wantTitle)}, nil)
	haveTitle, err := navigation.GetTitle(ctx)
	assert.Nil(t, err)
	assert.Equal(t, haveTitle, wantTitle)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/title", nil).Times(1).Return(nil, ErrInvalidResponse)
	haveTitle, err = navigation.GetTitle(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
	assert.Empty(t, haveTitle)
}

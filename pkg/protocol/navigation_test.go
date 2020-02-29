package protocol

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNavigation_NavigateTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	navigation := NewNavigation(cli, "123")
	botURL := "https://mediabuyerbot.com"
	someErr := errors.New("some error")
	cli.EXPECT().Post(ctx, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(`null`),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Equal(t, params["url"], botURL)
	})

	err := navigation.NavigateTo(ctx, botURL)
	assert.Nil(t, err)

	cli.EXPECT().Post(ctx, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(`null`),
		}, someErr)
	err = navigation.NavigateTo(ctx, botURL)
	assert.Equal(t, err, someErr)

	cli.EXPECT().Post(ctx, "/session/123/url", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Value:     []byte(``),
		}, nil)
	err = navigation.NavigateTo(ctx, botURL)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_GetCurrentURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	navigation := NewNavigation(cli, "123")
	botURL := "https://mediabuyerbot.com"
	cli.EXPECT().Get(ctx, "/session/123/url").Times(1).Return(
		&Response{
			SessionID: "132",
			Status:    0,
			Value:     []byte(botURL),
		}, nil)
	curURL, err := navigation.GetCurrentURL(ctx)
	assert.Nil(t, err)
	assert.Equal(t, curURL, botURL)

	cli.EXPECT().Get(ctx, "/session/123/url").Times(1).Return(
		nil, ErrInvalidResponse)
	curURL, err = navigation.GetCurrentURL(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
	assert.Empty(t, curURL)
}

func TestNavigation_Back(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	navigation := NewNavigation(cli, "123")

	cli.EXPECT().Post(ctx, "/session/123/back", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Back(ctx)
	assert.Nil(t, err)

	someErr := errors.New("some error")
	cli.EXPECT().Post(ctx, "/session/123/back", gomock.Any()).Times(1).Return(nil, someErr)
	err = navigation.Back(ctx)
	assert.Equal(t, err, someErr)

	cli.EXPECT().Post(ctx, "/session/123/back", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Back(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_Forward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	navigation := NewNavigation(cli, "123")

	cli.EXPECT().Post(ctx, "/session/123/forward", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Forward(ctx)
	assert.Nil(t, err)

	someErr := errors.New("some error")
	cli.EXPECT().Post(ctx, "/session/123/forward", gomock.Any()).Times(1).Return(nil, someErr)
	err = navigation.Forward(ctx)
	assert.Equal(t, err, someErr)

	cli.EXPECT().Post(ctx, "/session/123/forward", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Forward(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	navigation := NewNavigation(cli, "123")

	cli.EXPECT().Post(ctx, "/session/123/refresh", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte("null"),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err := navigation.Refresh(ctx)
	assert.Nil(t, err)

	someErr := errors.New("some error")
	cli.EXPECT().Post(ctx, "/session/123/refresh", gomock.Any()).Times(1).Return(nil, someErr)
	err = navigation.Refresh(ctx)
	assert.Equal(t, err, someErr)

	cli.EXPECT().Post(ctx, "/session/123/refresh", gomock.Any()).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(""),
		}, nil).Do(func(_ context.Context, p string, params map[string]interface{}) {
		assert.Nil(t, params)
	})
	err = navigation.Refresh(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestNavigation_GetTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cli := NewMockClient(ctrl)

	botTitle := "MediaBuyerBot"
	navigation := NewNavigation(cli, "123")
	cli.EXPECT().Get(ctx, "/session/123/title").Times(1).Return(
		&Response{Value: []byte(botTitle)}, nil)
	title, err := navigation.GetTitle(ctx)
	assert.Nil(t, err)
	assert.Equal(t, title, botTitle)

	cli.EXPECT().Get(ctx, "/session/123/title").Times(1).Return(nil, ErrInvalidResponse)
	title, err = navigation.GetTitle(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
	assert.Empty(t, title)
}

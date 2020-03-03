package protocol

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var cookieErr = &Error{
	Code:    "cookie code",
	Message: "cookie msg",
}

func newCookie(t *testing.T, sessID string) (Cookies, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cookie := NewCookies(cli, sessID)
	return cookie, cli, func() {
		ctrl.Finish()
	}
}

func TestCookies_AddCookie(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	wantCookie := MakeCookie()
	wantCookie.
		SetName("cookie").
		SetValue("cookie").
		SetPath("/path/path").
		SetDomain("mediabuyerbot.com")

	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		pc, ok := p["cookie"].(Params)
		assert.True(t, ok)
		assert.NotEmpty(t, pc)
		assert.Equal(t, pc["name"], wantCookie.Name())
		assert.Equal(t, pc["value"], wantCookie.Value())
		assert.Equal(t, pc["path"], wantCookie.Path())
		assert.Equal(t, pc["domain"], wantCookie.Domain())
	})
	err := cookies.Add(ctx, wantCookie)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(nil, cookieErr)
	err = cookies.Add(ctx, wantCookie)
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.Add(ctx, wantCookie)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestCookies_DeleteCookie(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := cookies.Delete(ctx, "test")
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie/test", nil).Times(1).Return(nil, cookieErr)
	err = cookies.Delete(ctx, "test")
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.Delete(ctx, "test")
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestCookies_DeleteAllCookies(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := cookies.DeleteAll(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie", nil).Times(1).Return(nil, cookieErr)
	err = cookies.DeleteAll(ctx)
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.DeleteAll(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestCookies_GetAllCookies(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()
	wantCookies := []Cookie{
		Cookie{
			CookieNameKey:  "test",
			CookieValueKey: "test",
		},
		Cookie{
			CookieNameKey:  "test1",
			CookieValueKey: "test1",
		},
	}

	// returns success
	value, err := json.Marshal(wantCookies)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: value,
		}, nil)
	haveCookies, err := cookies.All(ctx)
	assert.Nil(t, err)
	assert.ElementsMatch(t, haveCookies, wantCookies)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(nil, cookieErr)
	haveCookies, err = cookies.All(ctx)
	assert.Len(t, haveCookies, 0)
	assert.Equal(t, cookieErr, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: []byte(`[{":123}]`),
		}, nil)
	haveCookies, err = cookies.All(ctx)
	assert.Error(t, err)
	assert.Len(t, haveCookies, 0)
}

func TestCookies_GetNamedCookie(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()

	wantCookie := MakeCookie().SetName("test").SetValue("test")

	// returns success
	value, err := json.Marshal(wantCookie)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: value,
		}, nil)
	haveCookie, err := cookies.Get(ctx, "test")
	assert.Nil(t, err)
	assert.Equal(t, wantCookie.Name(), haveCookie.Name())
	assert.Equal(t, wantCookie.Value(), haveCookie.Value())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(nil, cookieErr)
	haveCookie, err = cookies.Get(ctx, "test")
	assert.Equal(t, cookieErr, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: []byte(`[`),
		}, nil)
	haveCookie, err = cookies.Get(ctx, "test")
	assert.Error(t, err)
	assert.Empty(t, haveCookie.Value())
}

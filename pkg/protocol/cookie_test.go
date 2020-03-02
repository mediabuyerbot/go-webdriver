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
	wantCookie := Cookie{
		Name:   "cookie",
		Value:  "cookie",
		Path:   "/path/path",
		Domain: "mediabuyerbot.com",
	}
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, p["name"], wantCookie.Name)
		assert.Equal(t, p["value"], wantCookie.Value)
		assert.Equal(t, p["path"], wantCookie.Path)
		assert.Equal(t, p["domain"], wantCookie.Domain)
		assert.Equal(t, p["secure"], wantCookie.Secure)
		assert.Equal(t, p["expiry"], wantCookie.Expiry)
		assert.Equal(t, p["httpOnly"], wantCookie.HttpOnly)
		assert.Equal(t, p["sameSite"], wantCookie.SameSite)
	})
	err := cookies.AddCookie(ctx, wantCookie)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(nil, cookieErr)
	err = cookies.AddCookie(ctx, wantCookie)
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/cookie", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.AddCookie(ctx, wantCookie)
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
	err := cookies.DeleteCookie(ctx, "test")
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie/test", nil).Times(1).Return(nil, cookieErr)
	err = cookies.DeleteCookie(ctx, "test")
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.DeleteCookie(ctx, "test")
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
	err := cookies.DeleteAllCookies(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie", nil).Times(1).Return(nil, cookieErr)
	err = cookies.DeleteAllCookies(ctx)
	assert.Equal(t, err, cookieErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cookies.DeleteAllCookies(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestCookies_GetAllCookies(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()
	wantCookies := []Cookie{
		Cookie{
			Name:  "test",
			Value: "test",
		},
		Cookie{
			Name:  "test1",
			Value: "test2",
		},
	}

	// returns success
	value, err := json.Marshal(wantCookies)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: value,
		}, nil)
	haveCookies, err := cookies.GetAllCookies(ctx)
	assert.Nil(t, err)
	assert.ElementsMatch(t, haveCookies, wantCookies)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(nil, cookieErr)
	haveCookies, err = cookies.GetAllCookies(ctx)
	assert.Len(t, haveCookies, 0)
	assert.Equal(t, cookieErr, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie", nil).Times(1).Return(
		&Response{
			Value: []byte(`[{":123}]`),
		}, nil)
	haveCookies, err = cookies.GetAllCookies(ctx)
	assert.Error(t, err)
	assert.Len(t, haveCookies, 0)
}

func TestCookies_GetNamedCookie(t *testing.T) {
	cookies, cli, done := newCookie(t, "123")
	defer done()

	ctx := context.TODO()

	wantCookie := Cookie{
		Name:  "test",
		Value: "test",
	}

	// returns success
	value, err := json.Marshal(wantCookie)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: value,
		}, nil)
	haveCookie, err := cookies.GetNamedCookie(ctx, "test")
	assert.Nil(t, err)
	assert.Equal(t, wantCookie.Name, haveCookie.Name)
	assert.Equal(t, wantCookie.Value, haveCookie.Value)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(nil, cookieErr)
	haveCookie, err = cookies.GetNamedCookie(ctx, "test")
	assert.Equal(t, cookieErr, err)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/cookie/test", nil).Times(1).Return(
		&Response{
			Value: []byte(`[`),
		}, nil)
	haveCookie, err = cookies.GetNamedCookie(ctx, "test")
	assert.Error(t, err)
	assert.Empty(t, haveCookie.Value)
}

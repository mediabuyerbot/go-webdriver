package protocol

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var contextErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newContext(t *testing.T, sessID string) (Context, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewContext(cli, sessID)
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestSessionContext_Types(t *testing.T) {
	assert.Equal(t, WindowHandle("133").String(), "133")
	assert.Equal(t, FrameHandle("133").String(), "133")
	assert.Equal(t, WindowType("tab").String(), "tab")
	assert.Equal(t, WindowType("window").String(), "window")
	assert.Equal(t, Window{Handle: WindowHandle("133")}.String(), "133")
	assert.Equal(t, Frame{FrameHandle("133")}.String(), "133")
}

func TestSessionContext_CloseWindow(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/window", nil).Times(1).Return(
		&Response{
			SessionID: "",
			Value:     []byte(`["CDwindow-2AAB7036D24C6759EED24640084481B7", "CDwindow-2AAB7036D24C6759EED24640084481B8"]`),
		}, nil)

	windowHandles, err := cx.CloseWindow(ctx)
	assert.Nil(t, err)
	assert.Len(t, windowHandles, 2, "window handles length")
	assert.ElementsMatch(t, windowHandles, []WindowHandle{"CDwindow-2AAB7036D24C6759EED24640084481B7", "CDwindow-2AAB7036D24C6759EED24640084481B8"})

	// returns error
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/window", nil).Times(1).Return(nil, contextErr)
	windowHandles, err = cx.CloseWindow(ctx)
	assert.Error(t, err)
	assert.Equal(t, err, contextErr)
	assert.Nil(t, windowHandles)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodDelete, "/session/123/window", nil).Times(1).Return(
		&Response{
			SessionID: "",
			Value:     []byte(`[CDwindow-2AAB7036D24C6759EED24640084481B7"]`),
		}, nil)
	windowHandles, err = cx.CloseWindow(ctx)
	assert.Error(t, err)
}

func TestSessionContext_Fullscreen(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/fullscreen", nil).Times(1).Return(
		&Response{
			Value: []byte(`{
                  "height": 1050,
                  "width": 1680,
                  "x": -2000,
                  "y": 4000
           }`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Nil(t, p)
	})

	r, err := cx.Fullscreen(ctx)
	assert.Nil(t, err)
	assert.Equal(t, r.Height, 1050)
	assert.Equal(t, r.Width, 1680)
	assert.Equal(t, r.X, -2000)
	assert.Equal(t, r.Y, 4000)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/fullscreen", nil).Times(1).Return(nil, contextErr)
	r, err = cx.Fullscreen(ctx)
	assert.Error(t, err)

	// returns success (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/fullscreen", nil).Times(1).Return(
		&Response{
			Value: []byte(`{
                  "height: 1050,
                  "width: 1680,
                  "x: -2000,
                  "y: 4000
           }`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Nil(t, p)
	})

	_, err = cx.Fullscreen(ctx)
	assert.Error(t, err)
}

func TestSessionContext_Maximize(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/maximize", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := cx.Maximize(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/maximize", nil).Times(1).Return(nil,
		contextErr)
	err = cx.Maximize(ctx)
	assert.Error(t, err)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/maximize", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cx.Maximize(ctx)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestSessionContext_Minimize(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/minimize", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := cx.Minimize(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/minimize", nil).Times(1).Return(nil,
		contextErr)
	err = cx.Minimize(ctx)
	assert.Error(t, err)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/minimize", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = cx.Minimize(ctx)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestSessionContext_GetRect(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/rect", nil).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value: []byte(`
                {
                  "height": 1006,
                  "width": 1200,
                  "x": 22,
                  "y": 23
                }
           `),
		}, nil)

	r, err := cx.GetRect(ctx)
	assert.Nil(t, err)
	assert.Equal(t, r.Height, 1006)
	assert.Equal(t, r.Width, 1200)
	assert.Equal(t, r.X, 22)
	assert.Equal(t, r.Y, 23)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/rect", nil).Times(1).Return(nil, contextErr)
	r, err = cx.GetRect(ctx)
	assert.Error(t, err)
	assert.Equal(t, err, contextErr)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/rect", nil).Times(1).Return(
		&Response{
			SessionID: "123",
			Status:    0,
			Value:     []byte(`{null`),
		}, nil)
	r, err = cx.GetRect(ctx)
	assert.Error(t, err)
	assert.Equal(t, r.Height, 0)
	assert.Equal(t, r.Width, 0)
	assert.Equal(t, r.X, 0)
	assert.Equal(t, r.Y, 0)
}

func TestSessionContext_SetRect(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	newRect := Rect{
		Width:  200,
		Height: 200,
		X:      1,
		Y:      2,
	}
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/rect", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, p["width"], newRect.Width)
		assert.Equal(t, p["height"], newRect.Height)
		assert.Equal(t, p["x"], newRect.X)
		assert.Equal(t, p["y"], newRect.Y)
	})
	err := cx.SetRect(ctx, newRect)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/rect", gomock.Any()).Times(1).Return(nil, contextErr)
	err = cx.SetRect(ctx, newRect)
	assert.Error(t, err)
	assert.Equal(t, err, contextErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/rect", gomock.Any()).Times(1).Return(
		&Response{Value: []byte(`{}`)}, nil)
	err = cx.SetRect(ctx, newRect)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestSessionContext_SwitchToFrame(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	frameID := FrameHandle("34")
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, p["id"], frameID)
	})
	err := cx.SwitchToFrame(ctx, frameID)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame", gomock.Any()).Times(1).Return(nil, contextErr)
	err = cx.SwitchToFrame(ctx, frameID)
	assert.Equal(t, err, contextErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame", gomock.Any()).Times(1).Return(&Response{Value: []byte(`{}`)}, nil)
	err = cx.SwitchToFrame(ctx, frameID)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestSessionContext_SwitchToParentFrame(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame/parent", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Nil(t, p)
	})
	err := cx.SwitchToParentFrame(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame/parent", nil).Times(1).Return(nil, contextErr)
	err = cx.SwitchToParentFrame(ctx)
	assert.Equal(t, err, contextErr)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/frame/parent", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Nil(t, p)
	})
	err = cx.SwitchToParentFrame(ctx)
	assert.Equal(t, err, ErrInvalidResponse)
}

func TestSessionContext_GetWindowHandle(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botWindowHandle := WindowHandle("CDwindow-1BCAB31FFE62561727B38152C27A7B88")
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window", nil).Times(1).Return(
		&Response{
			Value: []byte(botWindowHandle),
		}, nil)
	handle, err := cx.GetWindowHandle(ctx)
	assert.Nil(t, err)
	assert.Equal(t, handle, botWindowHandle)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window", nil).Times(1).Return(nil, contextErr)
	handle, err = cx.GetWindowHandle(ctx)
	assert.Empty(t, handle)
	assert.Error(t, err)
}

func TestSessionContext_GetWindowHandles(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botWindowHandles := []WindowHandle{
		"CDwindow-1BCAB31FFE62561727B38152C27A7B88",
		"CDwindow-1BCAB31FFE62561727B38152C27A7B89",
	}
	value, err := json.Marshal(botWindowHandles)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/handles", nil).Times(1).Return(&Response{
		Value: value,
	}, nil)
	handles, err := cx.GetWindowHandles(ctx)
	assert.Nil(t, err)
	assert.ElementsMatch(t, botWindowHandles, handles)

	// returns success (empty list)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/handles", nil).Times(1).Return(&Response{
		Value: []byte(`[]`),
	}, nil)
	handles, err = cx.GetWindowHandles(ctx)
	assert.Nil(t, err)
	assert.ElementsMatch(t, handles, []WindowHandle{})

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/handles", nil).Times(1).Return(nil, contextErr)
	handles, err = cx.GetWindowHandles(ctx)
	assert.Equal(t, err, contextErr)
	assert.ElementsMatch(t, handles, []WindowHandle{})

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/window/handles", nil).Times(1).Return(&Response{Value: []byte(`[`)}, nil)
	handles, err = cx.GetWindowHandles(ctx)
	assert.Error(t, err)
	assert.Empty(t, handles)
}

func TestSessionContext_NewWindow(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botWindow := Window{
		Handle: WindowHandle("CDwindow-5B07A00849E4B4DB05A83B074747A172"),
		Type:   Tab,
	}
	value, err := json.Marshal(botWindow)
	assert.Nil(t, err)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/new", gomock.Any()).Times(1).Return(
		&Response{
			Value: value,
		}, nil)
	win, err := cx.NewWindow(ctx)
	assert.Nil(t, err)
	assert.Equal(t, win.Handle, botWindow.Handle)
	assert.Equal(t, win.Type, botWindow.Type)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/new", gomock.Any()).Times(1).Return(nil, contextErr)
	win, err = cx.NewWindow(ctx)
	assert.Equal(t, err, contextErr)
	assert.Equal(t, win.Handle, Window{}.Handle)

	// returns error (invalid response)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window/new", gomock.Any()).Times(1).Return(&Response{Value: []byte(`{`)}, nil)
	win, err = cx.NewWindow(ctx)
	assert.Error(t, err)
	assert.Equal(t, win.Handle, Window{}.Handle)
}

func TestSessionContext_SwitchToWindow(t *testing.T) {
	cx, cli, done := newContext(t, "123")
	defer done()

	ctx := context.TODO()

	// returns success
	botWindowHandle := WindowHandle("CDwindow-5B07A00849E4B4DB05A83B074747A172")
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window", gomock.Any()).Times(1).Return(
		&Response{Value: []byte(`null`)}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, p["name"], botWindowHandle)
	})
	err := cx.SwitchToWindow(ctx, botWindowHandle)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window", gomock.Any()).Times(1).Return(nil, contextErr)
	err = cx.SwitchToWindow(ctx, botWindowHandle)
	assert.Equal(t, err, contextErr)

	// returns error (bad JSON format)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/window", gomock.Any()).Times(1).Return(
		&Response{Value: []byte(`{}`)}, nil)
	err = cx.SwitchToWindow(ctx, botWindowHandle)
	assert.Equal(t, err, ErrInvalidResponse)
}

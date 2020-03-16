package w3c

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Context represents a window handle with a unique identifier.
type Context interface {

	// GetWindowHandle returns the current window handle.
	GetWindowHandle(context.Context) (WindowHandle, error)

	// GetWindowHandles returns the list of all window handles available to the session.
	GetWindowHandles(context.Context) ([]WindowHandle, error)

	// CloseWindow close the current window.
	CloseWindow(context.Context) ([]WindowHandle, error)

	// NewWindow create a new top-level browsing context.
	NewWindow(context.Context, WindowType) (Window, error)

	// SwitchToWindow switching window will select the current top-level browsing context used as the target
	// for all subsequent commands. In a tabbed browser, this will typically make the tab containing
	// the browsing context the selected tab.
	SwitchToWindow(context.Context, WindowHandle) error

	// SwitchToFrame change focus to another frame on the page.
	SwitchToFrame(context.Context, FrameHandle) error

	// SwitchToParentFrame change focus back to parent frame.
	SwitchToParentFrame(context.Context) error

	// SetRect alters the size and the position of the operating system window
	// corresponding to the current top-level browsing context.
	SetRect(context.Context, Rect) (Rect, error)

	// returns the size and position on the screen of the operating system
	// window corresponding to the current top-level browsing context.
	GetRect(context.Context) (Rect, error)

	// Maximize invokes the window manager-specific “maximize” operation,
	// if any, on the window containing the current top-level browsing context.
	// This typically increases the window to the maximum available size without going full-screen.
	Maximize(context.Context) (Rect, error)

	// Minimize invokes the window manager-specific “minimize” operation,
	// if any, on the window containing the current top-level browsing context.
	// This typically hides the window in the system tray.
	Minimize(context.Context) (Rect, error)

	// Fullscreen fullscreen mode.
	Fullscreen(context.Context) (Rect, error)
}

const (
	Tab WindowType = "tab"
	Win WindowType = "window"
)

type (
	WindowType   string
	WindowHandle string
	FrameHandle  string
)

func (wh WindowHandle) String() string {
	return string(wh)
}

func (wh WindowHandle) IsEmpty() bool {
	return len(wh) == 0
}

func (wt WindowType) String() string {
	return string(wt)
}

func (wt WindowType) Validate() error {
	switch wt {
	case Tab, Win:
		return nil
	default:
		return errors.New("w3c: unknown window type")
	}
}

func (fh FrameHandle) String() string {
	return string(fh)
}

type Rect struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Y      int `json:"y"`
	X      int `json:"x"`
}

type Window struct {
	Handle WindowHandle `json:"handle"`
	Type   WindowType   `json:"type"`
}

func (w Window) String() string {
	return string(w.Handle)
}

type Frame struct {
	Handle FrameHandle `json:"handle"`
}

func (f Frame) String() string {
	return f.Handle.String()
}

type sessionContext struct {
	id      string
	request Doer
}

// NewContext creates a new instance of Context.
func NewContext(doer Doer, sessID string) Context {
	return &sessionContext{
		id:      sessID,
		request: doer,
	}
}

func (c *sessionContext) NewWindow(ctx context.Context, wt WindowType) (w Window, err error) {
	if err := wt.Validate(); err != nil {
		return w, err
	}
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window/new", Params{"type": wt})
	if err != nil {
		return w, err
	}
	if err := json.Unmarshal(resp.Value, &w); err != nil {
		return w, err
	}
	return w, nil
}

func (c *sessionContext) GetWindowHandle(ctx context.Context) (wh WindowHandle, err error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/window", nil)
	if err != nil {
		return wh, err
	}
	if err := json.Unmarshal(resp.Value, &wh); err != nil {
		return wh, err
	}
	return wh, nil
}

func (c *sessionContext) GetWindowHandles(ctx context.Context) ([]WindowHandle, error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/window/handles", nil)
	if err != nil {
		return nil, err
	}
	var handles []string
	if err := json.Unmarshal(resp.Value, &handles); err != nil {
		return nil, err
	}
	windowHandles := make([]WindowHandle, len(handles))
	for i, handle := range handles {
		windowHandles[i] = WindowHandle(handle)
	}
	return windowHandles, nil
}

func (c *sessionContext) CloseWindow(ctx context.Context) ([]WindowHandle, error) {
	resp, err := c.request.Do(ctx, http.MethodDelete, "/session/"+c.id+"/window", nil)
	if err != nil {
		return nil, err
	}
	var handles []string
	if err := json.Unmarshal(resp.Value, &handles); err != nil {
		return nil, err
	}
	windowHandles := make([]WindowHandle, len(handles))
	for i, handle := range handles {
		windowHandles[i] = WindowHandle(handle)
	}
	return windowHandles, nil
}

func (c *sessionContext) SwitchToWindow(ctx context.Context, h WindowHandle) error {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window", Params{"handle": h})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SwitchToFrame(ctx context.Context, h FrameHandle) error {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/frame", Params{"id": h})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SwitchToParentFrame(ctx context.Context) error {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/frame/parent", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SetRect(ctx context.Context, r Rect) (winRect Rect, err error) {
	p := Params{
		"width":  r.Width,
		"height": r.Height,
		"x":      r.X,
		"y":      r.Y,
	}
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window/rect", p)
	if err != nil {
		return winRect, err
	}
	if err := json.Unmarshal(resp.Value, &winRect); err != nil {
		return winRect, err
	}
	return winRect, nil
}

func (c *sessionContext) GetRect(ctx context.Context) (r Rect, err error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/window/rect", nil)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

func (c *sessionContext) Maximize(ctx context.Context) (r Rect, err error) {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window/maximize", nil)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

func (c *sessionContext) Minimize(ctx context.Context) (r Rect, err error) {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window/minimize", nil)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

func (c *sessionContext) Fullscreen(ctx context.Context) (r Rect, err error) {
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/window/fullscreen", nil)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

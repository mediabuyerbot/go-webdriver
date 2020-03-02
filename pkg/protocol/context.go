package protocol

import (
	"context"
	"encoding/json"
)

// Context represents a window handle with a unique identifier.
type Context interface {

	// GetWindowHandle returns the current window handle.
	GetWindowHandle(context.Context) (WindowHandle, error)

	// GetWindowHandles returns the list of all window handles available to the session.
	GetWindowHandles(context.Context) ([]WindowHandle, error)

	// CloseWindow close the current window.
	CloseWindow(context.Context) ([]WindowHandle, error)

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
	SetRect(context.Context, Rect) error

	// returns the size and position on the screen of the operating system
	// window corresponding to the current top-level browsing context.
	GetRect(context.Context) (Rect, error)

	// Maximize invokes the window manager-specific “maximize” operation,
	// if any, on the window containing the current top-level browsing context.
	// This typically increases the window to the maximum available size without going full-screen.
	Maximize(context.Context) error

	// Minimize invokes the window manager-specific “minimize” operation,
	// if any, on the window containing the current top-level browsing context.
	// This typically hides the window in the system tray.
	Minimize(context.Context) error

	// Fullscreen fullscreen mode.
	Fullscreen(context.Context) error
}

type Rect struct {
	Width  int  `json:"width"`
	Height int  `json:"height"`
	Y      uint `json:"y"`
	X      uint `json:"x"`
}

type WindowHandle string

func (wh WindowHandle) String() string {
	return string(wh)
}

type FrameHandle string

func (fh FrameHandle) String() string {
	return string(fh)
}

type sessionContext struct {
	id     string
	client Client
}

// NewContext creates a new instance of Context.
func NewContext(cli Client, sessID string) Context {
	return &sessionContext{
		id:     sessID,
		client: cli,
	}
}

func (c *sessionContext) GetWindowHandle(ctx context.Context) (wh WindowHandle, err error) {
	resp, err := c.client.Get(ctx, "/session/"+c.id+"/window")
	if err != nil {
		return wh, err
	}
	return WindowHandle(resp.Value), nil
}

func (c *sessionContext) GetWindowHandles(ctx context.Context) ([]WindowHandle, error) {
	resp, err := c.client.Get(ctx, "/session/"+c.id+"/window/handles")
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
	resp, err := c.client.Delete(ctx, "/session/"+c.id+"/window")
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
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/window", params{"name": h})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SwitchToFrame(ctx context.Context, h FrameHandle) error {
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/frame", params{"id": h})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SwitchToParentFrame(ctx context.Context) error {
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/frame/parent", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) SetRect(ctx context.Context, r Rect) error {
	params := params{
		"width":  r.Width,
		"height": r.Height,
		"x":      r.X,
		"y":      r.Y,
	}
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/window/rect", params)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) GetRect(ctx context.Context) (r Rect, err error) {
	resp, err := c.client.Get(ctx, "/session/"+c.id+"/window/rect")
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

func (c *sessionContext) Maximize(ctx context.Context) error {
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/window/maximize", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) Minimize(ctx context.Context) error {
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/window/minimize", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *sessionContext) Fullscreen(ctx context.Context) error {
	resp, err := c.client.Post(ctx, "/session/"+c.id+"/window/fullscreen ", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

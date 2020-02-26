package protocol

import (
	"context"
	"encoding/json"
)

const (
	// Maximized the window is maximized.
	Maximized WindowState = "maximize"

	// Minimized the window is iconified.
	Minimized WindowState = "minimized"

	// Normal the window is shown normally.
	Normal WindowState = "normal"

	// Fullscreen the window is in full screen mode.
	Fullscreen WindowState = "fullscreen"
)

type WindowState string

type Rect struct {
	Width  int  `json:"width"`
	Height int  `json:"height"`
	Y      uint `json:"y"`
	X      uint `json:"x"`
}

type Window struct {
	windowID  string
	sessionID string
	client    Client
}

func (w Window) ID() string {
	return w.windowID
}

func (w Window) SessionID() string {
	return w.sessionID
}

// SetRect change the size of the specified window.
// https://www.w3.org/TR/webdriver1/#set-window-rect
func (w Window) SetRect(ctx context.Context, r Rect) error {
	params := params{
		"width":  r.Width,
		"height": r.Height,
		"x":      r.X,
		"y":      r.Y,
	}
	if _, err := w.client.Post(ctx, "/session/"+w.sessionID+"/window/rect", params); err != nil {
		return err
	}
	return nil
}

// GetRect get the size of the specified window.
// https://www.w3.org/TR/webdriver1/#get-window-rect
func (w Window) GetRect(ctx context.Context) (r Rect, err error) {
	resp, err := w.client.Get(ctx, "/session/"+w.sessionID+"/window/rect")
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}

// Maximize the specified window if not already maximized.
// https://www.w3.org/TR/webdriver1/#maximize-window
func (w Window) Maximize(ctx context.Context) error {
	if _, err := w.client.Post(ctx, "/session/"+w.sessionID+"/window/maximize", nil); err != nil {
		return err
	}
	return nil
}

// Minimize the specified window if not already minimized.
// https://www.w3.org/TR/webdriver1/#minimize-window
func (w Window) Minimize(ctx context.Context) error {
	if _, err := w.client.Post(ctx, "/session/"+w.sessionID+"/window/minimize", nil); err != nil {
		return err
	}
	return nil
}

// Fullscreen the specified window.
// https://www.w3.org/TR/webdriver1/#fullscreen-window
func (w Window) Fullscreen(ctx context.Context) error {
	if _, err := w.client.Post(ctx, "/session/"+w.sessionID+"/window/fullscreen ", nil); err != nil {
		return err
	}
	return nil
}

package protocol

import "context"

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

type Size struct {
	Width  int
	Height int
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

func (w Window) SetSize(ctx context.Context, s Size) error {
	return nil
}

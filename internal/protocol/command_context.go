package protocol

import (
	"context"
	"encoding/json"
)

// CommandContext represents the protocol window handle.
type CommandContext struct {
	id     string
	client Client
}

// NewCommandContext creates a new instance of CommandContext.
func NewCommandContext(cli Client, id string) *CommandContext {
	return &CommandContext{
		id:     id,
		client: cli,
	}
}

// GetWindowHandle returns the current window handle.
// https://www.w3.org/TR/webdriver1/#get-window-handle
func (cc *CommandContext) GetWindowHandle(ctx context.Context) (wh Window, err error) {
	resp, err := cc.client.Get(ctx, "/session/"+cc.id+"/window")
	if err != nil {
		return wh, err
	}
	return Window{
		sessionID: cc.id,
		windowID:  string(resp.Value),
		client:    cc.client,
	}, nil
}

// GetWindowHandles returns the list of all window handles available to the session.
// https://www.w3.org/TR/webdriver1/#get-window-handles
func (cc *CommandContext) GetWindowHandles(ctx context.Context) ([]Window, error) {
	resp, err := cc.client.Get(ctx, "/session/"+cc.id+"/window/handles")
	if err != nil {
		return nil, err
	}

	var windowIDs []string
	if err := json.Unmarshal(resp.Value, &windowIDs); err != nil {
		return nil, err
	}

	windows := make([]Window, len(windowIDs))
	for i, wid := range windowIDs {
		windows[i] = Window{
			client:    cc.client,
			sessionID: cc.id,
			windowID:  wid,
		}
	}
	return windows, nil
}

// CloseWindow close the current window.
// https://www.w3.org/TR/webdriver1/#close-window
func (cc *CommandContext) CloseWindow(ctx context.Context) error {
	if _, err := cc.client.Delete(ctx, "/session/"+cc.id+"/window"); err != nil {
		return err
	}
	return nil
}

// SwitchToWindow change focus to another window. The window to change focus to may be specified by its server
// assigned window handle, or by the value of its name attribute.
// https://www.w3.org/TR/webdriver1/#switch-to-window
func (cc *CommandContext) SwitchToWindow(ctx context.Context, name string) error {
	if _, err := cc.client.Post(ctx, "/session/"+cc.id+"/window", params{"name": name}); err != nil {
		return err
	}
	return nil
}

// SwitchToFrame change focus to another frame on the page.
// https://www.w3.org/TR/webdriver1/#switch-to-frame
func (cc *CommandContext) SwitchToFrame(ctx context.Context, frameID string) error {
	if _, err := cc.client.Post(ctx, "/session/"+cc.id+"/frame", params{"id": frameID}); err != nil {
		return err
	}
	return nil
}

// SwitchToParentFrame change focus back to parent frame.
// https://www.w3.org/TR/webdriver1/#switch-to-parent-frame
func (cc *CommandContext) SwitchToParentFrame(ctx context.Context) error {
	if _, err := cc.client.Post(ctx, "/session/"+cc.id+"/frame/parent", nil); err != nil {
		return err
	}
	return nil
}

package protocol

import (
	"context"
	"net/http"
)

type Alert interface {

	// Dismiss dismisses a simple dialog if present.
	// A request to dismiss an alert user prompt, which may not necessarily have a dismiss button,
	// has the same effect as accepting it.
	Dismiss(ctx context.Context) error

	// Accept accepts the currently displayed alert dialog.
	Accept(ctx context.Context) error

	// Text gets the text of the currently displayed JavaScript alert(), confirm(), or prompt() dialog.
	Text(ctx context.Context) (string, error)

	// SetText sets the text field of a window.prompt user prompt to the given value.
	SetText(ctx context.Context, text string) error
}

type alert struct {
	id      string
	request Doer
}

// NewAlert creates a new instance of Alert.
func NewAlert(doer Doer, sessID string) Alert {
	return &alert{
		id:      sessID,
		request: doer,
	}
}

func (a *alert) Dismiss(ctx context.Context) error {
	resp, err := a.request.Do(ctx, http.MethodPost, "/session/"+a.id+"/alert/dismiss", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (a *alert) Accept(ctx context.Context) error {
	resp, err := a.request.Do(ctx, http.MethodPost, "/session/"+a.id+"/alert/accept", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (a *alert) Text(ctx context.Context) (string, error) {
	resp, err := a.request.Do(ctx, http.MethodGet, "/session/"+a.id+"/alert/text", nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

func (a *alert) SetText(ctx context.Context, text string) error {
	resp, err := a.request.Do(ctx, http.MethodPost, "/session/"+a.id+"/alert/text", Params{"text": text})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

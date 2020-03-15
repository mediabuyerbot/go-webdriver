package w3c

import (
	"context"
	"net/http"
)

// Navigation represents a navigation of the current session context
type Navigation interface {

	// NavigateTo navigate to a new URL.
	NavigateTo(ctx context.Context, url string) (err error)

	// GetCurrentURL returns the URL of the current page.
	GetCurrentURL(ctx context.Context) (url string, err error)

	// Back navigate backwards in the browser history, if possible.
	Back(ctx context.Context) error

	// Forward navigate forwards in the browser history, if possible.
	Forward(ctx context.Context) error

	// Refresh refresh the current page.
	Refresh(ctx context.Context) error

	// GetTitle returns the current page title.
	GetTitle(ctx context.Context) (title string, err error)
}

type navigation struct {
	id      string
	request Doer
}

// NewNavigation creates a new instance of Navigation.
func NewNavigation(doer Doer, sessID string) Navigation {
	return &navigation{
		id:      sessID,
		request: doer,
	}
}

func (n *navigation) NavigateTo(ctx context.Context, url string) (err error) {
	resp, err := n.request.Do(ctx, http.MethodPost, "/session/"+n.id+"/url", Params{"url": url})
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (n *navigation) GetCurrentURL(ctx context.Context) (s string, err error) {
	resp, err := n.request.Do(ctx, http.MethodGet, "/session/"+n.id+"/url", nil)
	if err != nil {
		return s, err
	}
	return string(resp.Value), nil
}

func (n *navigation) Back(ctx context.Context) error {
	resp, err := n.request.Do(ctx, http.MethodPost, "/session/"+n.id+"/back", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (n *navigation) Forward(ctx context.Context) error {
	resp, err := n.request.Do(ctx, http.MethodPost, "/session/"+n.id+"/forward", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (n *navigation) Refresh(ctx context.Context) error {
	resp, err := n.request.Do(ctx, http.MethodPost, "/session/"+n.id+"/refresh", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (n *navigation) GetTitle(ctx context.Context) (title string, err error) {
	resp, err := n.request.Do(ctx, http.MethodGet, "/session/"+n.id+"/title", nil)
	if err != nil {
		return title, err
	}
	return string(resp.Value), nil
}

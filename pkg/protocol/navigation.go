package protocol

import "context"

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
	id     string
	client Client
}

// NewNavigation creates a new instance of Navigation.
func NewNavigation(cli Client, sessID string) Navigation {
	return &navigation{
		id:     sessID,
		client: cli,
	}
}

func (n *navigation) NavigateTo(ctx context.Context, url string) (err error) {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/url", params{"url": url}); err != nil {
		return err
	}
	return nil
}

func (n *navigation) GetCurrentURL(ctx context.Context) (s string, err error) {
	resp, err := n.client.Get(ctx, "/session/"+n.id+"/url")
	if err != nil {
		return s, err
	}
	return string(resp.Value), nil
}

func (n *navigation) Back(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/back", nil); err != nil {
		return err
	}
	return nil
}

func (n *navigation) Forward(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/forward", nil); err != nil {
		return err
	}
	return nil
}

func (n *navigation) Refresh(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/refresh", nil); err != nil {
		return err
	}
	return nil
}

func (n *navigation) GetTitle(ctx context.Context) (t string, err error) {
	resp, err := n.client.Get(ctx, "/session/"+n.id+"/title")
	if err != nil {
		return t, err
	}
	return string(resp.Value), nil
}

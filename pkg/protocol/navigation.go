package protocol

import "context"

// Navigation represents the navigation of the current top-level browsing context to new URLs
// and introspection of the document currently loaded in this browsing context.
// https://www.w3.org/TR/webdriver1/#navigation
type Navigation struct {
	id     string
	client Client
}

// NewNavigation creates a new instance of Navigation.
func NewNavigation(cli Client, sessID string) *Navigation {
	return &Navigation{
		id:     sessID,
		client: cli,
	}
}

// Url navigate to a new URL.
// https://www.w3.org/TR/webdriver1/#navigate-to
func (n *Navigation) Url(ctx context.Context, url string) error {
	params := params{"url": url}
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/url", params); err != nil {
		return err
	}
	return nil
}

// GetCurrentURL returns the URL of the current page.
// https://www.w3.org/TR/webdriver1/#get-current-url
func (n *Navigation) GetCurrentURL(ctx context.Context) (s string, err error) {
	resp, err := n.client.Get(ctx, "/session/"+n.id+"/url")
	if err != nil {
		return s, err
	}
	return string(resp.Value), nil
}

// Back navigate backwards in the browser history, if possible.
// https://www.w3.org/TR/webdriver1/#back
func (n *Navigation) Back(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/back", nil); err != nil {
		return err
	}
	return nil
}

// Forward navigate forwards in the browser history, if possible.
// https://www.w3.org/TR/webdriver1/#forward
func (n *Navigation) Forward(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/forward", nil); err != nil {
		return err
	}
	return nil
}

// Refresh refresh the current page.
// https://www.w3.org/TR/webdriver1/#refresh
func (n *Navigation) Refresh(ctx context.Context) error {
	if _, err := n.client.Post(ctx, "/session/"+n.id+"/refresh", nil); err != nil {
		return err
	}
	return nil
}

// GetTitle returns the current page title.
// https://www.w3.org/TR/webdriver1/#get-title
func (n *Navigation) GetTitle(ctx context.Context) (t string, err error) {
	resp, err := n.client.Get(ctx, "/session/"+n.id+"/title")
	if err != nil {
		return t, err
	}
	return string(resp.Value), nil
}

package webdriver

import (
	"context"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
)

// Remote creates a new instance of the remote browser.
func Remote(addr string, opts w3c.BrowserOptions) (*Browser, error) {
	sess, err := NewSession(addr, opts)
	if err != nil {
		return nil, err
	}
	return &Browser{
		ctx:    context.Background(),
		driver: nil,
		sess:   sess,
	}, nil
}

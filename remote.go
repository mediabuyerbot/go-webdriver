package webdriver

import (
	"context"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3cproto"
)

// OpenRemoteBrowser creates a new instance of the remote browser.
func OpenRemoteBrowser(ctx context.Context, addr string, opts w3cproto.BrowserOptions) (*Browser, error) {
	sess, err := NewSession(ctx, addr, opts)
	if err != nil {
		return nil, err
	}
	return &Browser{
		ctx:    context.Background(),
		driver: nil,
		sess:   sess,
	}, nil
}

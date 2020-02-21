package webdriver

import (
	"context"
	"time"

	"github.com/mediabuyerbot/go-webdriver/internal/protocol"
	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
)

const (
	DefaultRetryCount = 3
	DefaultTimeout    = 15 * time.Second
)

type Browser struct {
	session    *protocol.Session
	navigation *protocol.Navigation
	commandCtx *protocol.CommandContext
}

func NewBrowserFromClient(client httpclient.Client, desired, required protocol.Capabilities) (*Browser, error) {
	cli := protocol.NewClient(client)
	sess, err := protocol.NewSession(cli, desired, required)
	if err != nil {
		return nil, err
	}

	browser := Browser{
		session:    sess,
		navigation: protocol.NewNavigation(cli, sess.ID()),
		commandCtx: protocol.NewCommandContext(cli, sess.ID()),
	}
	return &browser, nil
}

func NewBrowser(addr string, desired, required protocol.Capabilities) (*Browser, error) {
	client, err := httpclient.NewClient(addr,
		httpclient.WithRetryCount(DefaultRetryCount),
		httpclient.WithTimeout(DefaultTimeout),
	)
	if err != nil {
		return nil, err
	}
	return NewBrowserFromClient(client, desired, required)
}

func (b *Browser) Session() *protocol.Session {
	return b.session
}

func (b *Browser) Navigation() *protocol.Navigation {
	return b.navigation
}

func (b *Browser) CommandContext() *protocol.CommandContext {
	return b.commandCtx
}

func (b *Browser) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

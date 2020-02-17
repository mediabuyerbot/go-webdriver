package webdriver

import (
	"context"
	"time"

	"github.com/mediabuyerbot/go-webdriver/httpclient"
	"github.com/mediabuyerbot/go-webdriver/internal/protocol"
)

const (
	DefaultRetryCount = 3
	DefaultTimeout    = 15 * time.Second
)

type Browser struct {
	session *protocol.Session
}

func NewBrowserFromClient(client httpclient.Client, desired, required protocol.Capabilities) (*Browser, error) {
	cli := protocol.NewClient(client)
	sess, err := protocol.NewSession(cli, desired, required)
	if err != nil {
		return nil, err
	}

	browser := Browser{
		session: sess,
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

func (b *Browser) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

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

type Session struct {
	session    *protocol.Session
	navigation *protocol.Navigation
	commandCtx *protocol.CommandContext
}

func NewSessionFromClient(client httpclient.Client, desired, required protocol.Capabilities) (*Session, error) {
	cli := protocol.NewClient(client)
	sess, err := protocol.NewSession(cli, desired, required)
	if err != nil {
		return nil, err
	}

	browser := Session{
		session:    sess,
		navigation: protocol.NewNavigation(cli, sess.ID()),
		commandCtx: protocol.NewCommandContext(cli, sess.ID()),
	}
	return &browser, nil
}

func NewSession(addr string, desired, required protocol.Capabilities) (*Session, error) {
	client, err := httpclient.NewClient(addr,
		httpclient.WithRetryCount(DefaultRetryCount),
		httpclient.WithTimeout(DefaultTimeout),
	)
	if err != nil {
		return nil, err
	}
	return NewSessionFromClient(client, desired, required)
}

func (b *Session) Session() *protocol.Session {
	return b.session
}

func (b *Session) Navigation() *protocol.Navigation {
	return b.navigation
}

func (b *Session) CommandContext() *protocol.CommandContext {
	return b.commandCtx
}

func (b *Session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

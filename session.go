package webdriver

import (
	"context"
	"net/http"
	"time"

	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
	"github.com/mediabuyerbot/go-webdriver/pkg/protocol"
)

const (
	DefaultRetryCount = 3
	DefaultBackOff    = 5
	DefaultTimeout    = 15 * time.Second
)

// Session represents a session is equivalent to a single instantiation of a particular user agent,
// including all its child browsers.
type Session interface {
	// SessionID returns the unique session id.
	SessionID() string

	// Session returns an implemented session protocol.
	Session() protocol.Session

	// Timeouts returns an implemented timeouts protocol.
	Timeouts() protocol.Timeouts

	// Close close the current session.
	Close(ctx context.Context) error
}

type (
	DesiredCapabilities  map[string]interface{}
	RequiredCapabilities map[string]interface{}
)

type session struct {
	session    protocol.Session
	timeouts   protocol.Timeouts
	navigation *protocol.Navigation
	commandCtx *protocol.CommandContext
}

func NewSessionFromClient(client httpclient.Client, d DesiredCapabilities, r RequiredCapabilities) (Session, error) {
	cli := protocol.NewClient(client)
	sess, err := protocol.NewSession(cli, d, r)
	if err != nil {
		return nil, err
	}

	browser := session{
		session:  sess,
		timeouts: protocol.NewTimeouts(cli, sess.ID()),
		// navigation: protocol.NewNavigation(cli, sess.ID()),
		// commandCtx: protocol.NewCommandContext(cli, sess.ID()),
	}
	return &browser, nil
}

func NewSession(addr string, d DesiredCapabilities, r RequiredCapabilities) (Session, error) {
	client, err := httpclient.NewClient(addr,
		httpclient.WithRetryCount(DefaultRetryCount),
		httpclient.WithTimeout(DefaultTimeout),
		httpclient.WithBackOff(func(attemptNum int, resp *http.Response) time.Duration {
			return DefaultBackOff * time.Second
		}),
	)
	if err != nil {
		return nil, err
	}
	return NewSessionFromClient(client, d, r)
}

func (b *session) SessionID() string {
	return b.session.ID()
}

func (b *session) Session() protocol.Session {
	return b.session
}

func (b *session) Timeouts() protocol.Timeouts {
	return b.timeouts
}

//func (b *Session) Navigation() *protocol.Navigation {
//	return b.navigation
//}
//
//func (b *Session) CommandContext() *protocol.CommandContext {
//	return b.commandCtx
//}

func (b *session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

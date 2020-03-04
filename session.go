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

	// Session returns a session protocol.
	Session() protocol.Session

	// Timeouts returns a timeouts protocol.
	Timeouts() protocol.Timeouts

	// Navigation returns a navigation protocol.
	Navigation() protocol.Navigation

	// Context returns a context protocol.
	Context() protocol.Context

	// Cookies returns a cookies protocol.
	Cookies() protocol.Cookies

	// Document returns a document protocol.
	Document() protocol.Document

	// ScreenCapture returns a screen capture protocol.
	ScreenCapture() protocol.ScreenCapture

	// Close close the current session.
	Close(ctx context.Context) error
}

type (
	DesiredCapabilities  map[string]interface{}
	RequiredCapabilities map[string]interface{}
)

type session struct {
	session       protocol.Session
	timeouts      protocol.Timeouts
	navigation    protocol.Navigation
	context       protocol.Context
	cookies       protocol.Cookies
	document      protocol.Document
	screenCapture protocol.ScreenCapture
}

func NewSessionFromClient(client httpclient.Client, d DesiredCapabilities, r RequiredCapabilities) (Session, error) {
	cli := protocol.NewTransport(client)
	sess, err := protocol.NewSession(cli, d, r)
	if err != nil {
		return nil, err
	}

	browser := session{
		session:       sess,
		timeouts:      protocol.NewTimeouts(cli, sess.ID()),
		navigation:    protocol.NewNavigation(cli, sess.ID()),
		context:       protocol.NewContext(cli, sess.ID()),
		cookies:       protocol.NewCookies(cli, sess.ID()),
		document:      protocol.NewDocument(cli, sess.ID()),
		screenCapture: protocol.NewScreenCapture(cli, sess.ID()),
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

// SessionID returns the unique session id.
func (b *session) SessionID() string {
	return b.session.ID()
}

// Session returns a session protocol.
func (b *session) Session() protocol.Session {
	return b.session
}

// Timeouts returns a timeouts protocol.
func (b *session) Timeouts() protocol.Timeouts {
	return b.timeouts
}

// Navigation returns a navigation protocol.
func (b *session) Navigation() protocol.Navigation {
	return b.navigation
}

// Context returns a context protocol.
func (b *session) Context() protocol.Context {
	return b.context
}

// Cookies returns a cookies protocol.
func (b *session) Cookies() protocol.Cookies {
	return b.cookies
}

// Document returns a document protocol.
func (b *session) Document() protocol.Document {
	return b.document
}

// ScreenCapture returns a screen capture protocol.
func (b *session) ScreenCapture() protocol.ScreenCapture {
	return b.screenCapture
}

// Close close the current session.
func (b *session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

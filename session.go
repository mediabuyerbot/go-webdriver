package webdriver

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"time"

	"github.com/mediabuyerbot/go-webdriver/pkg/httpclient"
	"github.com/mediabuyerbot/go-webdriver/pkg/protocol"
)

const (
	DefaultRetryCount = 3
	DefaultBackOff    = 5
	DefaultTimeout    = 15 * time.Second
)

type (
	DesiredCapabilities  map[string]interface{}
	RequiredCapabilities map[string]interface{}
)

type Session struct {
	session       protocol.Session
	timeouts      protocol.Timeouts
	navigation    protocol.Navigation
	context       protocol.Context
	cookies       protocol.Cookies
	document      protocol.Document
	screenCapture protocol.ScreenCapture
}

func NewSessionFromClient(client httpclient.Client, d DesiredCapabilities, r RequiredCapabilities) (*Session, error) {
	cli := protocol.NewTransport(client)
	sess, err := protocol.NewSession(cli, d, r)
	if err != nil {
		return nil, err
	}

	browser := Session{
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

func NewSession(addr string, d DesiredCapabilities, r RequiredCapabilities) (*Session, error) {
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
func (b *Session) SessionID() string {
	return b.session.ID()
}

// ScreenshotPNG take a screenshot of the current page.
func (b *Session) ScreenshotPNG(ctx context.Context) (image.Image, error) {
	reader, err := b.screenCapture.Take(ctx)
	if err != nil {
		return nil, err
	}
	return png.Decode(reader)
}

// ScreenshotJPG take a screenshot of the current page.
func (b *Session) ScreenshotJPG(ctx context.Context) (image.Image, error) {
	reader, err := b.screenCapture.Take(ctx)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(reader)
}

// Url navigate to a new URL.
func (b *Session) Url(ctx context.Context, u *url.URL) error {
	return b.navigation.NavigateTo(ctx, u.String())
}

// Title returns the current page title.
func (b *Session) Title(ctx context.Context) (string, error) {
	return b.navigation.GetTitle(ctx)
}

// ExecScript inject a snippet of JavaScript into the page for execution in the context of the currently
// selected frame. The executed script is assumed to be synchronous and the result of evaluating the script
// is returned to the client. The script argument defines the script to execute in the form of a function body.
// The value returned by that function will be returned to the client. The function will be invoked with
// the provided args array and the values may be accessed via the arguments object in the order specified.
// Arguments may be any JSON-primitive, array, or JSON object.
func (b *Session) ExecScript(ctx context.Context, script string, args []interface{}) ([]byte, error) {
	return b.document.ExecuteScript(ctx, script, args)
}

// ExecAsyncScript  inject a snippet of JavaScript into the page for execution in the context of the
// currently selected frame. The executed script is assumed to be asynchronous and must signal that
// is done by invoking the provided callback, which is always provided as the final argument to the function.
// The value to this callback will be returned to the client. Asynchronous script commands may not span page loads.
// If an unload event is fired while waiting for a script result, an error should be returned to the client.
// The script argument defines the script to execute in teh form of a function body. The function will be invoked
// with the provided args array and the values may be accessed via the arguments object in the order specified.
// The final argument will always be a callback function that must be invoked to signal that the script has finished.
// Arguments may be any JSON-primitive, array, or JSON object.
func (b *Session) ExecAsyncScript(ctx context.Context, script string, args []interface{}) ([]byte, error) {
	return b.document.ExecuteAsyncScript(ctx, script, args)
}

// Session returns a session protocol.
func (b *Session) Session() protocol.Session {
	return b.session
}

// Timeouts returns a timeouts protocol.
func (b *Session) Timeouts() protocol.Timeouts {
	return b.timeouts
}

// Navigation returns a navigation protocol.
func (b *Session) Navigation() protocol.Navigation {
	return b.navigation
}

// Context returns a context protocol.
func (b *Session) Context() protocol.Context {
	return b.context
}

// Cookies returns a cookies protocol.
func (b *Session) Cookies() protocol.Cookies {
	return b.cookies
}

// Document returns a document protocol.
func (b *Session) Document() protocol.Document {
	return b.document
}

// ScreenCapture returns a screen capture protocol.
func (b *Session) ScreenCapture() protocol.ScreenCapture {
	return b.screenCapture
}

// Close close the current session.
func (b *Session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

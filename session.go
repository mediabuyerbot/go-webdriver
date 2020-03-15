package webdriver

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
	"github.com/mediabuyerbot/httpclient"
)

type Session struct {
	session       w3c.Session
	timeouts      w3c.Timeouts
	navigation    w3c.Navigation
	context       w3c.Context
	cookies       w3c.Cookies
	document      w3c.Document
	screenCapture w3c.ScreenCapture
	elements      w3c.Elements
}

func NewSessionFromClient(client httpclient.Client, opts w3c.BrowserOptions) (*Session, error) {
	cli := w3c.WithClient(client)
	sess, err := w3c.NewSession(cli, opts)
	if err != nil {
		return nil, err
	}

	browser := Session{
		session:       sess,
		timeouts:      w3c.NewTimeouts(cli, sess.ID()),
		navigation:    w3c.NewNavigation(cli, sess.ID()),
		context:       w3c.NewContext(cli, sess.ID()),
		cookies:       w3c.NewCookies(cli, sess.ID()),
		document:      w3c.NewDocument(cli, sess.ID()),
		elements:      w3c.NewElements(cli, sess.ID()),
		screenCapture: w3c.NewScreenCapture(cli, sess.ID()),
	}
	return &browser, nil
}

func NewSession(addr string, opts w3c.BrowserOptions) (*Session, error) {
	client, err := httpclient.New(httpclient.WithBaseURL(addr))
	if err != nil {
		return nil, err
	}
	return NewSessionFromClient(client, opts)
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
func (b *Session) Url(ctx context.Context, url string) error {
	return b.navigation.NavigateTo(ctx, url)
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
func (b *Session) Session() w3c.Session {
	return b.session
}

// Timeouts returns a timeouts protocol.
func (b *Session) Timeouts() w3c.Timeouts {
	return b.timeouts
}

// Navigation returns a navigation protocol.
func (b *Session) Navigation() w3c.Navigation {
	return b.navigation
}

// Context returns a context protocol.
func (b *Session) Context() w3c.Context {
	return b.context
}

// Cookies returns a cookies protocol.
func (b *Session) Cookies() w3c.Cookies {
	return b.cookies
}

// Document returns a document protocol.
func (b *Session) Document() w3c.Document {
	return b.document
}

// ScreenCapture returns a screen capture protocol.
func (b *Session) ScreenCapture() w3c.ScreenCapture {
	return b.screenCapture
}

// Elements returns an elements protocol.
func (b *Session) Elements() w3c.Elements {
	return b.elements
}

// Close close the current session.
func (b *Session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

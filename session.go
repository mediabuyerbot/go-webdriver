package webdriver

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3cproto"
	"github.com/mediabuyerbot/httpclient"
)

type Session struct {
	session       w3cproto.Session
	timeouts      w3cproto.Timeouts
	navigation    w3cproto.Navigation
	context       w3cproto.Context
	cookies       w3cproto.Cookies
	document      w3cproto.Document
	screenCapture w3cproto.ScreenCapture
	elements      w3cproto.Elements
}

func NewSessionFromClient(ctx context.Context, client httpclient.Client, opts w3cproto.BrowserOptions) (*Session, error) {
	cli := w3cproto.WithClient(client)
	sess, err := w3cproto.NewSession(ctx, cli, opts)
	if err != nil {
		return nil, err
	}

	browser := Session{
		session:       sess,
		timeouts:      w3cproto.NewTimeouts(cli, sess.ID()),
		navigation:    w3cproto.NewNavigation(cli, sess.ID()),
		context:       w3cproto.NewContext(cli, sess.ID()),
		cookies:       w3cproto.NewCookies(cli, sess.ID()),
		document:      w3cproto.NewDocument(cli, sess.ID()),
		elements:      w3cproto.NewElements(cli, sess.ID()),
		screenCapture: w3cproto.NewScreenCapture(cli, sess.ID()),
	}
	return &browser, nil
}

func NewSession(ctx context.Context, addr string, opts w3cproto.BrowserOptions) (*Session, error) {
	client, err := httpclient.New(httpclient.WithBaseURL(addr))
	if err != nil {
		return nil, err
	}
	return NewSessionFromClient(ctx, client, opts)
}

// SessionID returns the unique session id.
func (b *Session) SessionID() string {
	return b.session.ID()
}

// ScreenshotPNG takes a screenshot of the current page.
func (b *Session) ScreenshotPNG(ctx context.Context) (image.Image, error) {
	reader, err := b.screenCapture.Take(ctx)
	if err != nil {
		return nil, err
	}
	return png.Decode(reader)
}

// ScreenshotJPG takes a screenshot of the current page.
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
func (b *Session) Session() w3cproto.Session {
	return b.session
}

// Timeouts returns a timeouts protocol.
func (b *Session) Timeouts() w3cproto.Timeouts {
	return b.timeouts
}

// Navigation returns a navigation protocol.
func (b *Session) Navigation() w3cproto.Navigation {
	return b.navigation
}

// Context returns a context protocol.
func (b *Session) Context() w3cproto.Context {
	return b.context
}

// Cookies returns a cookies protocol.
func (b *Session) Cookies() w3cproto.Cookies {
	return b.cookies
}

// Document returns a document protocol.
func (b *Session) Document() w3cproto.Document {
	return b.document
}

// ScreenCapture returns a screen capture protocol.
func (b *Session) ScreenCapture() w3cproto.ScreenCapture {
	return b.screenCapture
}

// Elements returns an elements protocol.
func (b *Session) Elements() w3cproto.Elements {
	return b.elements
}

// Close close the current session.
func (b *Session) Close(ctx context.Context) error {
	return b.session.Delete(ctx)
}

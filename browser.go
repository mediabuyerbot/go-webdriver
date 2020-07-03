package webdriver

import (
	"context"
	"image"
	"io"
	"net"
	"strings"
	"time"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3cproto"
)

type Driver interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type Browser struct {
	ctx    context.Context
	sess   *Session
	driver Driver
}

func (b *Browser) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// UID returns the unique browser id.
func (b *Browser) UID() string {
	return b.sess.Session().ID()
}

// Execute inject a snippet of JavaScript into the page for execution in the context of the currently
// selected frame. The executed script is assumed to be synchronous and the result of evaluating the script
// is returned to the client. The script argument defines the script to execute in the form of a function body.
// The value returned by that function will be returned to the client. The function will be invoked with
// the provided args array and the values may be accessed via the arguments object in the order specified.
// Arguments may be any JSON-primitive, array, or JSON object.
func (b *Browser) Execute(script string, args []interface{}) ([]byte, error) {
	return b.sess.Document().ExecuteScript(b.ctx, script, args)
}

// ExecuteAsync inject a snippet of JavaScript into the page for execution in the context of the
// currently selected frame. The executed script is assumed to be asynchronous and must signal that
// is done by invoking the provided callback, which is always provided as the final argument to the function.
// The value to this callback will be returned to the client. Asynchronous script commands may not span page loads.
// If an unload event is fired while waiting for a script result, an error should be returned to the client.
// The script argument defines the script to execute in teh form of a function body. The function will be invoked
// with the provided args array and the values may be accessed via the arguments object in the order specified.
// The final argument will always be a callback function that must be invoked to signal that the script has finished.
// Arguments may be any JSON-primitive, array, or JSON object.
func (b *Browser) ExecuteAsync(script string, args []interface{}) ([]byte, error) {
	return b.sess.Document().ExecuteAsyncScript(b.ctx, script, args)
}

// Source returns a string serialization of the DOM of the current browsing context active document.
func (b *Browser) Source() (source string, err error) {
	return b.sess.Document().GetPageSource(b.ctx)
}

// Cookies returns an all cookies visible to the current page.
func (b *Browser) Cookies() ([]w3cproto.Cookie, error) {
	return b.sess.Cookies().All(b.ctx)
}

// GetCookie returns a cookie by name visible to the current page.
func (b *Browser) GetCookie(name string) (w3cproto.Cookie, error) {
	return b.sess.Cookies().Get(b.ctx, name)
}

// AddCookie adds a cookie.
func (b *Browser) AddCookie(c w3cproto.Cookie) error {
	return b.sess.Cookies().Add(b.ctx, c)
}

// DeleteCookie deletes cookies by name visible to the current page.
func (b *Browser) DeleteCookie(name string) error {
	return b.sess.Cookies().Delete(b.ctx, name)
}

// DeleteCookies deletes all cookies visible to the current page.
func (b *Browser) DeleteCookies() error {
	return b.sess.Cookies().DeleteAll(b.ctx)
}

// NavigateTo navigates to a new URL.
func (b *Browser) NavigateTo(u string) error {
	return b.sess.Navigation().NavigateTo(b.ctx, u)
}

// Url  navigates to a new URL (alias for NavigateTo).
func (b *Browser) Url(u string) error {
	return b.sess.Navigation().NavigateTo(b.ctx, u)
}

// CurrentURL returns the URL of the current page.
func (b *Browser) CurrentURL() (url string, err error) {
	return b.sess.Navigation().GetCurrentURL(b.ctx)
}

// Back navigate backwards in the browser history, if possible.
func (b *Browser) Back() error {
	return b.sess.Navigation().Back(b.ctx)
}

// Refresh refresh the current page.
func (b *Browser) Refresh(ctx context.Context) error {
	return b.sess.Navigation().Refresh(b.ctx)
}

// Title returns the current page title.
func (b *Browser) Title(ctx context.Context) (title string, err error) {
	return b.sess.Navigation().GetTitle(b.ctx)
}

// Forward navigate forwards in the browser history, if possible.
func (b *Browser) Forward() error {
	return b.sess.Navigation().Forward(b.ctx)
}

// GetTimeout returns the timeouts implicit, pageLoad, script.
func (b *Browser) GetTimeout() (w3cproto.Timeout, error) {
	return b.sess.Timeouts().Get(b.ctx)
}

// SetImplicitTimeout sets the amount of time the browser should wait when
// searching for elements. The timeout will be rounded to nearest millisecond.
func (b *Browser) SetImplicitTimeout(d time.Duration) error {
	return b.sess.Timeouts().SetImplicit(b.ctx, d)
}

// SetPageLoadTimeout sets the amount of time the browser should wait when
// loading a page. The timeout will be rounded to nearest millisecond.
func (b *Browser) SetPageLoadTimeout(d time.Duration) error {
	return b.sess.Timeouts().SetPageLoad(b.ctx, d)
}

// SetScriptTimeout sets the amount of time that asynchronous scripts
// are permitted to run before they are aborted. The timeout will be rounded
// to nearest millisecond.
func (b *Browser) SetScriptTimeout(d time.Duration) error {
	return b.sess.Timeouts().SetScript(b.ctx, d)
}

// Capabilities returns the browser capabilities.
func (b *Browser) Capabilities() w3cproto.Capabilities {
	return b.sess.Session().Capabilities()
}

// Status returns information about whether a browser  is in a state
// in which it can create new sessions, but may additionally include arbitrary
// meta information that is specific to the implementation.
func (b *Browser) Status() (w3cproto.Status, error) {
	return b.sess.Session().Status(b.ctx)
}

// FindElementByID finds an element on the page, starting from the document root.
func (b *Browser) FindElementByID(id string) (we WebElement, err error) {
	if len(id) == 0 {
		return we, w3cproto.ErrInvalidArguments
	}
	if !strings.HasPrefix(id, "#") {
		id = "#" + id
	}
	w3cWebElem, err := b.sess.Elements().FindOne(b.ctx, w3cproto.ByCSSSelector, id)
	if err != nil {
		return we, err
	}
	return WebElement{
		elem: w3cWebElem,
		ctx:  b.ctx,
		q: selector{
			id:       id,
			strategy: w3cproto.ByCSSSelector,
		},
	}, nil
}

// FindElementByXPATH finds an element on the page, starting from the document root.
func (b *Browser) FindElementByXPATH(xpath string) (we WebElement, err error) {
	w3cWebElem, err := b.sess.Elements().FindOne(b.ctx, w3cproto.ByXPATH, xpath)
	if err != nil {
		return we, err
	}
	return WebElement{
		elem: w3cWebElem,
		ctx:  b.ctx,
		q: selector{
			id:       xpath,
			strategy: w3cproto.ByXPATH,
		},
	}, nil
}

// FindElementByLinkText finds an element on the page, starting from the document root.
func (b *Browser) FindElementByLinkText(text string) (we WebElement, err error) {
	w3cWebElem, err := b.sess.Elements().FindOne(b.ctx, w3cproto.ByLinkText, text)
	if err != nil {
		return we, err
	}
	return WebElement{
		elem: w3cWebElem,
		ctx:  b.ctx,
		q: selector{
			id:       text,
			strategy: w3cproto.ByLinkText,
		},
	}, nil
}

// Windows returns the list of all window handles(ids) available to the session.
func (b *Browser) Windows() (ids []w3cproto.WindowHandle, err error) {
	handles, err := b.sess.Context().GetWindowHandles(b.ctx)
	if err != nil {
		return nil, err
	}
	return handles, err
}

// ActiveWindow returns the ID of current window handle.
func (b *Browser) ActiveWindow() (id w3cproto.WindowHandle, err error) {
	handle, err := b.sess.Context().GetWindowHandle(b.ctx)
	if err != nil {
		return id, err
	}
	return handle, nil
}

// CloseActiveWindow closes the current window.
func (b *Browser) CloseActiveWindow() error {
	_, err := b.sess.Context().CloseWindow(b.ctx)
	return err
}

// CloseWindow closes a window.
func (b *Browser) CloseWindow(id w3cproto.WindowHandle) error {
	if id.IsEmpty() {
		return nil
	}
	if err := b.SwitchTo(id); err != nil {
		return err
	}
	return b.CloseActiveWindow()
}

// OpenTab creates a new tab.
func (b *Browser) OpenTab() (id w3cproto.WindowHandle, err error) {
	win, err := b.sess.Context().NewWindow(b.ctx, w3cproto.Tab)
	if err != nil {
		return id, err
	}
	return win.Handle, nil
}

// OpenWindow creates a new window.
func (b *Browser) OpenWindow() (id w3cproto.WindowHandle, err error) {
	win, err := b.sess.Context().NewWindow(b.ctx, w3cproto.Win)
	if err != nil {
		return id, err
	}
	return win.Handle, nil
}

// SwitchTo switches between tabs or windows.
func (b *Browser) SwitchTo(id w3cproto.WindowHandle) error {
	if id.IsEmpty() {
		return w3cproto.ErrUnknownWindowHandler
	}
	return b.sess.Context().SwitchToWindow(b.ctx, id)
}

// SwitchToFrame changes focus to another frame on the page.
func (b *Browser) SwitchToFrame(id w3cproto.FrameHandle) error {
	return b.sess.Context().SwitchToFrame(b.ctx, id)
}

// SwitchToParentFrame changes focus back to parent frame.
func (b *Browser) SwitchToParentFrame() error {
	return b.sess.Context().SwitchToParentFrame(b.ctx)
}

// ResizeWindow alters the size or position of the operating system window.
func (b *Browser) ResizeWindow(r w3cproto.Rect) (winRect w3cproto.Rect, err error) {
	return b.sess.Context().SetRect(b.ctx, r)
}

// MoveTo alters the position of the operating system window.
func (b *Browser) MoveTo(x, y int) error {
	rect, err := b.sess.Context().GetRect(b.ctx)
	if err != nil {
		return err
	}

	rect.X = x
	rect.Y = y

	if _, err := b.sess.Context().SetRect(b.ctx, rect); err != nil {
		return err
	}
	return nil
}

// ResizeTo alters the size of the operating system window.
func (b *Browser) ResizeTo(width, height int) error {
	rect, err := b.sess.Context().GetRect(b.ctx)
	if err != nil {
		return err
	}

	rect.Width = width
	rect.Height = height

	if _, err := b.sess.Context().SetRect(b.ctx, rect); err != nil {
		return err
	}
	return nil
}

// ScreenSize returns a window size on the screen of the operating system.
func (b *Browser) ScreenSize() (width int, height int, err error) {
	rect, err := b.sess.Context().GetRect(b.ctx)
	if err != nil {
		return width, height, err
	}
	return rect.Width, rect.Height, nil
}

// ScreenPosition returns a window position on the screen of the operating system.
func (b *Browser) ScreenPosition() (x int, y int, err error) {
	rect, err := b.sess.Context().GetRect(b.ctx)
	if err != nil {
		return x, y, err
	}
	return rect.X, rect.Y, nil
}

// Maximize increases the window to the maximum available size without going full-screen.
func (b *Browser) Maximize() error {
	if _, err := b.sess.Context().Maximize(b.ctx); err != nil {
		return err
	}
	return nil
}

// Minimize decreases the window to the minimum available size.
func (b *Browser) Minimize() error {
	if _, err := b.sess.Context().Minimize(b.ctx); err != nil {
		return err
	}
	return nil
}

// Fullscreen resizes the window to full screen.
func (b *Browser) Fullscreen() error {
	if _, err := b.sess.Context().Fullscreen(b.ctx); err != nil {
		return err
	}
	return nil
}

// ScreenshotJPG takes a screenshot of the current page.
func (b *Browser) ScreenshotJPG() (image.Image, error) {
	return b.sess.ScreenshotJPG(b.ctx)
}

// ScreenshotPNG takes a screenshot of the current page.
func (b *Browser) ScreenshotPNG() (image.Image, error) {
	return b.sess.ScreenshotPNG(b.ctx)
}

// Screenshot takes a screenshot of the current page.
func (b *Browser) Screenshot() (io.Reader, error) {
	return b.sess.ScreenCapture().Take(b.ctx)
}

// ElementScreenshot takes a screenshot of the element on the current page.
func (b *Browser) ElementScreenshot(elementID string) (io.Reader, error) {
	return b.sess.ScreenCapture().TakeElement(b.ctx, elementID)
}

func (b *Browser) Close() error {
	defer func() {
		if b.driver != nil {
			_ = b.driver.Stop(b.ctx)
		}
	}()
	return b.sess.Close(b.ctx)
}

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

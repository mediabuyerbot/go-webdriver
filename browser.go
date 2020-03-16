package webdriver

import (
	"context"
	"net"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
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

// Windows returns the list of all window handles(ids) available to the session.
func (b *Browser) Windows() (ids []w3c.WindowHandle, err error) {
	handles, err := b.sess.Context().GetWindowHandles(b.ctx)
	if err != nil {
		return nil, err
	}
	return handles, err
}

// ActiveWindow returns the ID of current window handle.
func (b *Browser) ActiveWindow() (id w3c.WindowHandle, err error) {
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
func (b *Browser) CloseWindow(id w3c.WindowHandle) error {
	if id.IsEmpty() {
		return nil
	}
	if err := b.SwitchTo(id); err != nil {
		return err
	}
	return b.CloseActiveWindow()
}

// OpenTab creates a new tab.
func (b *Browser) OpenTab() (id w3c.WindowHandle, err error) {
	win, err := b.sess.Context().NewWindow(b.ctx, w3c.Tab)
	if err != nil {
		return id, err
	}
	return win.Handle, nil
}

// OpenWindow creates a new window.
func (b *Browser) OpenWindow() (id w3c.WindowHandle, err error) {
	win, err := b.sess.Context().NewWindow(b.ctx, w3c.Win)
	if err != nil {
		return id, err
	}
	return win.Handle, nil
}

// SwitchTo switches between tabs or windows.
func (b *Browser) SwitchTo(id w3c.WindowHandle) error {
	if id.IsEmpty() {
		return w3c.ErrUnknownWindowHandler
	}
	return b.sess.Context().SwitchToWindow(b.ctx, id)
}

// SwitchToFrame changes focus to another frame on the page.
func (b *Browser) SwitchToFrame(id w3c.FrameHandle) error {
	return b.sess.Context().SwitchToFrame(b.ctx, id)
}

// SwitchToParentFrame changes focus back to parent frame.
func (b *Browser) SwitchToParentFrame() error {
	return b.sess.Context().SwitchToParentFrame(b.ctx)
}

// ResizeWindow alters the size or position of the operating system window.
func (b *Browser) ResizeWindow(r w3c.Rect) (winRect w3c.Rect, err error) {
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

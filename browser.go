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

// CloseActiveWindow close the current window.
func (b *Browser) CloseActiveWindow() error {
	_, err := b.sess.Context().CloseWindow(b.ctx)
	return err
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
	return b.sess.Context().SwitchToWindow(b.ctx, id)
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

package webdriver

import (
	"github.com/mediabuyerbot/go-webdriver/httpclient"
	"github.com/mediabuyerbot/go-webdriver/internal/protocol"
)

type Session struct {
	sessionID string
	client    httpclient.Client

	cookies *protocol.Cookies
	window  *protocol.Window
	element *protocol.Element
}

func (b *Session) SessionID() string {
	return b.sessionID
}

func (b *Session) Cookies() *protocol.Cookies {
	return b.cookies
}

func (b *Session) Window() *protocol.Window {
	return b.window
}

func (b *Session) Element() *protocol.Element {
	return b.element
}

func NewSession(client httpclient.Client, desired, required protocol.Capabilities) (*Session, error) {
	return nil, nil
}

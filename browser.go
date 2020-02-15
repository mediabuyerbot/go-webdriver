package webdriver

type Browser struct {
	sessionID string
	client    Client

	cookies *Cookies
	window  *Window
	element *Element
}

func (b *Browser) ID() string {
	return b.sessionID
}

func (b *Browser) Cookies() *Cookies {
	return b.cookies
}

func (b *Browser) Window() *Window {
	return b.window
}

func (b *Browser) Element() *Element {
	return b.element
}

func NewBrowser(addr string, desired, required Capabilities) (*Browser, error) {
	return &Browser{
		cookies: &Cookies{sessionID: ""},
		window:  &Window{sessionID: ""},
		element: &Element{sessionID: ""},
	}, nil
}

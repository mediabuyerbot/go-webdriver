package webdriver

type Browser struct {
	sessionID string

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

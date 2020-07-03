package w3cproto

type BrowserOptions interface {
	FirstMatch() []Capabilities
	AlwaysMatch() Capabilities
}

type browserOptions struct {
	firstMatch  []Capabilities
	alwaysMatch Capabilities
}

// NewBrowserOptions returns an instance of BrowserOptions.
func NewBrowserOptions(alwaysMatch Capabilities, firstMatch []Capabilities) BrowserOptions {
	return &browserOptions{
		firstMatch:  firstMatch,
		alwaysMatch: alwaysMatch,
	}
}

func (o *browserOptions) FirstMatch() []Capabilities {
	return o.firstMatch
}

func (o *browserOptions) AlwaysMatch() Capabilities {
	return o.alwaysMatch
}

package webdriver

type Cookie struct {
	Name   string
	Value  string
	Path   string
	Domain string
	Secure bool
	Expiry int
}

type Cookies struct {
	sessionID string
}

func (c *Cookie) All() ([]Cookie, error) {
	return nil, nil
}

func (c *Cookie) Get(name string) (Cookie, error) {
	return Cookie{}, nil
}

func (c *Cookie) Set(v Cookie) error {
	return nil
}

func (c *Cookie) DeleteAll() error {
	return nil
}

func (c *Cookie) DeleteByName(name string) error {
	return nil
}

package protocol

import "github.com/mediabuyerbot/go-webdriver/pkg/httpclient"

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
	client    httpclient.Client
}

func (c *Cookies) All() ([]Cookie, error) {
	return nil, nil
}

func (c *Cookies) Get(name string) (Cookie, error) {
	return Cookie{}, nil
}

func (c *Cookies) Set(v Cookie) error {
	return nil
}

func (c *Cookies) DeleteAll() error {
	return nil
}

func (c *Cookies) DeleteByName(name string) error {
	return nil
}

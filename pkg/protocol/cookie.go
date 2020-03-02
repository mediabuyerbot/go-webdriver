package protocol

import (
	"context"
	"net/http"
)

// Cookies represents a cookies protocol.
type Cookies interface {

	// GetAllCookies returns an all cookies visible to the current page.
	GetAllCookies(context.Context) ([]Cookie, error)

	// GetNamedCookie returns a cookie by name visible to the current page.
	GetNamedCookie(context.Context, string) (Cookie, error)

	// AddCookie add a cookie.
	AddCookie(context.Context, Cookie) error

	// DeleteCookie delete cookies by name visible to the current page.
	DeleteCookie(context.Context, string) error

	// DeleteAllCookies delete all cookies visible to the current page.
	DeleteAllCookies(context.Context) error
}

// Cookie represents a name-value pair described in [RFC6265].
type Cookie struct {
	// The name of the cookie.
	Name string `json:"name"`
	// The cookie value.
	Value string `json:"value"`
	// The cookie path. Defaults to "/" if omitted when adding a cookie.
	Path string `json:"path"`
	// The domain the cookie is visible to. Defaults to the current browsing context’s active
	// document’s URL domain if omitted when adding a cookie.
	Domain string `json:"domain"`
	// Whether the cookie is a secure cookie. Defaults to false if omitted when adding a cookie.
	Secure bool `json:"secure"`
	// When the cookie expires, specified in seconds since Unix Epoch. Must not be set if omitted when adding a cookie.
	Expiry int `json:"expiry"`
	// Whether the cookie is an HTTP only cookie. Defaults to false if omitted when adding a cookie.
	HttpOnly bool `json:"httpOnly"`
	// Whether the cookie applies to a SameSite policy. Defaults to None if omitted when adding a cookie.
	// Can be set to either "Lax" or "Strict".
	SameSite string `json:"sameSite"`
}

type cookies struct {
	request Doer
	id      string
}

// NewCookies creates a new instance of Cookie.
func NewCookies(doer Doer, sessID string) Cookies {
	return &cookies{
		id:      sessID,
		request: doer,
	}
}

func (c *cookies) GetAllCookies(ctx context.Context) (cookies []Cookie, err error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/cookie", nil)
	if err != nil {
		return nil, err
	}
	resp = resp
	return nil, nil
}

func (c *cookies) GetNamedCookie(ctx context.Context, name string) (cookie Cookie, err error) {
	return Cookie{}, nil
}

func (c *cookies) AddCookie(ctx context.Context, cookie Cookie) error {
	return nil
}

func (c *cookies) DeleteCookie(ctx context.Context, name string) error {
	return nil
}

func (c *cookies) DeleteAllCookies(ctx context.Context) error {
	return nil
}

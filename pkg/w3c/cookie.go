package w3c

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	// CookieNameKey the name of the cookie.
	CookieNameKey = "name"
	// CookieValueKey the cookie value.
	CookieValueKey = "value"
	// CookiePathKey  the cookie path. Defaults to "/" if omitted when adding a cookie.
	CookiePathKey = "path"
	// CookieDomainKey the domain the cookie is visible to. Defaults to the current browsing context’s active
	// document’s URL domain if omitted when adding a cookie.
	CookieDomainKey = "domain"
	// CookieSecureKey  whether the cookie is a secure cookie. Defaults to false if omitted when adding a cookie.
	CookieSecureKey = "secure"
	// CookieExpiryKey when the cookie expires, specified in seconds since Unix Epoch. Must not be set if omitted when adding a cookie.
	CookieExpiryKey = "expiry"
	// CookieHttpOnlyKey whether the cookie is an HTTP only cookie. Defaults to false if omitted when adding a cookie.
	CookieHttpOnlyKey = "httpOnly"
	// CookieSameSiteKey whether the cookie applies to a SameSite policy. Defaults to None if omitted when adding a cookie.
	// Can be set to either "Lax" or "Strict".
	CookieSameSiteKey = "sameSite"
)

var (
	ErrInvalidCookie = errors.New("protocol: cookie name and value is required")
)

// Cookies represents a cookies protocol.
type Cookies interface {

	// All returns an all cookies visible to the current page.
	All(context.Context) ([]Cookie, error)

	// Get returns a cookie by name visible to the current page.
	Get(context.Context, string) (Cookie, error)

	// Add adds a cookie.
	Add(context.Context, Cookie) error

	// Delete deletes cookies by name visible to the current page.
	Delete(context.Context, string) error

	// DeleteAll deletes all cookies visible to the current page.
	DeleteAll(context.Context) error
}

type Cookie map[string]interface{}

func MakeCookie() Cookie {
	return make(Cookie)
}

func (c Cookie) SetName(v string) Cookie {
	c[CookieNameKey] = v
	return c
}

func (c Cookie) Name() string {
	v, ok := c[CookieNameKey]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

func (c Cookie) SetValue(v interface{}) Cookie {
	c[CookieValueKey] = v
	return c
}

func (c Cookie) Value() interface{} {
	v, ok := c[CookieValueKey]
	if !ok {
		return nil
	}
	return v
}

func (c Cookie) SetPath(v string) Cookie {
	c[CookiePathKey] = v
	return c
}

func (c Cookie) Path() string {
	v, ok := c[CookiePathKey]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

func (c Cookie) SetDomain(v string) Cookie {
	c[CookieDomainKey] = v
	return c
}

func (c Cookie) Domain() string {
	v, ok := c[CookieDomainKey]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

func (c Cookie) SetSecure(v bool) Cookie {
	c[CookieSecureKey] = v
	return c
}

func (c Cookie) Secure() bool {
	v, ok := c[CookieSecureKey]
	if !ok {
		return false
	}
	n, ok := v.(bool)
	if !ok {
		return false
	}
	return n
}

func (c Cookie) SetExpiry(v int64) Cookie {
	c[CookieExpiryKey] = v
	return c
}

func (c Cookie) Expiry() int64 {
	v, ok := c[CookieExpiryKey]
	if !ok {
		return 0
	}
	n, ok := v.(int64)
	if !ok {
		return 0
	}
	return n
}

func (c Cookie) SetHttpOnly(v bool) Cookie {
	c[CookieHttpOnlyKey] = v
	return c
}

func (c Cookie) HttpOnly() bool {
	v, ok := c[CookieHttpOnlyKey]
	if !ok {
		return false
	}
	n, ok := v.(bool)
	if !ok {
		return false
	}
	return n
}

func (c Cookie) SetSameSite(v string) Cookie {
	c[CookieSameSiteKey] = v
	return c
}

func (c Cookie) SameSite() string {
	v, ok := c[CookieSameSiteKey]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

func (c Cookie) Get(key string) interface{} {
	v, ok := c[key]
	if !ok {
		return nil
	}
	return v
}

func (c Cookie) Set(key string, value interface{}) Cookie {
	c[key] = value
	return c
}

func (c Cookie) Delete(key string) {
	delete(c, key)
}

func (c Cookie) ToParams() Params {
	p := make(Params)
	for k, v := range c {
		p.Set(k, v)
	}
	return p
}

func (c Cookie) Validate() error {
	_, ok := c[CookieNameKey]
	if !ok {
		return ErrInvalidCookie
	}
	_, ok = c[CookieValueKey]
	if !ok {
		return ErrInvalidCookie
	}
	return nil
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

func (c *cookies) All(ctx context.Context) (cookies []Cookie, err error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/cookie", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resp.Value, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

func (c *cookies) Get(ctx context.Context, name string) (cookie Cookie, err error) {
	resp, err := c.request.Do(ctx, http.MethodGet, "/session/"+c.id+"/cookie/"+name, nil)
	if err != nil {
		return cookie, err
	}
	if err := json.Unmarshal(resp.Value, &cookie); err != nil {
		return cookie, err
	}
	return cookie, nil
}

func (c *cookies) Add(ctx context.Context, cookie Cookie) error {
	if err := cookie.Validate(); err != nil {
		return err
	}
	p := Params{
		"cookie": cookie.ToParams(),
	}
	resp, err := c.request.Do(ctx, http.MethodPost, "/session/"+c.id+"/cookie", p)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *cookies) Delete(ctx context.Context, name string) error {
	resp, err := c.request.Do(ctx, http.MethodDelete, "/session/"+c.id+"/cookie/"+name, nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (c *cookies) DeleteAll(ctx context.Context) error {
	resp, err := c.request.Do(ctx, http.MethodDelete, "/session/"+c.id+"/cookie", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

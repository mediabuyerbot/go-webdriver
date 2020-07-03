package w3cproto

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	CapabilityBrowserName               = "browserName"
	CapabilityBrowserVersion            = "browserVersion"
	CapabilityPlatformName              = "platformName"
	CapabilityAcceptInsecureCerts       = "acceptInsecureCerts"
	CapabilityPageLoadStrategy          = "pageLoadStrategy"
	CapabilityProxy                     = "proxy"
	CapabilitySetWindowRect             = "setWindowRect"
	CapabilityTimeouts                  = "timeouts"
	CapabilityUnhandledPromptBehavior   = "unhandledPromptBehavior"
	CapabilityStrictFileInteractability = "strictFileInteractability"
)

const (
	Linux   Platform = "linux"
	Mac     Platform = "mac"
	Windows Platform = "windows"
)

type Platform string

func (p Platform) Validate() error {
	switch p {
	case Linux, Mac, Windows:
		return nil
	default:
		return fmt.Errorf("unknown platform %s", p)
	}
}

func (p Platform) String() string {
	return string(p)
}

type Capabilities map[string]interface{}

func MakeCapabilities() Capabilities {
	return make(Capabilities)
}

func SetBrowserName(c Capabilities, name string) error {
	c[CapabilityBrowserName] = name
	return nil
}

func GetBrowserName(c Capabilities) string {
	return c.GetString(CapabilityBrowserName)
}

func SetBrowserVersion(c Capabilities, v string) error {
	c[CapabilityBrowserVersion] = v
	return nil
}

func GetBrowserVersion(c Capabilities) string {
	return c.GetString(CapabilityBrowserVersion)
}

func SetPlatformName(c Capabilities, platform Platform) error {
	c[CapabilityPlatformName] = platform.String()
	return nil
}

func GetPlatformName(c Capabilities) string {
	return c.GetString(CapabilityPlatformName)
}

func SetAcceptInsecureCerts(c Capabilities, flag bool) error {
	c[CapabilityAcceptInsecureCerts] = flag
	return nil
}

func GetAcceptInsecureCerts(c Capabilities) bool {
	return c.GetBool(CapabilityAcceptInsecureCerts)
}

func SetPageLoadStrategy(c Capabilities, strategy string) error {
	c[CapabilityPageLoadStrategy] = strategy
	return nil
}

func SetWindowRect(c Capabilities, flag bool) error {
	c[CapabilitySetWindowRect] = flag
	return nil
}

func GetWindowRect(c Capabilities) bool {
	return c.GetBool(CapabilitySetWindowRect)
}

func SetUnhandledPromptBehavior(c Capabilities, prompt string) error {
	c[CapabilityUnhandledPromptBehavior] = prompt
	return nil
}

func GetUnhandledPromptBehavior(c Capabilities) string {
	return c.GetString(CapabilityUnhandledPromptBehavior)
}

func SetStrictFileInteractability(c Capabilities, flag bool) error {
	c[CapabilityStrictFileInteractability] = flag
	return nil
}

func GetStrictFileInteractability(c Capabilities) bool {
	return c.GetBool(CapabilityStrictFileInteractability)
}

func SetProxy(c Capabilities, p *Proxy) error {
	proxyCap := MakeCapabilities()
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &proxyCap); err != nil {
		return err
	}
	c[CapabilityProxy] = proxyCap
	return nil
}

func GetProxy(c Capabilities) *Proxy {
	if !c.Has(CapabilityProxy) {
		return nil
	}
	b, err := json.Marshal(c[CapabilityProxy])
	if err != nil {
		return nil
	}
	proxy := &Proxy{}
	if err := json.Unmarshal(b, proxy); err != nil {
		return nil
	}
	return proxy
}

func SetTimeout(c Capabilities, t Timeout) error {
	timeoutCap := MakeCapabilities()
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &timeoutCap); err != nil {
		return err
	}
	c[CapabilityTimeouts] = timeoutCap
	return nil
}

func GetTimeout(c Capabilities) *Timeout {
	if !c.Has(CapabilityTimeouts) {
		return nil
	}
	b, err := json.Marshal(c[CapabilityTimeouts])
	if err != nil {
		return nil
	}
	t := &Timeout{}
	if err := json.Unmarshal(b, t); err != nil {
		return nil
	}
	return t
}

func (c Capabilities) Section(key string) Capabilities {
	v, ok := c[key]
	if !ok {
		return nil
	}
	switch x := v.(type) {
	case map[string]interface{}:
		return Capabilities(x)
	case Capabilities:
		return x
	default:
		return nil
	}
}

func (c Capabilities) Has(key string) bool {
	_, ok := c[key]
	return ok
}

func (c Capabilities) Set(k string, v interface{}) Capabilities {
	c[k] = v
	return c
}

func (c Capabilities) GetString(key string) (s string) {
	v, ok := c[key]
	if !ok {
		return
	}
	switch x := v.(type) {
	case []byte:
		return string(x)
	case string:
		return x
	default:
		s, ok = v.(string)
		if !ok {
			return
		}
		return s
	}
}

func (c Capabilities) GetStringSlice(key string) (s []string) {
	v, ok := c[key]
	if !ok {
		return nil
	}
	s, err := toSlice(v)
	if err != nil {
		return nil
	}
	return s
}

func (c Capabilities) GetBool(key string) (b bool) {
	v, ok := c[key]
	if !ok {
		return
	}
	b, ok = v.(bool)
	if !ok {
		return
	}
	return b
}

func (c Capabilities) GetUint(key string) (i uint) {
	v, ok := c[key]
	if !ok {
		return
	}
	switch x := v.(type) {
	case float64:
		return uint(x)
	case uint:
		return x
	case int:
		return uint(x)
	default:
		i, ok = v.(uint)
		if !ok {
			return
		}
		return i
	}
}

func (c Capabilities) GetInt(key string) (i int) {
	v, ok := c[key]
	if !ok {
		return
	}
	switch x := v.(type) {
	case float64:
		return int(x)
	case uint:
		return int(x)
	case int:
		return x
	default:
		f, ok := v.(float64)
		if !ok {
			return
		}
		return int(f)
	}
}

func (c Capabilities) GetFloat(key string) (f float64) {
	v, ok := c[key]
	if !ok {
		return
	}
	switch x := v.(type) {
	case uint:
		return float64(x)
	case int:
		return float64(x)
	case float64:
		return x
	default:
		f, ok = v.(float64)
		if !ok {
			return
		}
		return f
	}
}

func toSlice(slice interface{}) ([]string, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("given a non-slice type")
	}
	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		str, ok := v.(string)
		if !ok {
			return nil, errors.New("given a non-string type")
		}
		ret[i] = str
	}
	return ret, nil
}

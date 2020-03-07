package protocol

const (
	CapabilityBrowserName               = "browserName"
	CapabilityBrowserVersion            = "browserVersion"
	CapabilityPlatformName              = "platformName"
	CapabilityPlatformVersion           = "platformVersion"
	CapabilityAcceptInsecureCerts       = "acceptInsecureCerts"
	CapabilityPageLoadStrategy          = "pageLoadStrategy"
	CapabilityProxy                     = "proxy"
	CapabilitySetWindowRect             = "setWindowRect"
	CapabilityTimeouts                  = "timeouts"
	CapabilityUnhandledPromptBehavior   = "unhandledPromptBehavior"
	CapabilityStrictFileInteractability = "strictFileInteractability"
)

type Capabilities map[string]interface{}

func (o Capabilities) GetString(key string) (s string) {
	v, ok := o[key]
	if !ok {
		return
	}
	s, ok = v.(string)
	if !ok {
		return
	}
	return s
}

func (o Capabilities) GetStringSlice(key string) (s []string) {
	v, ok := o[key]
	if !ok {
		return
	}
	s, ok = v.([]string)
	if !ok {
		return
	}
	return s
}

func (o Capabilities) GetBool(key string) (b bool) {
	v, ok := o[key]
	if !ok {
		return
	}
	b, ok = v.(bool)
	if !ok {
		return
	}
	return b
}

func (o Capabilities) GetInt(key string) (i int) {
	v, ok := o[key]
	if !ok {
		return
	}
	t, ok := v.(float64)
	if !ok {
		return
	}
	return int(t)
}

func (o Capabilities) GetFloat(key string) (f float64) {
	v, ok := o[key]
	if !ok {
		return
	}
	f, ok = v.(float64)
	if !ok {
		return
	}
	return f
}

func (o Capabilities) BrowserName() string {
	return o.GetString(CapabilityBrowserName)
}

func (o Capabilities) BrowserVersion() string {
	return o.GetString(CapabilityBrowserVersion)
}

func (o Capabilities) PlatformName() string {
	return o.GetString(CapabilityPlatformName)
}

func (o Capabilities) PlatformVersion() string {
	return o.GetString(CapabilityPlatformVersion)
}

func (o Capabilities) AcceptInsecureCerts() bool {
	return o.GetBool(CapabilityAcceptInsecureCerts)
}

func (o Capabilities) PageLoadStrategy() string {
	return o.GetString(CapabilityPageLoadStrategy)
}

func (o Capabilities) WindowRect() bool {
	return o.GetBool(CapabilitySetWindowRect)
}

func (o Capabilities) Timeouts() (t Timeout) {
	v, ok := o[CapabilityTimeouts]
	if ok {
		doc, ok := v.(map[string]interface{})
		if !ok {
			return t
		}
		script, ok := doc["script"].(float64)
		if ok {
			t.Script = uint(script)
		}
		implicit, ok := doc["implicit"].(float64)
		if ok {
			t.Implicit = uint(implicit)
		}
		pageLoad, ok := doc["pageLoad"].(float64)
		if ok {
			t.PageLoad = uint(pageLoad)
		}
	}
	return t
}

func (o Capabilities) UnhandledPromptBehavior() string {
	return o.GetString(CapabilityUnhandledPromptBehavior)
}

func (o Capabilities) StrictFileInteractability() bool {
	return o.GetBool(CapabilityStrictFileInteractability)
}

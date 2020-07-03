package w3cproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetProxy(t *testing.T) {
	cap := MakeCapabilities()
	err := SetProxy(cap, &Proxy{
		Type:          ProxyAutodetectType,
		AutoconfigURL: "url",
		FTP:           "ftp",
		HTTP:          "http",
		SSL:           "ssl",
		SOCKS:         "socks",
		SOCKSUsername: "username",
		SOCKSPassword: "password",
		NoProxy: []string{
			"noproxy1",
			"noproxy2",
		},
		HTTPPort:  9090,
		SSLPort:   9091,
		SocksPort: 8000,
	})
	assert.Nil(t, err)
	proxy := cap.Section(CapabilityProxy)
	assert.NotNil(t, proxy)
	assert.Equal(t, string(ProxyAutodetectType), proxy.GetString("proxyType"))
	assert.Equal(t, 9090, proxy.GetInt("httpProxyPort"))
	assert.Equal(t, "url", proxy.GetString("proxyAutoconfigUrl"))
	assert.Equal(t, "ftp", proxy.GetString("ftpProxy"))
	assert.Equal(t, "http", proxy.GetString("httpProxy"))
	assert.Equal(t, "ssl", proxy.GetString("sslProxy"))
	assert.Equal(t, "socks", proxy.GetString("socksProxy"))
	assert.Equal(t, "username", proxy.GetString("socksUsername"))
	assert.Equal(t, "password", proxy.GetString("socksPassword"))
	assert.Len(t, proxy.GetStringSlice("noProxy"), 2)
	assert.Equal(t, 9091, proxy.GetInt("sslProxyPort"))
	assert.Equal(t, 8000, proxy.GetInt("socksProxyPort"))
}

func TestGetProxy(t *testing.T) {
	cap := MakeCapabilities()
	cap["proxy"] = map[string]interface{}{
		"proxyType":          ProxyAutodetectType,
		"proxyAutoconfigUrl": "proxyAutoconfigUrl",
		"ftpProxy":           "ftpProxy",
		"httpProxy":          "httpProxy",
		"sslProxy":           "sslProxy",
		"socksProxy":         "socksProxy",
		"socksUsername":      "socksUsername",
		"socksPassword":      "socksPassword",
		"noProxy":            []string{"noProxy"},
		"httpProxyPort":      9090,
		"sslProxyPort":       9090,
		"socksProxyPort":     9090,
	}
	proxy := GetProxy(cap)
	assert.NotNil(t, proxy)

	assert.Equal(t, ProxyAutodetectType, proxy.Type)
	assert.Equal(t, "proxyAutoconfigUrl", proxy.AutoconfigURL)
	assert.Equal(t, "ftpProxy", proxy.FTP)
	assert.Equal(t, "httpProxy", proxy.HTTP)
	assert.Equal(t, "sslProxy", proxy.SSL)
	assert.Equal(t, "socksProxy", proxy.SOCKS)
	assert.Equal(t, "socksUsername", proxy.SOCKSUsername)
	assert.Equal(t, "socksPassword", proxy.SOCKSPassword)
	assert.Len(t, proxy.NoProxy, 1)
	assert.Equal(t, 9090, proxy.HTTPPort)
	assert.Equal(t, 9090, proxy.SSLPort)
	assert.Equal(t, 9090, proxy.SocksPort)

	cap = MakeCapabilities()
	proxy = GetProxy(cap)
	assert.Nil(t, proxy)
}

func TestSetTimeout(t *testing.T) {
	cap := MakeCapabilities()
	val := uint(1000)
	err := SetTimeout(cap, Timeout{
		PageLoad: val,
		Script:   val,
		Implicit: val,
	})
	assert.Nil(t, err)
	timeout := cap.Section(CapabilityTimeouts)
	assert.Equal(t, val, timeout.GetUint(pageLoadTimeout))
	assert.Equal(t, val, timeout.GetUint(scriptTimeout))
	assert.Equal(t, val, timeout.GetUint(implicitTimeout))
}

func TestGetTimeout(t *testing.T) {
	cap := MakeCapabilities()
	val := uint(1000)
	cap[CapabilityTimeouts] = map[string]interface{}{
		pageLoadTimeout: val,
		scriptTimeout:   val,
		implicitTimeout: val,
	}
	timeout := GetTimeout(cap)
	assert.NotNil(t, timeout)
	assert.Equal(t, val, timeout.PageLoad)
	assert.Equal(t, val, timeout.Script)
	assert.Equal(t, val, timeout.Implicit)
}

func TestGetSetTimeout(t *testing.T) {
	cap := MakeCapabilities()
	val := uint(1000)
	t1 := Timeout{
		PageLoad: val,
		Script:   val,
		Implicit: val,
	}
	err := SetTimeout(cap, t1)
	assert.Nil(t, err)
	t2 := GetTimeout(cap)
	assert.NotNil(t, t2)
	assert.Equal(t, t1, *t2)
}

func TestCapabilities_GetBool(t *testing.T) {
	cap := MakeCapabilities()
	cap["flag"] = struct{}{}
	assert.False(t, cap.GetBool("flag"))
	cap["flag"] = ""
	assert.False(t, cap.GetBool("flag"))
	cap["flag"] = true
	assert.True(t, cap.GetBool("flag"))
	cap["flag"] = 1
	assert.False(t, cap.GetBool("flag"))
	cap["flag"] = interface{}(true)
	assert.True(t, cap.GetBool("flag"))
}

func TestCapabilities_GetFloat(t *testing.T) {
	cap := MakeCapabilities()
	cap["float"] = float64(9999)
	assert.Equal(t, cap["float"], cap.GetFloat("float"))
	cap["float"] = "9999"
	assert.Equal(t, float64(0), cap.GetFloat("float"))
	cap["float"] = 9999
	assert.Equal(t, float64(9999), cap.GetFloat("float"))
	cap["float"] = uint(9999)
	assert.Equal(t, float64(9999), cap.GetFloat("float"))
	cap["float"] = interface{}(9999)
	assert.Equal(t, float64(9999), cap.GetFloat("float"))
}

func TestCapabilities_GetInt(t *testing.T) {
	cap := MakeCapabilities()
	cap["int"] = 9999
	assert.Equal(t, cap["int"], cap.GetInt("int"))
	cap["int"] = uint(9999)
	assert.Equal(t, 9999, cap.GetInt("int"))
	cap["int"] = float64(9999)
	assert.Equal(t, 9999, cap.GetInt("int"))
	cap["int"] = uint32(9999)
	assert.Equal(t, 0, cap.GetInt("int"))
	cap["int"] = ""
	assert.Equal(t, 0, cap.GetInt("int"))
	cap["int"] = interface{}(9999)
	assert.Equal(t, 9999, cap.GetInt("int"))
}

func TestCapabilities_GetString(t *testing.T) {
	cap := MakeCapabilities()
	cap["string"] = "string"
	assert.Equal(t, cap["string"], cap.GetString("string"))
	cap["string"] = []byte("string")
	assert.Equal(t, "string", cap.GetString("string"))
	cap["string"] = struct{}{}
	assert.Equal(t, "", cap.GetString("string"))
	cap["string"] = Platform("string")
	assert.Equal(t, "", cap.GetString("string"))
	cap["string"] = interface{}("string")
	assert.Equal(t, "string", cap.GetString("string"))
}

func TestCapabilities_GetStringSlice(t *testing.T) {
	cap := MakeCapabilities()
	want := []string{"slice", "slice"}
	cap["slice"] = []string{"slice", "slice"}
	assert.Equal(t, cap["slice"], cap.GetStringSlice("slice"))
	cap["slice"] = []interface{}{"slice", "slice"}
	assert.Equal(t, want, cap.GetStringSlice("slice"))
	cap["slice"] = []interface{}{0, 0}
	assert.Nil(t, cap.GetStringSlice("slice"))
}

func TestCapabilities_GetUint(t *testing.T) {
	cap := MakeCapabilities()
	cap["uint"] = uint(9999)
	assert.Equal(t, cap["uint"], cap.GetUint("uint"))
	cap["uint"] = 9999
	assert.Equal(t, uint(9999), cap.GetUint("uint"))
	cap["uint"] = float64(9999)
	assert.Equal(t, uint(9999), cap.GetUint("uint"))
	cap["uint"] = interface{}(9999)
	assert.Equal(t, uint(9999), cap.GetUint("uint"))
}

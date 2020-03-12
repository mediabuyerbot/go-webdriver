package webdriver

import (
	"os"
	"testing"

	"github.com/mediabuyerbot/go-webdriver/pkg/protocol"

	"github.com/stretchr/testify/assert"
)

func TestChromeOptions(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.extensions)
	assert.Nil(t, builder.args)
	assert.Nil(t, builder.mobileEmulation)
	assert.Nil(t, builder.localState)
	assert.Nil(t, builder.perfLoggingPrefs)

	assert.Empty(t, builder.alwaysMatch)
	assert.Empty(t, builder.excludeSwitches)
	assert.Empty(t, builder.windowTypes)
	assert.Empty(t, builder.prefs)
	assert.Empty(t, builder.firstMatch)
}

func TestChromeOptionsBuilder_AddExcludeSwitches(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Empty(t, builder.excludeSwitches)
	builder.AddExcludeSwitches("param1", "param2", "param3")
	assert.Len(t, builder.excludeSwitches, 3)

	builder = ChromeCapabilities()
	builder.AddExcludeSwitches("", "", "")
	assert.Len(t, builder.excludeSwitches, 0)

	builder = ChromeCapabilities()
	builder.AddExcludeSwitches("param1")
	assert.Len(t, builder.excludeSwitches, 1)
	opts := builder.Build()
	assert.NotNil(t, opts.AlwaysMatch()[ChromeCapabilityExcludeSwitchesName])
	assert.Len(t, opts.AlwaysMatch()[ChromeCapabilityExcludeSwitchesName], 1)
}

func TestChromeOptionsBuilder_AddFirstMatch(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Empty(t, builder.firstMatch)
	builder.AddFirstMatch(protocol.CapabilityBrowserName, "chrome")
	builder.AddFirstMatch(protocol.CapabilityBrowserVersion, "123")
	assert.Len(t, builder.firstMatch, 2)
	opts := builder.Build()
	assert.Len(t, opts.FirstMatch(), 2)

	builder = ChromeCapabilities()
	builder.AddFirstMatch("", "value")
	builder.AddFirstMatch("", "")
	builder.AddFirstMatch("", "")
	opts = builder.Build()
	assert.Len(t, builder.firstMatch, 0)
	assert.Len(t, opts.FirstMatch(), 0)
}

func TestChromeOptionsBuilder_AddWindowTypes(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Empty(t, builder.windowTypes)
	builder.AddWindowTypes("iframe", "tab", "other")
	assert.Len(t, builder.windowTypes, 3)
	opts := builder.Build()
	assert.NotNil(t, opts.AlwaysMatch()[ChromeCapabilityWindowTypesName])
	assert.Len(t, opts.AlwaysMatch()[ChromeCapabilityWindowTypesName], 3)

	builder = ChromeCapabilities()
	builder.AddWindowTypes("", "", "")
	assert.Len(t, builder.windowTypes, 0)
	opts = builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityWindowTypesName])
}

func TestChromeOptionsBuilder_SetWindowRect(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilitySetWindowRect))
	builder.SetWindowRect(true)
	assert.True(t, builder.alwaysMatch.GetBool(protocol.CapabilitySetWindowRect))
	opts := builder.Build()
	assert.NotNil(t, opts.AlwaysMatch()[protocol.CapabilitySetWindowRect])
}

func TestChromeOptionsBuilder_Arguments(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.args)
	builder.Arguments()
	assert.NotNil(t, builder.args)
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityArgsName])

	builder = ChromeCapabilities()
	assert.Nil(t, builder.args)
	builder.Arguments().
		WithHeadless().
		WithStartMaximized().
		WithUserDataPath("/path/to/profile").
		Add("--ui-disabled-gpu")
	opts = builder.Build()
	assert.NotNil(t, builder.args)
	assert.Len(t, builder.args.opts, 4)
	assert.NotNil(t, opts.AlwaysMatch()[ChromeCapabilityArgsName])
	assert.Len(t, opts.AlwaysMatch()[ChromeCapabilityArgsName], 4)
}

func TestChromeOptionsBuilder_Extensions(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.extensions)
	builder.Extensions()
	assert.NotNil(t, builder.extensions)
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityExtensionName])

	builder = ChromeCapabilities()
	err := builder.Extensions().Add(
		// load unpacked extension
		"./testdata/chrome/extension",
		// load zip extension
		"./testdata/chrome/extension.zip",
		// load crx extension
		"./testdata/chrome/extension.crx",
	)
	assert.Nil(t, err)
	assert.NotNil(t, builder.extensions.opts)
	assert.Len(t, builder.extensions.opts, 3)
	assert.FileExists(t, "./testdata/chrome/extension.crx.pem")
	assert.Nil(t, os.Remove("./testdata/chrome/extension.crx.pem"))
	for _, b64 := range builder.extensions.opts {
		assert.NotZero(t, b64)
	}
	opts = builder.Build()
	assert.NotNil(t, opts.AlwaysMatch()[ChromeCapabilityExtensionName])

	builder = ChromeCapabilities()
	err = builder.Extensions().Add("/dev/null")
	assert.Error(t, err)
}

func TestChromeOptionsBuilder_LocalState(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.localState)
	builder.LocalState()
	assert.NotNil(t, builder.localState)
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityLocalStateName])

	builder = ChromeCapabilities()
	builder.LocalState().Set("key", "value")
	assert.NotNil(t, builder.localState)
	assert.Len(t, builder.localState.opts, 1)
	opts = builder.Build()
	assert.NotNil(t, opts.AlwaysMatch()[ChromeCapabilityLocalStateName])
	assert.Len(t, opts.AlwaysMatch()[ChromeCapabilityLocalStateName], 1)
}

func TestChromeOptionsBuilder_MobileEmulation(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.mobileEmulation)
	builder.MobileEmulation()
	assert.NotNil(t, builder.mobileEmulation)
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityMobileEmulationName])

	builder = ChromeCapabilities()
	dm := &DeviceMetrics{
		Width:      2000,
		Height:     4000,
		PixelRatio: 5,
		Touch:      true,
	}
	builder.MobileEmulation().
		SetDeviceMetrics(dm).
		SetUserAgent("user-agent blblblb").
		SetDeviceName("console").
		Set("key", "value")
	opts = builder.Build()
	alwaysMatch := opts.AlwaysMatch()
	emul := alwaysMatch.GetOpts(ChromeCapabilityMobileEmulationName)
	assert.NotNil(t, alwaysMatch[ChromeCapabilityMobileEmulationName])
	assert.Equal(t, emul.GetString("userAgent"), "user-agent blblblb")
	assert.Equal(t, emul.GetString("deviceName"), "console")
	assert.Equal(t, emul.GetString("key"), "value")
	metric := emul.GetOpts("deviceMetrics")
	assert.Equal(t, dm.Width, metric.GetUint("width"))
	assert.Equal(t, dm.Height, metric.GetUint("height"))
	assert.Equal(t, dm.PixelRatio, metric.GetFloat("pixelRatio"))
	assert.True(t, metric.GetBool("touch"))

	builder.MobileEmulation().SetDeviceMetrics(nil)
}

func TestChromeOptionsBuilder_PerfLoggingPreferences(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.perfLoggingPrefs)
	builder.PerfLoggingPreferences()
	assert.NotNil(t, builder.perfLoggingPrefs)
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityPerfLoggingPrefsName])

	builder = ChromeCapabilities()
	builder.PerfLoggingPreferences().
		TracingCategories("trace").
		BufferUsageReportingIntervalMillis(1).
		EnableNetwork(true).
		EnablePage(true).
		EnableTimeline(true).
		Set("key", "value")
	opts = builder.Build()
	assert.NotNil(t, builder.perfLoggingPrefs)
	alwaysMatch := opts.AlwaysMatch()
	prefLogs := alwaysMatch.GetOpts(ChromeCapabilityPerfLoggingPrefsName)
	assert.NotNil(t, prefLogs)
	assert.True(t, prefLogs.GetBool(perfLoggingEnableNetwork))
	assert.Equal(t, "trace", prefLogs.GetString(perfLoggingTracingCategories))
	assert.Equal(t, prefLogs.GetUint(perfLoggingBufferUsageReportingInterval), uint(1))
	assert.True(t, prefLogs.GetBool(perfLoggingEnablePage))
	assert.True(t, prefLogs.GetBool(perfLoggingEnableTimeline))
	assert.Equal(t, "value", prefLogs.GetString("key"))
}

func TestChromeOptionsBuilder_SetAcceptInsecureCerts(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityAcceptInsecureCerts))
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[protocol.CapabilityAcceptInsecureCerts])

	builder = ChromeCapabilities()
	builder.SetAcceptInsecureCerts(true)
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityAcceptInsecureCerts))
	assert.True(t, opts.AlwaysMatch().GetBool(protocol.CapabilityAcceptInsecureCerts))
}

func TestChromeOptionsBuilder_SetBinary(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(ChromeCapabilityBinaryName))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(ChromeCapabilityBinaryName))

	builder = ChromeCapabilities()
	builder.SetBinary("/path/to/chrome")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(ChromeCapabilityBinaryName))
	assert.Equal(t, "/path/to/chrome", opts.AlwaysMatch().GetString(ChromeCapabilityBinaryName))
}

func TestChromeOptionsBuilder_SetBrowserName(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityBrowserName))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityBrowserName))

	builder = ChromeCapabilities()
	builder.SetBrowserName("chrome")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityBrowserName))
	assert.Equal(t, "chrome", opts.AlwaysMatch().GetString(protocol.CapabilityBrowserName))
}

func TestChromeOptionsBuilder_SetBrowserVersion(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityBrowserVersion))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityBrowserVersion))

	builder = ChromeCapabilities()
	builder.SetBrowserVersion("chrome")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityBrowserVersion))
	assert.Equal(t, "chrome", opts.AlwaysMatch().GetString(protocol.CapabilityBrowserVersion))
}

func TestChromeOptionsBuilder_SetDebuggerAddr(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(ChromeCapabilityDebuggerAddressName))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(ChromeCapabilityDebuggerAddressName))

	builder = ChromeCapabilities()
	builder.SetDebuggerAddr("127.0.0.1:8784")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(ChromeCapabilityDebuggerAddressName))
	assert.Equal(t, "127.0.0.1:8784", opts.AlwaysMatch().GetString(ChromeCapabilityDebuggerAddressName))
}

func TestChromeOptionsBuilder_SetDetach(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(ChromeCapabilityDetachName))
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[ChromeCapabilityDetachName])

	builder = ChromeCapabilities()
	builder.SetDetach(true)
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(ChromeCapabilityDetachName))
	assert.True(t, opts.AlwaysMatch().GetBool(ChromeCapabilityDetachName))
}

func TestChromeOptionsBuilder_SetMiniDumpPath(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(ChromeCapabilityMiniDumpPathName))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(ChromeCapabilityMiniDumpPathName))

	builder = ChromeCapabilities()
	builder.SetMiniDumpPath("/path/to/dump")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(ChromeCapabilityMiniDumpPathName))
	assert.Equal(t, "/path/to/dump", opts.AlwaysMatch().GetString(ChromeCapabilityMiniDumpPathName))
}

func TestChromeOptionsBuilder_SetPageLoadStrategy(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityPageLoadStrategy))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityPageLoadStrategy))

	builder = ChromeCapabilities()
	builder.SetPageLoadStrategy("normal")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityPageLoadStrategy))
	assert.Equal(t, "normal", opts.AlwaysMatch().GetString(protocol.CapabilityPageLoadStrategy))
}

func TestChromeOptionsBuilder_SetPlatformName(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityPlatformName))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityPlatformName))

	builder = ChromeCapabilities()
	builder.SetPlatformName("platform")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityPlatformName))
	assert.Equal(t, "platform", opts.AlwaysMatch().GetString(protocol.CapabilityPlatformName))
}

func TestChromeOptionsBuilder_SetPlatformVersion(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityPlatformVersion))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityPlatformVersion))

	builder = ChromeCapabilities()
	builder.SetPlatformVersion("version")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityPlatformVersion))
	assert.Equal(t, "version", opts.AlwaysMatch().GetString(protocol.CapabilityPlatformVersion))
}

func TestChromeOptionsBuilder_SetProxy(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityProxy))
	opts := builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityProxy))

	builder = ChromeCapabilities()
	proxy := &protocol.Proxy{
		FTP:  "type",
		HTTP: "http",
	}
	builder.SetProxy(proxy)
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().HasNot(protocol.CapabilityProxy))
	assert.Equal(t, proxy, opts.Proxy())
}

func TestChromeOptionsBuilder_SetUnhandledPromptBehavior(t *testing.T) {
	builder := ChromeCapabilities()
	assert.True(t, builder.alwaysMatch.HasNot(protocol.CapabilityUnhandledPromptBehavior))
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch()[protocol.CapabilityUnhandledPromptBehavior])

	builder = ChromeCapabilities()
	builder.SetUnhandledPromptBehavior("string")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(protocol.CapabilityUnhandledPromptBehavior))
	assert.Equal(t, "string", opts.AlwaysMatch().GetString(protocol.CapabilityUnhandledPromptBehavior))
}

func TestChromeOptionsBuilder_Pref(t *testing.T) {
	builder := ChromeCapabilities()
	assert.Nil(t, builder.prefs)
	builder.Pref()
	opts := builder.Build()
	assert.Nil(t, opts.AlwaysMatch().GetOpts(ChromeCapabilityPreferencesName))

	builder = ChromeCapabilities()
	builder.Pref().Set("key", "value")
	opts = builder.Build()
	assert.True(t, opts.AlwaysMatch().Has(ChromeCapabilityPreferencesName))
	assert.Equal(t, "value", opts.AlwaysMatch().GetOpts(ChromeCapabilityPreferencesName).GetString("key"))
}

package webdriver

import (
	"encoding/base64"
	"os"
	"path"
	"testing"

	"github.com/mediabuyerbot/go-crx3"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
	"github.com/stretchr/testify/assert"
)

func TestChromeOptions(t *testing.T) {
	builder := ChromeOptions()

	extension := base64.StdEncoding.EncodeToString([]byte(`extension`))

	assert.Nil(t, builder.SetBrowserName("chrome"))
	assert.Nil(t, builder.SetBrowserVersion("91"))
	assert.Nil(t, builder.SetPlatformName("linux"))
	assert.Nil(t, builder.SetAcceptInsecureCerts(true))
	assert.Nil(t, builder.SetPageLoadStrategy("normal"))
	assert.Nil(t, builder.SetWindowRect(true))
	assert.Nil(t, builder.SetProxy(&w3c.Proxy{SocksPort: 8090}))
	assert.Nil(t, builder.SetUnhandledPromptBehavior("string"))
	assert.Nil(t, builder.SetTimeout(w3c.Timeout{Script: 9000}))
	assert.Nil(t, builder.SetDebuggerAddr("127.0.0.1:6666"))
	assert.Nil(t, builder.SetDetach(true))
	assert.Nil(t, builder.SetBinary("/path/to/chrome.bin"))
	assert.Nil(t, builder.SetMiniDumpPath("/path/to/dump"))
	assert.Nil(t, builder.SetLocalState("userProfile", "/profile"))
	assert.Nil(t, builder.SetPref("pref", "pref"))
	assert.Nil(t, builder.AddArgument("--headless", "--no-sniff"))
	assert.Nil(t, builder.AddExtension(extension))
	assert.Nil(t, builder.AddExcludeSwitches("--exclude", "--exclude2"))
	assert.Nil(t, builder.AddWindowTypes("window"))
	assert.Nil(t, builder.AddFirstMatch("browserName", "chrome"))

	dm := &DeviceMetrics{
		Width:      2000,
		Height:     4000,
		PixelRatio: 5,
		Touch:      true,
	}

	builder.MobileEmulation().
		SetDeviceName("name").
		SetUserAgent("userAgent").
		SetDeviceMetrics(dm).
		Set("customKey", "customValue").
		SetDeviceMetrics(nil)

	builder.PerfLoggingPreferences().
		TracingCategories("trace").
		BufferUsageReportingIntervalMillis(1).
		EnableNetwork(true).
		EnablePage(true).
		EnableTimeline(true).
		Set("key", "value")

	browserOptions := builder.Build()
	assert.NotNil(t, browserOptions)

	// always match
	alwaysMatch := browserOptions.AlwaysMatch()
	assert.Equal(t, "chrome", alwaysMatch.GetString(w3c.CapabilityBrowserName))
	assert.Equal(t, "91", alwaysMatch.GetString(w3c.CapabilityBrowserVersion))
	assert.Equal(t, "linux", alwaysMatch.GetString(w3c.CapabilityPlatformName))
	assert.True(t, alwaysMatch.GetBool(w3c.CapabilityAcceptInsecureCerts))
	assert.Equal(t, "normal", alwaysMatch.GetString(w3c.CapabilityPageLoadStrategy))
	assert.True(t, alwaysMatch.GetBool(w3c.CapabilitySetWindowRect))
	assert.Equal(t, 8090, alwaysMatch.Section("proxy").GetInt("socksProxyPort"))
	assert.Equal(t, "string", alwaysMatch.GetString(w3c.CapabilityUnhandledPromptBehavior))
	assert.Equal(t, uint(9000), alwaysMatch.Section(w3c.CapabilityTimeouts).GetUint("script"))
	assert.Equal(t, "127.0.0.1:6666", alwaysMatch.Section(ChromeOptionsKey).GetString(ChromeCapabilityDebuggerAddressName))
	assert.True(t, alwaysMatch.Section(ChromeOptionsKey).GetBool(ChromeCapabilityDetachName))
	assert.Equal(t, "/path/to/chrome.bin", alwaysMatch.Section(ChromeOptionsKey).GetString(ChromeCapabilityBinaryName))
	assert.Equal(t, "/path/to/dump", alwaysMatch.Section(ChromeOptionsKey).GetString(ChromeCapabilityMiniDumpPathName))
	assert.Equal(t, "/profile", alwaysMatch.Section(ChromeOptionsKey).Section("localState").GetString("userProfile"))
	assert.Equal(t, "pref", alwaysMatch.Section(ChromeOptionsKey).Section(ChromeCapabilityPreferencesName).GetString("pref"))
	assert.Len(t, alwaysMatch.Section(ChromeOptionsKey).GetStringSlice(ChromeCapabilityArgsName), 2)
	assert.Len(t, alwaysMatch.Section(ChromeOptionsKey).GetStringSlice(ChromeCapabilityExtensionName), 1)
	assert.Len(t, alwaysMatch.Section(ChromeOptionsKey).GetStringSlice(ChromeCapabilityWindowTypesName), 1)
	assert.Len(t, browserOptions.FirstMatch(), 1)

	// always match mobile emulation
	mobe := alwaysMatch.Section(ChromeOptionsKey).Section(ChromeCapabilityMobileEmulationName)
	assert.Equal(t, "name", mobe.GetString(mobileEmulationDeviceName))
	assert.Equal(t, "userAgent", mobe.GetString(mobileEmulationUserAgent))
	assert.Equal(t, "customValue", mobe.GetString("customKey"))
	assert.Equal(t, uint(2000), mobe.Section("deviceMetrics").GetUint("width"))
	assert.Equal(t, uint(4000), mobe.Section("deviceMetrics").GetUint("height"))
	assert.Equal(t, float64(5), mobe.Section("deviceMetrics").GetFloat("pixelRatio"))
	assert.True(t, mobe.Section("deviceMetrics").GetBool("touch"))

	// always match perfLogs
	perfLogs := alwaysMatch.Section(ChromeOptionsKey).Section(ChromeCapabilityPerfLoggingPrefsName)
	assert.Equal(t, "trace", perfLogs.GetString(perfLoggingTracingCategories))
	assert.Equal(t, uint(1), perfLogs.GetUint(perfLoggingBufferUsageReportingInterval))
	assert.True(t, perfLogs.GetBool(perfLoggingEnableNetwork))
	assert.True(t, perfLogs.GetBool(perfLoggingEnablePage))
	assert.True(t, perfLogs.GetBool(perfLoggingEnableTimeline))
	assert.Equal(t, "value", perfLogs.GetString("key"))

	// load bad extension
	err := builder.AddExtension("--")
	assert.Error(t, err)

	beforeExcludeSwitchesLen := len(builder.excludeSwitches)
	err = builder.AddExcludeSwitches("", "", "")
	assert.Nil(t, err)
	assert.Equal(t, beforeExcludeSwitchesLen, len(builder.excludeSwitches))

	beforeWindowTypesLen := len(builder.windowTypes)
	err = builder.AddWindowTypes("", "")
	assert.Nil(t, err)
	assert.Equal(t, beforeWindowTypesLen, len(builder.windowTypes))
}

func TestLoadChromeExtension(t *testing.T) {
	// load unpacked extension
	tmp := path.Join(os.TempDir(), "ext")
	err := os.Mkdir(tmp, os.ModePerm)
	assert.Nil(t, err)
	defer func() {
		err = os.RemoveAll(tmp)
		assert.Nil(t, err)
	}()
	b64, err := LoadChromeExtension(tmp)
	assert.Nil(t, err)
	assert.True(t, IsBase64(b64))

	// load zip extension
	err = crx3.Extension(tmp).Zip()
	assert.Nil(t, err)
	defer func() {
		err = os.Remove(tmp + ".zip")
	}()
	b64, err = LoadChromeExtension(tmp)
	assert.Nil(t, err)
	assert.True(t, IsBase64(b64))

	// load crx
	b64, err = LoadChromeExtension("./testdata/chrome/extension.crx")
	assert.Nil(t, err)
	assert.True(t, IsBase64(b64))

	// load bad crx
	b64, err = LoadChromeExtension("./testdata/chrome/none.crx")
	assert.Error(t, err)
	assert.False(t, IsBase64(b64))
}

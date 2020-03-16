package webdriver

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
)

const (
	// List of command-line arguments to use when starting Chrome. Arguments with an associated value
	// should be separated by a '=' sign (e.g., ['start-maximized', 'user-data-dir=/tmp/temp_profile']).
	// See here for a list of Chrome arguments.
	ChromeCapabilityArgsName = "args"

	// Path to the Chrome executable to use (on Mac OS X, this should be the actual binary,
	// not just the app. e.g., '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome')
	ChromeCapabilityBinaryName = "binary"

	// A list of Chrome extensions to install on startup. Each item in the list
	// should be a base-64 encoded packed Chrome extension (.crx)
	ChromeCapabilityExtensionName = "extensions"

	// A dictionary with each entry consisting of the name of the preference and its value.
	// These preferences are applied to the Local State file in the user data folder.
	ChromeCapabilityLocalStateName = "localState"

	// A dictionary with each entry consisting of the name of the preference and its value.
	// These preferences are only applied to the user profile in use.
	// See the 'Preferences' file in Chrome's user data directory for examples.
	ChromeCapabilityPreferencesName = "prefs"

	// If false, Chrome will be quit when ChromeDriver is killed, regardless of whether the session is quit.
	// If true, Chrome will only be quit if the session is quit (or closed). Note, if true, and the session is not quit,
	// ChromeDriver cannot clean up the temporary user data directory that the running Chrome instance is using.
	ChromeCapabilityDetachName = "detach"

	// An address of a Chrome debugger server to connect to, in the form of <hostname/ip:port>, e.g. '127.0.0.1:38947'
	ChromeCapabilityDebuggerAddressName = "debuggerAddress"

	// List of Chrome command line switches to exclude that ChromeDriver by default passes when starting Chrome.
	// Do not prefix switches with --.
	ChromeCapabilityExcludeSwitchesName = "excludeSwitches"

	// Directory to store Chrome minidumps . (Supported only on Linux.)
	ChromeCapabilityMiniDumpPathName = "minidumpPath"

	// A dictionary with either a value for “deviceName,” or values for “deviceMetrics” and “userAgent.”
	// Refer to Mobile Emulation for more information.
	ChromeCapabilityMobileEmulationName = "mobileEmulation"

	// An optional dictionary that specifies performance logging preferences. See below for more information.
	ChromeCapabilityPerfLoggingPrefsName = "perfLoggingPrefs"

	// A list of window types that will appear in the list of window handles.
	// For access to <webview> elements, include "webview" in this list.
	ChromeCapabilityWindowTypesName = "windowTypes"

	ChromeOptionsKey = "goog:chromeOptions"
)

const (
	zipExt = ".zip"
	crxExt = ".crx"
)

var ErrBase64Format = errors.New("webdriver: string does not match format base64")

type ChromeOptionsBuilder struct {
	//W3C Capabilities
	capabilities w3c.Capabilities

	// Chrome options
	chromeCapabilities w3c.Capabilities
	extensions         []string
	excludeSwitches    []string
	windowTypes        []string
	localState         w3c.Capabilities
	args               []string
	pref               w3c.Capabilities

	mobileEmulation *MobileEmulation
	perfLoggingPref *PerfLoggingPreferences

	firstMatch []w3c.Capabilities
}

func ChromeOptions() *ChromeOptionsBuilder {
	return &ChromeOptionsBuilder{
		capabilities: w3c.MakeCapabilities(),

		chromeCapabilities: w3c.MakeCapabilities(),
		extensions:         make([]string, 0),
		excludeSwitches:    make([]string, 0),
		windowTypes:        make([]string, 0),
		localState:         w3c.MakeCapabilities(),
		args:               make([]string, 0),
		pref:               w3c.MakeCapabilities(),

		firstMatch: make([]w3c.Capabilities, 0),
	}
}

func (b *ChromeOptionsBuilder) SetBrowserName(name string) *ChromeOptionsBuilder {
	_ = w3c.SetBrowserName(b.capabilities, name)
	return b
}

func (b *ChromeOptionsBuilder) SetBrowserVersion(version string) *ChromeOptionsBuilder {
	_ = w3c.SetBrowserVersion(b.capabilities, version)
	return b
}

func (b *ChromeOptionsBuilder) SetPlatformName(platform string) *ChromeOptionsBuilder {
	_ = w3c.SetPlatformName(b.capabilities, w3c.Platform(platform))
	return b
}

func (b *ChromeOptionsBuilder) SetAcceptInsecureCerts(flag bool) *ChromeOptionsBuilder {
	_ = w3c.SetAcceptInsecureCerts(b.capabilities, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetPageLoadStrategy(strategy string) *ChromeOptionsBuilder {
	_ = w3c.SetPageLoadStrategy(b.capabilities, strategy)
	return b
}

func (b *ChromeOptionsBuilder) SetWindowRect(flag bool) *ChromeOptionsBuilder {
	_ = w3c.SetWindowRect(b.capabilities, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetProxy(proxy *w3c.Proxy) *ChromeOptionsBuilder {
	_ = w3c.SetProxy(b.capabilities, proxy)
	return b
}

func (b *ChromeOptionsBuilder) SetUnhandledPromptBehavior(prompt string) *ChromeOptionsBuilder {
	_ = w3c.SetUnhandledPromptBehavior(b.capabilities, prompt)
	return b
}

func (b *ChromeOptionsBuilder) SetTimeout(timeout w3c.Timeout) *ChromeOptionsBuilder {
	_ = w3c.SetTimeout(b.capabilities, timeout)
	return b
}

func (b *ChromeOptionsBuilder) SetDebuggerAddr(addr string) *ChromeOptionsBuilder {
	_ = b.chromeCapabilities.Set(ChromeCapabilityDebuggerAddressName, addr)
	return b
}

func (b *ChromeOptionsBuilder) SetDetach(flag bool) *ChromeOptionsBuilder {
	_ = b.chromeCapabilities.Set(ChromeCapabilityDetachName, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetBinary(binPath string) *ChromeOptionsBuilder {
	b.chromeCapabilities.Set(ChromeCapabilityBinaryName, binPath)
	return b
}

func (b *ChromeOptionsBuilder) SetMiniDumpPath(path string) *ChromeOptionsBuilder {
	b.chromeCapabilities.Set(ChromeCapabilityMiniDumpPathName, path)
	return b
}

func (b *ChromeOptionsBuilder) SetLocalState(key string, value interface{}) *ChromeOptionsBuilder {
	b.localState.Set(key, value)
	return b
}

func (b *ChromeOptionsBuilder) SetPref(key string, value interface{}) *ChromeOptionsBuilder {
	b.pref.Set(key, value)
	return b
}

func (b *ChromeOptionsBuilder) AddArgument(arg ...string) *ChromeOptionsBuilder {
	b.args = append(b.args, arg...)
	return b
}

func (b *ChromeOptionsBuilder) AddExtension(base64 string) error {
	if ok := IsBase64(base64); !ok {
		return ErrBase64Format
	}
	b.extensions = append(b.extensions, base64)
	return nil
}

func (b *ChromeOptionsBuilder) AddExcludeSwitches(exclude ...string) *ChromeOptionsBuilder {
	for _, arg := range exclude {
		if len(arg) == 0 {
			continue
		}
		b.excludeSwitches = append(b.excludeSwitches, arg)
	}
	return b
}

func (b *ChromeOptionsBuilder) AddWindowTypes(types ...string) *ChromeOptionsBuilder {
	for _, arg := range types {
		if len(arg) == 0 {
			continue
		}
		b.windowTypes = append(b.windowTypes, arg)
	}
	return b
}

func (b *ChromeOptionsBuilder) AddFirstMatch(key string, value interface{}) *ChromeOptionsBuilder {
	if len(key) > 0 {
		cap := w3c.MakeCapabilities()
		cap.Set(key, value)
		b.firstMatch = append(b.firstMatch, cap)
	}
	return b
}

func (b *ChromeOptionsBuilder) MobileEmulation() *MobileEmulation {
	if b.mobileEmulation == nil {
		b.mobileEmulation = &MobileEmulation{opts: w3c.MakeCapabilities()}
	}
	return b.mobileEmulation
}

func (b *ChromeOptionsBuilder) PerfLoggingPreferences() *PerfLoggingPreferences {
	if b.perfLoggingPref == nil {
		b.perfLoggingPref = &PerfLoggingPreferences{opts: w3c.MakeCapabilities()}
	}
	return b.perfLoggingPref
}

func (b *ChromeOptionsBuilder) Build() w3c.BrowserOptions {
	if len(b.extensions) > 0 {
		b.chromeCapabilities[ChromeCapabilityExtensionName] = b.extensions
	}
	if len(b.localState) > 0 {
		b.chromeCapabilities[ChromeCapabilityLocalStateName] = b.localState
	}
	if len(b.excludeSwitches) > 0 {
		b.chromeCapabilities[ChromeCapabilityExcludeSwitchesName] = b.excludeSwitches
	}
	if len(b.windowTypes) > 0 {
		b.chromeCapabilities[ChromeCapabilityWindowTypesName] = b.windowTypes
	}
	if len(b.args) > 0 {
		b.chromeCapabilities[ChromeCapabilityArgsName] = b.args
	}
	if len(b.pref) > 0 {
		b.chromeCapabilities[ChromeCapabilityPreferencesName] = b.pref
	}
	if b.mobileEmulation != nil && len(b.mobileEmulation.opts) > 0 {
		b.chromeCapabilities[ChromeCapabilityMobileEmulationName] = b.mobileEmulation.opts
	}
	if b.perfLoggingPref != nil && len(b.perfLoggingPref.opts) > 0 {
		b.chromeCapabilities[ChromeCapabilityPerfLoggingPrefsName] = b.perfLoggingPref.opts
	}

	b.capabilities.Set(ChromeOptionsKey, b.chromeCapabilities)

	return w3c.NewBrowserOptions(b.capabilities, b.firstMatch)
}

func LoadChromeExtension(extensionPath string) (base64 string, err error) {
	extension := crx3.Extension(extensionPath)
	if extension.IsZip() || extension.IsDir() {
		err := extension.Pack(nil)
		if err != nil {
			return base64, err
		}

		crx := strings.TrimRight(extension.String(), "/")
		crx = strings.TrimRight(crx, zipExt)
		crx = crx + crxExt
		extension = crx3.Extension(crx)
	}

	if !extension.IsCRX3() {
		return base64, crx3.ErrUnsupportedFileFormat
	}

	b, err := extension.ToBase64()
	if err != nil {
		return base64, err
	}
	return string(b), nil
}

func IsBase64(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

const (
	mobileEmulationDeviceName    = "deviceName"
	mobileEmulationDeviceMetrics = "deviceMetrics"
	mobileEmulationUserAgent     = "userAgent"
)

type MobileEmulation struct {
	opts w3c.Capabilities
}

func (e *MobileEmulation) Set(key string, value interface{}) *MobileEmulation {
	e.opts.Set(key, value)
	return e
}

func (e *MobileEmulation) SetDeviceName(name string) *MobileEmulation {
	return e.Set(mobileEmulationDeviceName, name)
}

func (e *MobileEmulation) SetDeviceMetrics(m *DeviceMetrics) *MobileEmulation {
	if m == nil {
		return e
	}
	return e.Set(mobileEmulationDeviceMetrics, m.Capabilities())
}

func (e *MobileEmulation) SetUserAgent(agent string) *MobileEmulation {
	return e.Set(mobileEmulationUserAgent, agent)
}

const (
	perfLoggingEnableNetwork                = "enableNetwork"
	perfLoggingEnableTimeline               = "enableTimeline"
	perfLoggingEnablePage                   = "enablePage"
	perfLoggingTracingCategories            = "tracingCategories"
	perfLoggingBufferUsageReportingInterval = "bufferUsageReportingInterval"
)

type PerfLoggingPreferences struct {
	opts w3c.Capabilities
}

func (pp *PerfLoggingPreferences) Set(key string, value interface{}) *PerfLoggingPreferences {
	pp.opts.Set(key, value)
	return pp
}

func (pp *PerfLoggingPreferences) EnableNetwork(flag bool) *PerfLoggingPreferences {
	pp.opts.Set(perfLoggingEnableNetwork, flag)
	return pp
}

func (pp *PerfLoggingPreferences) EnableTimeline(flag bool) *PerfLoggingPreferences {
	pp.opts.Set(perfLoggingEnableTimeline, flag)
	return pp
}

func (pp *PerfLoggingPreferences) EnablePage(flag bool) *PerfLoggingPreferences {
	pp.opts.Set(perfLoggingEnablePage, flag)
	return pp
}

func (pp *PerfLoggingPreferences) TracingCategories(s string) *PerfLoggingPreferences {
	pp.opts.Set(perfLoggingTracingCategories, s)
	return pp
}

func (pp *PerfLoggingPreferences) BufferUsageReportingIntervalMillis(v uint) *PerfLoggingPreferences {
	pp.opts.Set(perfLoggingBufferUsageReportingInterval, v)
	return pp
}

type DeviceMetrics struct {
	Width      uint    `json:"width"`
	Height     uint    `json:"height"`
	PixelRatio float64 `json:"pixelRatio"`
	Touch      bool    `json:"touch,omitempty"`
}

func (dm *DeviceMetrics) Capabilities() w3c.Capabilities {
	return w3c.Capabilities{
		"width":      dm.Width,
		"height":     dm.Height,
		"pixelRatio": dm.PixelRatio,
		"touch":      dm.Touch,
	}
}

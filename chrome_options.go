package webdriver

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mediabuyerbot/go-webdriver/pkg/protocol"

	"github.com/mediabuyerbot/go-crx3"
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
)

type Builder interface {
	Build() protocol.Options
}

const (
	zipExt = ".zip"
	crxExt = ".crx"
)

type ChromeOptions struct {
	proxy       *protocol.Proxy
	firstMatch  []protocol.O
	alwaysMatch protocol.O
}

func (o *ChromeOptions) Proxy() *protocol.Proxy {
	return o.proxy
}

func (o *ChromeOptions) FirstMatch() []protocol.O {
	return o.firstMatch
}

func (o *ChromeOptions) AlwaysMatch() protocol.O {
	return o.alwaysMatch
}

type ChromeOptionsBuilder struct {
	extensions       *Extensions
	args             *Arguments
	excludeSwitches  []string
	mobileEmulation  *MobileEmulation
	localState       *LocalState
	perfLoggingPrefs *PerfLoggingPreferences
	windowTypes      []string

	prefs protocol.O
	proxy *protocol.Proxy

	firstMatch  []protocol.O
	alwaysMatch protocol.O
}

func ChromeCapabilities() *ChromeOptionsBuilder {
	return &ChromeOptionsBuilder{
		excludeSwitches: make([]string, 0),
		windowTypes:     make([]string, 0),
		prefs:           protocol.MakeOptions(),

		firstMatch:  make([]protocol.O, 0),
		alwaysMatch: make(protocol.O),
	}
}

func (b *ChromeOptionsBuilder) SetBrowserName(name string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityBrowserName, name)
	return b
}

func (b *ChromeOptionsBuilder) SetBrowserVersion(version string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityBrowserVersion, version)
	return b
}

func (b *ChromeOptionsBuilder) SetPlatformName(platform string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityPlatformName, platform)
	return b
}

func (b *ChromeOptionsBuilder) SetPlatformVersion(version string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityPlatformVersion, version)
	return b
}

func (b *ChromeOptionsBuilder) SetAcceptInsecureCerts(flag bool) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityAcceptInsecureCerts, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetPageLoadStrategy(strategy string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityPageLoadStrategy, strategy)
	return b
}

func (b *ChromeOptionsBuilder) SetWindowRect(flag bool) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilitySetWindowRect, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetProxy(p *protocol.Proxy) *ChromeOptionsBuilder {
	b.proxy = p
	return b
}

func (b *ChromeOptionsBuilder) SetDebuggerAddr(addr string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(ChromeCapabilityDebuggerAddressName, addr)
	return b
}

func (b *ChromeOptionsBuilder) SetUnhandledPromptBehavior(s string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(protocol.CapabilityUnhandledPromptBehavior, s)
	return b
}

func (b *ChromeOptionsBuilder) SetDetach(flag bool) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(ChromeCapabilityDetachName, flag)
	return b
}

func (b *ChromeOptionsBuilder) SetBinary(binPath string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(ChromeCapabilityBinaryName, binPath)
	return b
}

func (b *ChromeOptionsBuilder) SetMiniDumpPath(path string) *ChromeOptionsBuilder {
	b.alwaysMatch.Set(ChromeCapabilityMiniDumpPathName, path)
	return b
}

func (b *ChromeOptionsBuilder) AddExcludeSwitches(exclude ...string) *ChromeOptionsBuilder {
	b.excludeSwitches = append(b.excludeSwitches, exclude...)
	return b
}

func (b *ChromeOptionsBuilder) AddWindowTypes(types ...string) *ChromeOptionsBuilder {
	b.windowTypes = append(b.windowTypes, types...)
	return b
}

func (b *ChromeOptionsBuilder) AddFirstMatch(key string, value interface{}) *ChromeOptionsBuilder {
	b.firstMatch = append(b.firstMatch, protocol.O{key: value})
	return b
}

func (b *ChromeOptionsBuilder) LocalState() *LocalState {
	if b.localState == nil {
		b.localState = &LocalState{
			opts: protocol.MakeOptions(),
		}
	}
	return b.localState
}

func (b *ChromeOptionsBuilder) Arguments() *Arguments {
	if b.args == nil {
		b.args = &Arguments{
			opts: make([]string, 0, 32),
		}
	}
	return b.args
}

func (b *ChromeOptionsBuilder) Extensions() *Extensions {
	if b.extensions == nil {
		b.extensions = &Extensions{
			opts: make([]string, 0, 32),
		}
	}
	return b.extensions
}

func (b *ChromeOptionsBuilder) MobileEmulation() *MobileEmulation {
	if b.mobileEmulation == nil {
		b.mobileEmulation = &MobileEmulation{
			opts: protocol.MakeOptions(),
		}
	}
	return b.mobileEmulation
}

func (b *ChromeOptionsBuilder) PerfLoggingPreferences() *PerfLoggingPreferences {
	if b.perfLoggingPrefs == nil {
		b.perfLoggingPrefs = &PerfLoggingPreferences{
			opts: protocol.MakeOptions(),
		}
	}
	return b.perfLoggingPrefs
}

func (b *ChromeOptionsBuilder) Build() protocol.Options {

	// extensions
	if b.extensions != nil && len(b.extensions.opts) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityExtensionName, b.extensions.opts)
	}

	// local state
	if b.localState != nil && len(b.localState.opts) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityLocalStateName, b.localState.opts)
	}

	// proxy
	if b.proxy != nil {
		b.alwaysMatch.Set(protocol.CapabilityProxy, b.proxy)
	}

	// args
	if b.args != nil {
		b.alwaysMatch.Set(ChromeCapabilityArgsName, b.args.opts)
	}

	// exclude switches
	if len(b.excludeSwitches) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityExcludeSwitchesName, b.excludeSwitches)
	}

	// window types
	if len(b.windowTypes) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityWindowTypesName, b.windowTypes)
	}

	// mobile emulation
	if b.mobileEmulation != nil && len(b.mobileEmulation.opts) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityMobileEmulationName, b.mobileEmulation.opts)
	}

	// prefs
	if b.prefs != nil && len(b.prefs) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityPreferencesName, b.prefs)
	}

	// perfLoggingPrefs
	if b.perfLoggingPrefs != nil && len(b.perfLoggingPrefs.opts) > 0 {
		b.alwaysMatch.Set(ChromeCapabilityPerfLoggingPrefsName, b.perfLoggingPrefs.opts)
	}

	return &ChromeOptions{
		alwaysMatch: b.alwaysMatch,
		firstMatch:  b.firstMatch,
		proxy:       b.proxy,
	}
}

type Extensions struct {
	opts []string
}

func (e *Extensions) Add(filepath ...string) error {
	for _, fp := range filepath {
		base64, err := e.add(fp)
		if err != nil {
			return fmt.Errorf("extension %s error %v", fp, err)
		}
		e.opts = append(e.opts, base64)
	}
	return nil
}

func (e *Extensions) add(path string) (s string, err error) {
	extension := crx3.Extension(path)
	if extension.IsZip() || extension.IsDir() {
		err := extension.Pack(nil)
		if err != nil {
			return s, err
		}
		crx := strings.TrimRight(extension.String(), zipExt)
		crx = filepath.Join(crx, crxExt)
		if err := crx3.Extension(crx).Pack(nil); err != nil {
			return s, err
		}
		extension = crx3.Extension(crx)
	}

	if !extension.IsCRX3() {
		return s, crx3.ErrUnsupportedFileFormat
	}

	base64, err := extension.ToBase64()
	if err != nil {
		return s, err
	}
	return string(base64), nil
}

type Arguments struct {
	opts []string
}

// WithHeadless run in headless mode, i.e., without a UI or display server dependencies.
func (a *Arguments) WithHeadless() *Arguments {
	return a.Add("--headless")
}

// WithUserDataPath set a custom user profile to use.
func (a *Arguments) WithUserDataPath(path string) *Arguments {
	return a.Add("user-data-dir=" + path)
}

// WithStartMaximize starts the browser maximized, regardless of any previous settings.
func (a *Arguments) WithStartMaximized() *Arguments {
	return a.Add("--start-maximized")
}

func (a *Arguments) Add(v ...string) *Arguments {
	a.opts = append(a.opts, v...)
	return a
}

type LocalState struct {
	opts protocol.O
}

func (ls *LocalState) Set(key string, value interface{}) *LocalState {
	ls.opts.Set(key, value)
	return ls
}

const (
	mobileEmulationDeviceName    = "deviceName"
	mobileEmulationDeviceMetrics = "deviceMetrics"
	mobileEmulationUserAgent     = "userAgent"
)

type MobileEmulation struct {
	opts protocol.O
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
	return e.Set(mobileEmulationDeviceMetrics, m.ToOptions())
}

func (e *MobileEmulation) SetUserAgent(agent string) *MobileEmulation {
	return e.Set(mobileEmulationUserAgent, agent)
}

type PerfLoggingPreferences struct {
	opts protocol.O
}

func (pp *PerfLoggingPreferences) Set(key string, value interface{}) *PerfLoggingPreferences {
	pp.opts.Set(key, value)
	return pp
}

func (pp *PerfLoggingPreferences) EnableNetwork(flag bool) *PerfLoggingPreferences {
	pp.opts.Set("enableNetwork", flag)
	return pp
}

func (pp *PerfLoggingPreferences) EnableTimeline(flag bool) *PerfLoggingPreferences {
	pp.opts.Set("enableTimeline", flag)
	return pp
}

func (pp *PerfLoggingPreferences) EnablePage(flag bool) *PerfLoggingPreferences {
	pp.opts.Set("enablePage", flag)
	return pp
}

func (pp *PerfLoggingPreferences) TracingCategories(s string) *PerfLoggingPreferences {
	pp.opts.Set("tracingCategories", s)
	return pp
}

func (pp *PerfLoggingPreferences) BufferUsageReportingIntervalMillis(v uint) *PerfLoggingPreferences {
	pp.opts.Set("bufferUsageReportingInterval", v)
	return pp
}

type DeviceMetrics struct {
	Width      uint    `json:"width"`
	Height     uint    `json:"height"`
	PixelRatio float64 `json:"pixelRatio"`
	Touch      *bool   `json:"touch,omitempty"`
}

func (dm *DeviceMetrics) ToOptions() protocol.O {
	touch := true
	if dm.Touch != nil {
		touch = *dm.Touch
	}
	return protocol.O{
		"width":      dm.Width,
		"height":     dm.Height,
		"pixelRatio": dm.PixelRatio,
		"touch":      touch,
	}
}

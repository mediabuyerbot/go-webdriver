package geckodriver

import (
	"fmt"
	"math"
	"os/exec"
)

const (
	// Connect to an existing Firefox instance
	flagConnExisting = "--connect-existing"

	// Attach browser toolbox debugger for Firefox
	flagJsDebugger = "--jsdebugger"

	// Log level verbosity (-v for debug and -vv for trace level)
	flagVerbose = "-v"

	// Prints version and copying information
	flagVersion = "-V"
)

const (
	// Path to the Firefox binary
	argsBinary = "--binary=%s"

	// Set Gecko log level [possible values: fatal, error, warn, info, config, debug, trace]
	argsLogLevel = "--log=%s"

	// Host to use to connect to Gecko [default: 127.0.0.1]
	argsMarionetteHost = "--marionette-host=%s"

	// Port to use to connect to Gecko [default: system -allocated port]
	argsMarionettePort = "--marionette-port=%d"

	// Host IP to use for WebDriver server [default: 127.0.0.1]
	argsHost = "--host=%s"

	// Port to use for WebDriver server [default: 4444]
	argsPort = "--port=%d"
)

const (
	Fatal   LogLevel = "fatal"
	Error   LogLevel = "error"
	Warning LogLevel = "warn"
	Info    LogLevel = "info"
	Config  LogLevel = "config"
	Debug   LogLevel = "debug"
	Trace   LogLevel = "trace"
)

type LogLevel string

func (ll LogLevel) Validate() error {
	switch ll {
	case Fatal, Debug, Info, Warning, Error, Config, Trace:
		return nil
	default:
		return fmt.Errorf("unknown log level %s", ll)
	}
}

type flags struct {
	store map[string]string
}

func newFlags() *flags {
	return &flags{
		store: make(map[string]string),
	}
}

func (f *flags) ConnExisting() {
	f.store[flagConnExisting] = flagConnExisting
}

func (f *flags) JsDebugger() {
	f.store[flagJsDebugger] = flagJsDebugger
}

func (f *flags) Verbose() {
	f.store[flagVerbose] = flagVerbose
}

func (f *flags) Has(flag string) bool {
	_, ok := f.store[flag]
	return ok
}

func (f *flags) Remove(flag string) {
	delete(f.store, flag)
}

func (f *flags) Version() {
	f.store[flagVersion] = flagVersion
}

func (f *flags) Build() []string {
	flags := make([]string, len(f.store))
	i := 0
	for v := range f.store {
		flags[i] = v
		i++
	}
	return flags
}

type arguments struct {
	store map[string]string
}

func newArguments() *arguments {
	return &arguments{
		store: make(map[string]string),
	}
}

func (a *arguments) Has(flag string) bool {
	_, ok := a.store[flag]
	return ok
}

func (a *arguments) Remove(flag string) {
	delete(a.store, flag)
}

func (a *arguments) SetBinary(bin string) error {
	binPath, err := exec.LookPath(bin)
	if err != nil {
		return err
	}
	a.store[argsBinary] = fmt.Sprintf(argsBinary, binPath)
	return nil
}

func (a *arguments) SetLogLevel(ll LogLevel) error {
	if err := ll.Validate(); err != nil {
		return err
	}
	a.store[argsLogLevel] = fmt.Sprintf(argsLogLevel, string(ll))
	return nil
}

func (a *arguments) SetMarionetteHost(host string) {
	if len(host) == 0 {
		return
	}
	a.store[argsMarionetteHost] = fmt.Sprintf(argsMarionetteHost, host)
}

func (a *arguments) SetMarionettePort(port int) error {
	if err := isValidPortRange(port); err != nil {
		return err
	}
	a.store[argsMarionettePort] = fmt.Sprintf(argsMarionettePort, port)
	return nil
}

func (a *arguments) SetHost(host string) {
	a.store[argsHost] = fmt.Sprintf(argsHost, host)
}

func (a *arguments) SetPort(port int) error {
	if err := isValidPortRange(port); err != nil {
		return err
	}
	a.store[argsPort] = fmt.Sprintf(argsPort, port)
	return nil
}

func (a *arguments) Build() []string {
	flags := make([]string, len(a.store))
	i := 0
	for _, v := range a.store {
		flags[i] = v
		i++
	}
	return flags
}

func isValidPortRange(port int) error {
	if port <= 0 || port >= math.MaxUint16 {
		return fmt.Errorf("out of range port %d", port)
	}
	return nil
}

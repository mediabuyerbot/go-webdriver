package chromedriver

import (
	"fmt"
	"math"
	"os"
)

const (
	// port to listen on
	argsPort = "--port=%d"

	// adb server port
	argsADBPort = "--adb-port=%d"

	// set log level: ALL, DEBUG, INFO, WARNING, SEVERE, OFF
	argsLogLevel = "--log-level=%s"

	// write server log to file instead of stderr, increases log level to INFO
	argsLogPath = "--log-path=%s"

	// log verbosely (equivalent to --log-level=ALL)
	argsVerbose = "--verbose"

	// log nothing (equivalent to --log-level=OFF)
	argsSilent = "--silent"

	// append log file instead of rewriting
	argsAppendLog = "--append-log"

	// (experimental) log verbosely and don't truncate long strings so that the log can be replayed.
	argsReplayable = "--replayable"

	// print the version number and exit
	argsVersion = "--version"

	// base URL path prefix for commands, e.g. wd/url
	argsURLBase = "--url-base=%s"

	// comma-separated whitelist of remote IP addresses which are allowed to connect to ChromeDriver
	argsWhitelistedIps = "--whitelisted-ips=%s"

	// add readable timestamps to log
	argsReadableTimestamp = "--readable-timestamp"
)

const (
	// ALL, DEBUG, INFO, WARNING, SEVERE, OFF
	All     LogLevel = "ALL"
	Debug   LogLevel = "DEBUG"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Severe  LogLevel = "SEVERE"
	Off     LogLevel = "OFF"
)

type LogLevel string

func (ll LogLevel) Validate() error {
	switch ll {
	case All, Debug, Info, Warning, Severe, Off:
		return nil
	default:
		return fmt.Errorf("unknown log level %s", ll)
	}
}

type arguments struct {
	args map[string]string
}

func newArguments() *arguments {
	return &arguments{
		args: make(map[string]string),
	}
}

func (a *arguments) setPort(port int) error {
	if err := isValidPortRange(port); err != nil {
		return err
	}
	a.args[argsPort] = fmt.Sprintf(argsPort, port)
	return nil
}

func (a *arguments) setADBPort(port int) error {
	if err := isValidPortRange(port); err != nil {
		return err
	}
	a.args[argsADBPort] = fmt.Sprintf(argsADBPort, port)
	return nil
}

func (a *arguments) setVerbose() {
	a.args[argsVersion] = argsVerbose
}

func (a *arguments) setSilent() {
	a.args[argsSilent] = argsSilent
}

func (a *arguments) setAppendLog() {
	a.args[argsAppendLog] = argsAppendLog
}

func (a *arguments) setReplayable() {
	a.args[argsReplayable] = argsReplayable
}

func (a *arguments) showVersion() {
	a.args[argsVersion] = argsVersion
}

func (a *arguments) setReadableTimestamp() {
	a.args[argsReadableTimestamp] = argsReadableTimestamp
}

func (a *arguments) setBaseURL(url string) {
	a.args[argsURLBase] = fmt.Sprintf(argsURLBase, url)
}

func (a *arguments) setLogLevel(l LogLevel) error {
	if err := l.Validate(); err != nil {
		return err
	}
	a.args[argsLogLevel] = fmt.Sprintf(argsLogLevel, l)
	return nil
}

func (a *arguments) setWhitelistedIps(ips string) {
	a.args[argsWhitelistedIps] = fmt.Sprintf(argsWhitelistedIps, ips)
}

func (a *arguments) setLogPath(logPath string) error {
	if err := createFileIfNotExist(logPath); err != nil {
		return err
	}
	a.args[argsLogPath] = fmt.Sprintf(argsLogPath, logPath)
	return nil
}

func (a *arguments) build() []string {
	var (
		args = make([]string, len(a.args))
		idx  int
	)
	for _, arg := range a.args {
		args[idx] = arg
		idx++
	}
	return args
}

func createFileIfNotExist(filepath string) error {
	_, err := os.Stat(filepath)
	if err == nil {
		return nil
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func isValidPortRange(port int) error {
	if port <= 0 || port >= math.MaxUint16 {
		return fmt.Errorf("out of range port %d", port)
	}
	return nil
}

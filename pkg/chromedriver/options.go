package chromedriver

import (
	"io"
	"strings"
)

type Option func(cd *Process)

// WithPort port to listen on
func WithPort(port int) Option {
	return func(p *Process) {
		if err := p.args.setPort(port); err != nil {
			panic(err)
		}
	}
}

// WithAdbPort adb server port
func WithAdbPort(port int) Option {
	return func(p *Process) {
		if err := p.args.setADBPort(port); err != nil {
			panic(err)
		}
	}
}

// WithVerbose log verbosely
func WithVerbose() Option {
	return func(p *Process) {
		p.args.setVerbose()
	}
}

// WithSilent log nothing
func WithSilent() Option {
	return func(p *Process) {
		p.args.setSilent()
	}
}

// WithAppendLog append log file instead of rewriting
func WithAppendLog() Option {
	return func(p *Process) {
		p.args.setAppendLog()
	}
}

// WithReplayable (experimental) log verbosely and don't truncate long strings so that the log can be replayed.
func WithReplayable() Option {
	return func(p *Process) {
		p.args.setReplayable()
	}
}

// WithShowVersion print the version number and exit
func WithShowVersion() Option {
	return func(p *Process) {
		p.args.showVersion()
	}
}

// WithBaseURL base URL path prefix for commands, e.g. wd/url
func WithBaseURL(url string) Option {
	return func(p *Process) {
		p.args.setBaseURL(url)
	}
}

// WithLogLevel log level: ALL, DEBUG, INFO, WARNING, SEVERE, OFF
func WithLogLevel(ll LogLevel) Option {
	return func(p *Process) {
		if err := p.args.setLogLevel(ll); err != nil {
			panic(err)
		}
	}
}

// WithWhitelistedIps whitelist of remote IP addresses which are allowed to connect to ChromeDriver
func WithWhitelistedIps(ips []string) Option {
	return func(p *Process) {
		p.args.setWhitelistedIps(strings.Join(ips, ","))
	}
}

// WithLogPath write server log to file instead of stderr
func WithLogPath(path string) Option {
	return func(p *Process) {
		if err := p.args.setLogPath(path); err != nil {
			panic(err)
		}
	}
}

// WithReadableTimestamp add readable timestamps to log
func WithReadableTimestamp() Option {
	return func(p *Process) {
		p.args.setReadableTimestamp()
	}
}

func WithRunHook(hook Hook) Option {
	return func(p *Process) {
		p.runHook = hook
	}
}

func WithStopHook(hook Hook) Option {
	return func(p *Process) {
		p.stopHook = hook
	}
}

func WithStderr(writer io.Writer) Option {
	return func(p *Process) {
		p.stderr = writer
	}
}

func WithStdout(writer io.Writer) Option {
	return func(p *Process) {
		p.stdout = writer
	}
}

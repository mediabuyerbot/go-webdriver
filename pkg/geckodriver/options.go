package geckodriver

import "io"

type Option func(cd *Process)

func WithConnExisting() Option {
	return func(p *Process) {
		p.flags.ConnExisting()
	}
}

func WithJsDebugger() Option {
	return func(p *Process) {
		p.flags.JsDebugger()
	}
}

func WithVerbose() Option {
	return func(p *Process) {
		p.flags.Verbose()
	}
}

func WithShowVersion() Option {
	return func(p *Process) {
		p.flags.Version()
	}
}

func WithBinary(bin string) Option {
	return func(p *Process) {
		if err := p.args.SetBinary(bin); err != nil {
			panic(err)
		}
	}
}

func WithLogLevel(ll LogLevel) Option {
	return func(p *Process) {
		if err := p.args.SetLogLevel(ll); err != nil {
			panic(err)
		}
	}
}

func WithMarionetteHost(host string) Option {
	return func(p *Process) {
		p.args.SetMarionetteHost(host)
	}
}

func WithMarionettePort(port int) Option {
	return func(p *Process) {
		if err := p.args.SetMarionettePort(port); err != nil {
			panic(err)
		}
	}
}

func WithHost(host string) Option {
	return func(p *Process) {
		p.args.SetHost(host)
	}
}

func WithPort(port int) Option {
	return func(p *Process) {
		if err := p.args.SetPort(port); err != nil {
			panic(err)
		}
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

package chromedriver

import (
	"time"
)

type Option func(cd *ChromeDriver)

func WithStartTime(duration time.Duration) Option {
	return func(cd *ChromeDriver) {
		cd.startTimeout = duration
	}
}

func WithThreads(threads int) Option {
	return func(cd *ChromeDriver) {
		cd.threads = threads
	}
}

func WithPort(port int) Option {
	return func(cd *ChromeDriver) {
		cd.port = port
	}
}

func WithAdbPort(port int) Option {
	return func(cd *ChromeDriver) {
		cd.adbPort = port
	}
}

func WithHost(host string) Option {
	return func(cd *ChromeDriver) {
		cd.host = host
	}
}

func WithCommand(cmd string) Option {
	return func(cd *ChromeDriver) {
		cd.path = cmd
	}
}

func WithBaseURL(baseURL string) Option {
	return func(cd *ChromeDriver) {
		cd.baseURL = baseURL
	}
}

func WithLogLevel(logLevel LogLevel) Option {
	return func(cd *ChromeDriver) {
		cd.logLevel = logLevel
	}
}

func WithOnStartHook(h Hook) Option {
	return func(cd *ChromeDriver) {
		cd.onStartHook = h
	}
}

func WithOnCloseHook(h Hook) Option {
	return func(cd *ChromeDriver) {
		cd.onCloseHook = h
	}
}

package chromedriver

import "time"

type Option func(cd *ChromeDriver)

func WithChromeDriverStartTime(duration time.Duration) Option {
	return func(cd *ChromeDriver) {
		cd.startTimeout = duration
	}
}

func WithChromeDriverThreads(threads int) Option {
	return func(cd *ChromeDriver) {
		cd.threads = threads
	}
}

func WithChromeDriverPort(port int) Option {
	return func(cd *ChromeDriver) {
		cd.port = port
	}
}

func WithChromeDriverAdbPort(port int) Option {
	return func(cd *ChromeDriver) {
		cd.adbPort = port
	}
}

func WithChromeDriverHost(host string) Option {
	return func(cd *ChromeDriver) {
		cd.host = host
	}
}

func WithChromeDriverCommand(cmd string) Option {
	return func(cd *ChromeDriver) {
		cd.path = cmd
	}
}

func WithChromeDriverBaseURL(baseURL string) Option {
	return func(cd *ChromeDriver) {
		cd.baseURL = baseURL
	}
}

func WithChromeDriverLogLevel(logLevel LogLevel) Option {
	return func(cd *ChromeDriver) {
		cd.logLevel = logLevel
	}
}

package process

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
)

type LogLevelChromeDriver string

const (
	DefaultChromeDriverCommand   = "chromedriver"
	DefaultChromeDriverStartTime = time.Duration(20 * time.Second)
	DefaultChromeDriverThreads   = 6
	DefaultChromeDriverPort      = 9165
	DefaultChromeDriverHost      = "127.0.0.1"
	DefaultChromeDriverLogPath   = "chromedriver.log"

	// ALL, DEBUG, INFO, WARNING, SEVERE, OFF
	AllLogLevelChromeDriver     LogLevelChromeDriver = "ALL"
	DebugLogLevelChromeDriver   LogLevelChromeDriver = "DEBUG"
	InfoLogLevelChromeDriver    LogLevelChromeDriver = "INFO"
	WarningLogLevelChromeDriver LogLevelChromeDriver = "WARNING"
	SevereLogLevelChromeDriver  LogLevelChromeDriver = "SEVERE"
	OffLogLevelChromeDriver     LogLevelChromeDriver = "OFF"
)

var (
	ErrTimeOut                    = errors.New("process: timout error")
	ErrChromeDriverAlreadyRunning = errors.New("process: chromedriver already running")
	ErrChromeDriverNotRunning     = errors.New("process: chromedriver not running")
)

type ChromeDriverOption func(cd *chromeDriver)

type chromeDriver struct {
	// base URL path prefix for commands, e.g. wd/url
	baseURL string
	path    string
	threads int

	logPath  string
	logLevel LogLevelChromeDriver

	startTimeout time.Duration
	port         int
	adbPort      int
	host         string
	proc         *exec.Cmd
	out          bytes.Buffer
}

func NewChromeDriver(opts ...ChromeDriverOption) Process {
	cd := &chromeDriver{
		path:         DefaultChromeDriverCommand,
		startTimeout: DefaultChromeDriverStartTime,
		threads:      DefaultChromeDriverThreads,
		port:         DefaultChromeDriverPort,
		host:         DefaultChromeDriverHost,
		logLevel:     AllLogLevelChromeDriver,
	}

	for _, opt := range opts {
		opt(cd)
	}

	return cd
}

func (d *chromeDriver) Out() bytes.Buffer {
	return d.out
}

func (d *chromeDriver) Run(ctx context.Context) error {
	if d.proc != nil {
		return ErrChromeDriverAlreadyRunning
	}

	args := []string{
		"--port=" + strconv.Itoa(d.port),
		"--http-threads=" + strconv.Itoa(d.threads),
	}
	if len(d.baseURL) > 0 {
		args = append(args, "--url-base="+d.baseURL)
	}
	if d.adbPort > 0 {
		args = append(args, "--adb-port="+strconv.Itoa(d.adbPort))
	}
	if len(d.logPath) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		d.logPath = home + DefaultChromeDriverLogPath
	}

	d.proc = exec.Command(d.path, args...)
	d.proc.Stdout = &d.out
	if err := d.proc.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- d.proc.Wait()
	}()
	select {
	case <-ctx.Done():
		err := d.proc.Process.Kill()
		if err == nil {
			return ErrTimeOut
		}
		return err
	case err := <-done:
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *chromeDriver) Stop(_ context.Context) error {
	if d.proc == nil {
		return ErrChromeDriverNotRunning
	}
	defer func() {
		d.proc = nil
	}()
	return d.proc.Process.Signal(os.Interrupt)
}

func WithChromeDriverStartTime(duration time.Duration) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.startTimeout = duration
	}
}

func WithChromeDriverThreads(threads int) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.threads = threads
	}
}

func WithChromeDriverPort(port int) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.port = port
	}
}

func WithChromeDriverAdbPort(port int) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.adbPort = port
	}
}

func WithChromeDriverHost(host string) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.host = host
	}
}

func WithChromeDriverCommand(cmd string) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.path = cmd
	}
}

func WithChromeDriverBaseURL(baseURL string) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.baseURL = baseURL
	}
}

func WithChromeDriverLogLevel(logLevel LogLevelChromeDriver) ChromeDriverOption {
	return func(cd *chromeDriver) {
		cd.logLevel = logLevel
	}
}

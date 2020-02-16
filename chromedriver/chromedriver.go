package chromedriver

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

type LogLevel string

const (
	DefaultCommand   = "chromedriver"
	DefaultStartTime = time.Duration(20 * time.Second)
	DefaultThreads   = 6
	DefaultPort      = 9165
	DefaultHost      = "127.0.0.1"
	DefaultLogPath   = "chromedriver.log"

	// ALL, DEBUG, INFO, WARNING, SEVERE, OFF
	All     LogLevel = "ALL"
	Debug   LogLevel = "DEBUG"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Severe  LogLevel = "SEVERE"
	Off     LogLevel = "OFF"
)

var (
	ErrTimeOut                    = errors.New("chromedriver: timout error")
	ErrChromeDriverAlreadyRunning = errors.New("chromedriver: chromedriver already running")
	ErrChromeDriverNotRunning     = errors.New("chromedriver: chromedriver not running")
)

type ChromeDriver struct {
	// base URL path prefix for commands, e.g. wd/url
	baseURL string
	path    string
	threads int

	logPath  string
	logLevel LogLevel

	startTimeout time.Duration
	port         int
	adbPort      int
	host         string
	proc         *exec.Cmd
	out          bytes.Buffer
}

func New(opts ...Option) *ChromeDriver {
	cd := &ChromeDriver{
		path:         DefaultCommand,
		startTimeout: DefaultStartTime,
		threads:      DefaultThreads,
		port:         DefaultPort,
		host:         DefaultHost,
		logLevel:     All,
	}

	for _, opt := range opts {
		opt(cd)
	}

	return cd
}

func (d *ChromeDriver) Out() bytes.Buffer {
	return d.out
}

func (d *ChromeDriver) Run(ctx context.Context) error {
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
		d.logPath = home + DefaultLogPath
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

func (d *ChromeDriver) Stop(_ context.Context) error {
	if d.proc == nil {
		return ErrChromeDriverNotRunning
	}
	defer func() {
		d.proc = nil
	}()
	return d.proc.Process.Signal(os.Interrupt)
}

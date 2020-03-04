package chromedriver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
)

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

type LogLevel string
type Hook func(pid int)

type Error struct {
	Err    error
	Stdout string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v %s",
		e.Err, e.Stdout)
}

type ChromeDriver struct {
	// base URL path prefix for commands, e.g. wd/url
	baseURL string
	path    string
	threads int

	onStartHook Hook
	onCloseHook Hook

	logPath  string
	logLevel LogLevel

	startTimeout time.Duration
	port         int
	adbPort      int
	host         string
	cmd          *exec.Cmd
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

func (d *ChromeDriver) Run(ctx context.Context) error {
	if d.cmd != nil {
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

	var out bytes.Buffer
	d.cmd = exec.Command(d.path, args...)
	d.cmd.Stderr = &out
	if err := d.cmd.Start(); err != nil {
		return Error{
			Err:    err,
			Stdout: out.String(),
		}
	}

	pid := d.cmd.Process.Pid
	done := make(chan error, 1)
	go func() {
		if d.onStartHook != nil {
			d.onStartHook(pid)
		}
		done <- d.cmd.Wait()
	}()
	select {
	case <-ctx.Done():
		err := d.cmd.Process.Kill()
		if err == nil {
			return ErrTimeOut
		}
		return err
	case err := <-done:
		if err != nil {
			return Error{
				Err:    err,
				Stdout: out.String(),
			}
		}
		if d.onCloseHook != nil {
			d.onCloseHook(pid)
		}
	}
	return nil
}

func (d *ChromeDriver) Stop(_ context.Context) error {
	if d.cmd == nil {
		return ErrChromeDriverNotRunning
	}
	defer func() {
		d.cmd = nil
	}()
	return d.cmd.Process.Signal(os.Interrupt)
}

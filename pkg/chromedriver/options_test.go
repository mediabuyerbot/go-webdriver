package chromedriver

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"

	"github.com/stretchr/testify/assert"
)

func newProc(t *testing.T, opts ...Option) *Process {
	proc, err := New(bin.ChromeDriver64(), opts...)
	assert.Nil(t, err)
	return proc
}

func TestWithAdbPort(t *testing.T) {
	proc := newProc(t, WithAdbPort(1234))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--adb-port=1234", proc.Args()[0])

	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "out of range")
	}()
	newProc(t, WithAdbPort(-1))
}

func TestWithAppendLog(t *testing.T) {
	proc := newProc(t, WithAppendLog())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, argsAppendLog, proc.Args()[0])
}

func TestWithBaseURL(t *testing.T) {
	proc := newProc(t, WithBaseURL("url"))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--url-base=url", proc.Args()[0])
}

func TestWithLogLevel(t *testing.T) {
	proc := newProc(t, WithLogLevel(Info))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--log-level=INFO", proc.Args()[0])

	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	proc = newProc(t, WithLogLevel(LogLevel("HOHOH")))

}

func TestWithLogPath(t *testing.T) {
	logPath := os.TempDir()
	proc := newProc(t, WithLogPath(logPath))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, fmt.Sprintf("--log-path=%s", logPath), proc.Args()[0])

	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	proc = newProc(t, WithLogPath("/path/to/"))
}

func TestWithPort(t *testing.T) {
	proc := newProc(t, WithPort(123))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--port=123", proc.Args()[0])

	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "out of range")
	}()
	newProc(t, WithPort(999999999))
}

func TestWithReadableTimestamp(t *testing.T) {
	proc := newProc(t, WithReadableTimestamp())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--readable-timestamp", proc.Args()[0])
}

func TestWithReplayable(t *testing.T) {
	proc := newProc(t, WithReplayable())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--replayable", proc.Args()[0])
}

func TestWithSilent(t *testing.T) {
	proc := newProc(t, WithSilent())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--silent", proc.Args()[0])
}

func TestWithVerbose(t *testing.T) {
	proc := newProc(t, WithVerbose())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--verbose", proc.Args()[0])
}

func TestWithVersion(t *testing.T) {
	proc := newProc(t, WithShowVersion())
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--version", proc.Args()[0])
}

func TestWithWhitelistedIps(t *testing.T) {
	proc := newProc(t, WithWhitelistedIps([]string{"ip", "ip2"}))
	assert.Len(t, proc.Args(), 1)
	assert.Equal(t, "--whitelisted-ips=ip,ip2", proc.Args()[0])
}

func TestWithRunHook(t *testing.T) {
	proc := newProc(t, WithRunHook(func(pid int) {}))
	assert.NotNil(t, proc.runHook)
}

func TestWithStopHook(t *testing.T) {
	proc := newProc(t, WithStopHook(func(pid int) {}))
	assert.NotNil(t, proc.stopHook)
}

func TestWithStderr(t *testing.T) {
	proc := newProc(t, WithStderr(bytes.NewBuffer([]byte(``))))
	assert.NotNil(t, proc.stderr)
}

func TestWithStdout(t *testing.T) {
	proc := newProc(t, WithStdout(bytes.NewBuffer([]byte(``))))
	assert.NotNil(t, proc.stdout)
}

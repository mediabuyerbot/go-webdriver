package geckodriver

import (
	"fmt"
	"testing"

	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"
	"github.com/stretchr/testify/assert"
)

func newProc(t *testing.T, opts ...Option) *Process {
	proc, err := New(bin.GeckoDriver64(), opts...)
	assert.Nil(t, err)
	return proc
}

func TestOptions(t *testing.T) {
	fakeBinPath := bin.GeckoDriver64()

	proc := newProc(t,
		WithConnExisting(),
		WithJsDebugger(),
		WithVerbose(),
		WithShowVersion(),
		WithBinary(fakeBinPath),
		WithLogLevel(Debug),
		WithMarionetteHost("host"),
		WithMarionettePort(1234),
		WithHost("host"),
		WithPort(4444),
	)

	args := proc.Args()
	assert.True(t, findArg(args, flagConnExisting))
	assert.True(t, findArg(args, flagJsDebugger))
	assert.True(t, findArg(args, flagVersion))
	assert.True(t, findArg(args, flagVerbose))
	assert.True(t, findArg(args, fmt.Sprintf(argsBinary, fakeBinPath)))
	assert.True(t, findArg(args, fmt.Sprintf(argsLogLevel, Debug)))
	assert.True(t, findArg(args, fmt.Sprintf(argsMarionetteHost, "host")))
	assert.True(t, findArg(args, fmt.Sprintf(argsMarionettePort, 1234)))
	assert.True(t, findArg(args, fmt.Sprintf(argsHost, "host")))
	assert.True(t, findArg(args, fmt.Sprintf(argsPort, 4444)))
}

func TestOptions_SetBinaryPanic(t *testing.T) {
	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "no such file or directory")
	}()
	newProc(t, WithBinary("/path-to-not-exists-bin"))
}

func TestOptions_SetLogLevel(t *testing.T) {
	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "unknown log level")
	}()
	newProc(t, WithLogLevel(LogLevel("$$$")))
}

func TestOptions_SetMarionettePort(t *testing.T) {
	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "out of range")
	}()
	newProc(t, WithMarionettePort(-1))
}

func TestOptions_SetPort(t *testing.T) {
	defer func() {
		err := recover()
		assert.Error(t, err.(error))
		assert.Contains(t, err.(error).Error(), "out of range")
	}()
	newProc(t, WithPort(700000))
}

func findArg(args []string, arg string) bool {
	for _, a := range args {
		if arg == a {
			return true
		}
	}
	return false
}

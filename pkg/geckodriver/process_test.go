package geckodriver

import (
	"bytes"
	"context"
	"testing"
	"time"

	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"
	"github.com/mitchellh/go-ps"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	ctx := context.Background()
	driver, err := New(bin.GeckoDriver64(),
		WithLogLevel(Error),
		WithPort(4444),
		WithStdout(stdout),
		WithStderr(stderr),
		WithRunHook(func(pid int) {
			assert.NotZero(t, pid)
			pid = pid
		}),
		WithStopHook(func(pid int) {
			assert.NotZero(t, pid)
			assert.False(t, procIsRunning(pid))
		}),
	)
	assert.Nil(t, err)
	assert.NotNil(t, driver)

	go func() {
		time.Sleep(2 * time.Second)
		err = driver.Stop(ctx)
		assert.Nil(t, err)
		err = driver.Stop(ctx)
		assert.Nil(t, err)
	}()

	err = driver.Run(ctx)
	if err != nil {
		t.Log(stderr.String())
	}
	assert.Nil(t, err)
}

func TestProcessWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	driver, err := New(bin.GeckoDriver64(),
		WithPort(4998),
		WithRunHook(func(pid int) {
			assert.NotZero(t, pid)
		}),
		WithStopHook(func(pid int) {
			assert.NotZero(t, pid)
			assert.False(t, procIsRunning(pid))
		}),
	)
	assert.Nil(t, err)
	assert.NotNil(t, driver)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	err = driver.Run(ctx)
	assert.Nil(t, err)
}

func procIsRunning(pid int) bool {
	osp, err := ps.FindProcess(pid)
	if err != nil || osp == nil {
		return false
	}
	return osp.Pid() == pid
}

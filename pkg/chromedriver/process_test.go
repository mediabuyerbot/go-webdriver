package chromedriver

import (
	"context"
	"testing"
	"time"

	"github.com/mitchellh/go-ps"

	"github.com/stretchr/testify/assert"

	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"
)

func TestProcess(t *testing.T) {
	ctx := context.Background()
	driver, err := New(bin.ChromeDriver64(),
		WithPort(4987),
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
	assert.Nil(t, err)
}

func TestProcessWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	driver, err := New(bin.ChromeDriver64(),
		WithPort(4988),
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

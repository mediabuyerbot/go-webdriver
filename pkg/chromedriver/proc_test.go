package chromedriver

import (
	"context"
	"testing"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/stretchr/testify/assert"
)

func TestRunAndStop(t *testing.T) {
	if t.Skipped() {
		return
	}

	ctx := context.Background()
	onstart := make(chan int)
	driver := New(
		WithOnStartHook(func(pid int) {
			onstart <- pid
		}),
	)
	go func() {
		if err := driver.Run(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	pid := <-onstart

	time.Sleep(500 * time.Millisecond)
	proc, err := ps.FindProcess(pid)
	assert.Nil(t, err)
	assert.Equal(t, proc.Pid(), pid)

	time.Sleep(500 * time.Millisecond)

	if err := driver.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

package chromedriver

import (
	"context"
	"testing"

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
	proc, err := ps.FindProcess(pid)
	assert.Nil(t, err)
	assert.Equal(t, proc.Pid(), pid)

	if err := driver.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

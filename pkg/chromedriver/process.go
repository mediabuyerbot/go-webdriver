package chromedriver

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Hook func(pid int)

type Process struct {
	args *arguments
	bin  string
	cmd  *exec.Cmd

	stdout io.Writer
	stderr io.Writer

	runHook  Hook
	stopHook Hook

	lock sync.RWMutex
}

// New returns a new instance of Process.
func New(bin string, opts ...Option) (*Process, error) {
	driverPath, err := exec.LookPath(bin)
	if err != nil {
		return nil, err
	}

	cd := &Process{
		args: newArguments(),
		bin:  driverPath,
	}

	for _, opt := range opts {
		opt(cd)
	}

	return cd, nil
}

func (d *Process) Run(ctx context.Context) error {
	d.lock.Lock()
	d.cmd = exec.CommandContext(ctx, d.bin, d.args.build()...)

	d.setWriters()

	if err := d.cmd.Start(); err != nil {
		d.lock.Unlock()
		return err
	}
	d.lock.Unlock()

	if d.runHook != nil {
		go func() {
			time.Sleep(time.Second)
			d.runHook(d.cmd.Process.Pid)
		}()
	}

	err := d.cmd.Wait()
	if err != nil && d.isSignaled() {
		err = nil
	}

	if d.stopHook != nil {
		d.stopHook(d.cmd.Process.Pid)
	}

	d.lock.Lock()
	d.cmd = nil
	d.lock.Unlock()

	return err
}

func (d *Process) Stop(_ context.Context) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.cmd == nil {
		return nil
	}
	return d.cmd.Process.Signal(os.Interrupt)
}

func (d *Process) setWriters() {
	if d.stderr != nil {
		d.cmd.Stderr = d.stderr
	}
	if d.stdout != nil {
		d.cmd.Stdout = d.stdout
	}
}

func (d *Process) isSignaled() bool {
	if d.cmd == nil {
		return false
	}
	status, ok := d.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		return false
	}
	return status.ExitStatus() == -1 && status.Signaled()
}

func (d *Process) Args() []string {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.args.build()
}

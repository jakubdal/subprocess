package subprocess

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Process is used to manipulate underlying, running os process
//
// Process is not thread-safe
type Process struct {
	ctx     context.Context
	errChan chan<- error

	name          string
	additionalEnv []string
	args          []string

	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer

	stopped bool
}

// NewProcess starts a new process with context `ctx`, running program `name`
// with current environment, with `additionalEnv` variables added to it (if
// they're not empty), using provided `args` as program arguments.
//
// If program exits with non-zero exit code, that was not caused by a call to
// Process' `Stop` function, it will be restarted. All errors from those
// restarts will be sent to `errChan` if it's non-nil
func NewProcess(
	ctx context.Context,
	errChan chan<- error,

	name string,
	additionalEnv []string,
	args ...string,
) (*Process, error) {
	p := &Process{
		ctx:     ctx,
		errChan: errChan,

		name:          name,
		additionalEnv: additionalEnv,
		args:          args,
	}

	if err := p.start(); err != nil {
		return nil, errors.Wrap(err, "start process")
	}

	go func() {
		for {
			if err := p.cmd.Wait(); err == nil {
				return
			}

			if p.stopped {
				return
			}

			if err := p.start(); err != nil {
				if p.errChan != nil {
					p.errChan <- errors.Wrap(err, "restart process")
				}
			}
		}
	}()

	return p, nil
}

func (p *Process) start() error {
	p.cmd = exec.Command(p.name, p.args...)
	p.cmd.Env = append(os.Environ(), p.additionalEnv...)
	if err := p.cmd.Start(); err != nil {
		return errors.Wrap(err, "cmd.Start")
	}

	return nil
}

// Stop stops underlying process from running
func (p *Process) Stop() {
	p.stopped = true
	p.cmd.Process.Kill()
}

// Signal relays provided signal to underlying os process
func (p *Process) Signal(sig os.Signal) error {
	return errors.Wrap(p.cmd.Process.Signal(sig), "p.cmd.Process.Signal")
}

// ReadStdout reads data from stdout into provided `buf`
func (p *Process) ReadStdout(buf []byte) (int, error) {
	pipe, err := p.cmd.StdoutPipe()
	if err != nil {
		return 0, errors.Wrap(err, "cmd.StdoutPipe")
	}
	return pipe.Read(buf)
}

// ReadStderr reads data from stderr into provided `buf`
func (p *Process) ReadStderr(buf []byte) (int, error) {
	pipe, err := p.cmd.StderrPipe()
	if err != nil {
		return 0, errors.Wrap(err, "cmd.StderrPipe")
	}
	return pipe.Read(buf)
}

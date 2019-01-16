package subprocess

import (
	"context"
	"io"
	"os"
)

type process interface {
	Start(ctx context.Context) error
	Stdout() io.Reader
	Stderr() io.Reader
	Wait() error
	Stop()
	Signal(sig os.Signal) error
}

type RestartingProcess struct {
	proc process
}

// WithRestartOnFailure wraps underlying process to ensure it will exit with a
// zero status code
func WithRestartOnFailure(proc process) *RestartingProcess {
	return &RestartingProcess{proc: proc}
}

func (p *RestartingProcess) Start(ctx context.Context) error {
	// TODO restart until wait returns zero
	return p.proc.Start(ctx)
}

func (p *RestartingProcess) Stdout() io.Reader {
	return p.proc.Stdout()
}

func (p *RestartingProcess) Stderr() io.Reader {
	return p.proc.Stderr()
}

func (p *RestartingProcess) Wait() error {
	// TODO wait on this process to zero-exit
	return p.proc.Wait()
}

func (p *RestartingProcess) Stop() {
	// TODO break this process' restarting loop
	p.proc.Stop()
}

func (p *RestartingProcess) Signal(sig os.Signal) error {
	return p.proc.Signal(sig)
}

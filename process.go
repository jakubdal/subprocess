package subprocess

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// ProcessOpts contains argument and environment variable options for
// starting a process
type ProcessOpts struct {
	AdditionalEnv []string
	CLIArgs       []string
}

// Process is used to manipulate underlying, running os process
//
// Process is not thread-safe
type Process struct {
	ctx            context.Context
	name           string
	processOpts    ProcessOpts
	descriptorOpts DescriptorOpts

	cmd *exec.Cmd
}

// NewProcess starts a new process with context `ctx`, running program `name`
func NewProcess(
	ctx context.Context,
	name string,

	processOpts *ProcessOpts,
	descriptorOpts *DescriptorOpts,
) (*Process, error) {
	p := &Process{
		ctx:  ctx,
		name: name,
	}
	if processOpts != nil {
		p.processOpts = *processOpts
	}
	if descriptorOpts != nil {
		p.descriptorOpts = *descriptorOpts
	}

	return p, nil
}

func (p *Process) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = p.ctx
	}
	p.cmd = exec.CommandContext(ctx, p.name, p.processOpts.CLIArgs...)
	p.cmd.Env = append(os.Environ(), p.processOpts.AdditionalEnv...)
	p.descriptorOpts.SetDescriptors(p.cmd)

	if err := p.cmd.Start(); err != nil {
		return errors.Wrap(err, "cmd.Start")
	}

	return nil
}

// Stdout returns reader of stdout
func (p *Process) Stdout() io.Reader {
	if reader, ok := p.cmd.Stdout.(io.Reader); ok {
		return reader
	}
	return nil
}

// Stderr returns reader of stderr
func (p *Process) Stderr() io.Reader {
	if reader, ok := p.cmd.Stderr.(io.Reader); ok {
		return reader
	}
	return nil
}

func (p *Process) Wait() error {
	if err := p.cmd.Wait(); err != nil {
		errors.Wrap(err, "cmd.Wait")
	}
	return nil
}

// Stop stops underlying process from running
func (p *Process) Stop() {
	p.cmd.Process.Kill()
}

// Signal relays provided signal to underlying os process
func (p *Process) Signal(sig os.Signal) error {
	return errors.Wrap(p.cmd.Process.Signal(sig), "p.cmd.Process.Signal")
}

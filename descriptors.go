package subprocess

import (
	"io"
	"os/exec"
)

// DescriptorOpts contains options for substituting process' i/o descriptors
//
// Any `nil` descriptor will not overwrite cmd's default descriptors
type DescriptorOpts struct {
	// Stdin will be set at Stdin in started process
	Stdin io.Reader
	// Stdout will be set as Stdout in started process
	Stdout io.Writer
	// Stderr will be set as Stderr in started process
	Stderr io.Writer
}

func (opts *DescriptorOpts) SetDescriptors(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	}
	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	}
}

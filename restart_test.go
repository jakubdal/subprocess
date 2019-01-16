package subprocess

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRestartOnCrash(t *testing.T) {
	testDir := filepath.Join("e2e_testsuite", "crash_test")

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(
			testDir,
			os.ModePerm,
		); err != nil {
			t.Fatalf("os.Mkdir: %v", err)
		}
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Fatalf("os.Remove: %v", err)
		}
	}()

	proc, err := NewProcess(
		context.Background(),
		"python3",
		&ProcessOpts{CLIArgs: []string{
			filepath.Join("e2e_testsuite", "crash_three_times.py"),
		}},
		&DescriptorOpts{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
	)
	assertNil(t, err)
	restarter := WithRestartOnFailure(proc)

	assertNil(t, restarter.Start(nil))

	assertNil(t, restarter.Wait())

	stdout := restarter.Stdout().(*bytes.Buffer).String()
	if stdout != "success" {
		t.Errorf("invalid output: %v", stdout)
	}
}

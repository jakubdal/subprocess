package subprocess

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTouchFile(t *testing.T) {
	testDir, touchedFile := "create_list_test", "test_touch_file"

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, os.ModePerm); err != nil {
			t.Fatalf("os.Mkdir: %v", err)
		}
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Fatalf("os.Remove: %v", err)
		}
	}()

	toucher, err := NewProcess(
		context.Background(),
		"touch",
		&ProcessOpts{
			CLIArgs: []string{filepath.Join(testDir, touchedFile)},
		},
		nil,
	)
	if err != nil {
		t.Fatalf("create toucher process")
	}

	toucher.Start(nil)
	toucher.Wait()

	if _, err := os.Stat(
		filepath.Join(testDir, touchedFile),
	); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}

func TestPrintStderrStdout(t *testing.T) {
	printer, err := NewProcess(
		context.Background(),
		"python3",
		&ProcessOpts{CLIArgs: []string{
			filepath.Join("e2e_testsuite", "print.py"),
		}},
		&DescriptorOpts{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
	)
	assertNil(t, err)

	assertNil(t, printer.Start(nil))
	printer.Wait()

	stdoutContent := printer.Stdout().(*bytes.Buffer).String()
	if stdoutContent != "Hello, stdout!" {
		t.Errorf("invalid stdout content: %v", stdoutContent)
	}

	stderrContent := printer.Stderr().(*bytes.Buffer).String()
	if stderrContent != "Hello, stderr!" {
		t.Errorf("invalid stderr content: %v", stderrContent)
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

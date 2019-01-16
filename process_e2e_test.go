package subprocess

import (
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

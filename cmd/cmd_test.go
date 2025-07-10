package cmd

// These tests exercise the CLI commands defined in this package.  They verify
// that command handlers invoke the underlying packages correctly and handle
// common error scenarios.

import (
	"os"
	"testing"

	"codex/memory"
)

// TestAddCommand verifies that the `add` CLI subcommand writes a memory entry
// with the provided importance score. It exercises the integration between the
// cmd layer and the memory package.
func TestAddCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	addCmd.Flags().Set("importance", "4")
	addCmd.Run(addCmd, []string{"proj", "user", "hello"})

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	entries, err := memory.LastNEntries(db, "proj", 1)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(entries) != 1 || entries[0].Content != "hello" || entries[0].Importance != 4 {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

// TestExecuteInvalidCommand ensures that invoking an unknown command returns an
// error from cobra. This guards the command parsing entry point.
func TestExecuteInvalidCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"nonexist"})
	if err := Execute(); err == nil {
		t.Fatalf("expected error for invalid command")
	}
}

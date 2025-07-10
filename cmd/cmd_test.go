package cmd

import (
	"os"
	"testing"

	"local-ai/memory"
)

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

func TestExecuteInvalidCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"nonexist"})
	if err := Execute(); err == nil {
		t.Fatalf("expected error for invalid command")
	}
}

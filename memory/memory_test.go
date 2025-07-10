package memory

import (
	"os"
	"testing"
	"time"
)

func TestAddAndRetrieveEntries(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	db, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	if err := AddEntry(db, "proj", "system", "one"); err != nil {
		t.Fatalf("AddEntry error: %v", err)
	}
	time.Sleep(1 * time.Second)
	if err := AddEntry(db, "proj", "user", "two"); err != nil {
		t.Fatalf("AddEntry error: %v", err)
	}
	time.Sleep(1 * time.Second)
	if err := AddEntry(db, "proj", "assistant", "three", 5); err != nil {
		t.Fatalf("AddEntry error: %v", err)
	}

	last, err := LastNEntries(db, "proj", 2)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(last) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(last))
	}
	if last[0].Content != "three" || last[1].Content != "two" {
		t.Fatalf("unexpected order: %+v", last)
	}

	top, err := TopImportantEntries(db, "proj", 1)
	if err != nil {
		t.Fatalf("TopImportantEntries error: %v", err)
	}
	if len(top) != 1 || top[0].Content != "three" {
		t.Fatalf("unexpected top entry: %+v", top)
	}
}

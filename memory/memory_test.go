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

func TestAddMemory(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	if err := AddMemory("proj", "user", "hello"); err != nil {
		t.Fatalf("AddMemory error: %v", err)
	}

	db, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	entries, err := LastNEntries(db, "proj", 1)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(entries) != 1 || entries[0].Content != "hello" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestProjectManagement(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	db, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	if err := AddProject(db, "p1"); err != nil {
		t.Fatalf("add project: %v", err)
	}
	if err := AddProject(db, "p2"); err != nil {
		t.Fatalf("add project: %v", err)
	}
	if err := SetActiveProject(db, "p2"); err != nil {
		t.Fatalf("set active: %v", err)
	}
	list, err := ListProjects(db)
	if err != nil || len(list) != 2 {
		t.Fatalf("list err: %v len:%d", err, len(list))
	}
	active, err := GetActiveProject(db)
	if err != nil || active != "p2" {
		t.Fatalf("active err:%v val:%s", err, active)
	}
	if err := DeleteProject(db, "p2"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	active, _ = GetActiveProject(db)
	if active != "" {
		t.Fatalf("expected no active project")
	}
}

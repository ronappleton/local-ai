package memory

import (
	"os"
	"testing"
	"time"
)

// TestAddAndRetrieveEntries ensures entries can be stored and queried in time
// order and by importance. It validates the core persistence logic used by the
// assistant for long term memory.
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

// TestAddMemory covers the higher-level AddMemory helper which opens the
// database internally. AI Awareness: this path is used when the CLI records a
// message without directly handling the database handle.
func TestAddMemory(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	if err := AddMemory("proj", "user", "hello", 2); err != nil {
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
	if len(entries) != 1 || entries[0].Content != "hello" || entries[0].Importance != 2 {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

// TestProjectManagement verifies creating, listing and deleting projects along
// with storing which one is active. This mimics how the assistant maintains
// separate conversation contexts.
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

// TestRenameProject checks that renaming a project cascades to memories and the
// active project setting. This protects data consistency across the database.
func TestRenameProject(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	db, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	if err := AddProject(db, "old"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := SetActiveProject(db, "old"); err != nil {
		t.Fatalf("set active: %v", err)
	}
	if err := AddEntry(db, "old", "user", "hello"); err != nil {
		t.Fatalf("add entry: %v", err)
	}

	if err := RenameProject(db, "old", "new"); err != nil {
		t.Fatalf("rename: %v", err)
	}

	list, _ := ListProjects(db)
	if len(list) != 1 || list[0] != "new" {
		t.Fatalf("unexpected list: %+v", list)
	}
	active, _ := GetActiveProject(db)
	if active != "new" {
		t.Fatalf("active not updated: %s", active)
	}
	entries, _ := LastNEntries(db, "new", 1)
	if len(entries) != 1 || entries[0].Project != "new" {
		t.Fatalf("entries not renamed: %+v", entries)
	}
}

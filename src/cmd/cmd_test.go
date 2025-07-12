package cmd

// These tests exercise the CLI commands defined in this package.  They verify
// that command handlers invoke the underlying packages correctly and handle
// common error scenarios.

import (
	"codex/src/auth"
	"codex/src/memory"
	"os"
	"testing"
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
func TestCreateUserCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	if err := createUserCmd.RunE(createUserCmd, []string{"bob", "b@c.com", "pwd"}); err != nil {
		t.Fatalf("command error: %v", err)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()
	u, err := auth.GetByUsername(db, "bob")
	if err != nil || u == nil {
		t.Fatalf("user not created: %v", err)
	}
}

func TestPromoteUserCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	// create user first
	if err := createUserCmd.RunE(createUserCmd, []string{"eve", "e@c.com", "pwd"}); err != nil {
		t.Fatalf("create error: %v", err)
	}
	// promote
	if err := promoteUserCmd.RunE(promoteUserCmd, []string{"eve"}); err != nil {
		t.Fatalf("promote error: %v", err)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()
	u, err := auth.GetByUsername(db, "eve")
	if err != nil || !u.Admin || !u.Verified {
		t.Fatalf("user not promoted/verified: %+v err:%v", u, err)
	}
}

func TestListUsersCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	// create a couple users
	if err := createUserCmd.RunE(createUserCmd, []string{"john", "j@c.com", "pwd"}); err != nil {
		t.Fatalf("create error: %v", err)
	}
	if err := createUserCmd.RunE(createUserCmd, []string{"kate", "k@c.com", "pwd"}); err != nil {
		t.Fatalf("create error: %v", err)
	}

	if err := listUsersCmd.RunE(listUsersCmd, []string{}); err != nil {
		t.Fatalf("list error: %v", err)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()
	users, err := auth.List(db)
	if err != nil || len(users) != 2 {
		t.Fatalf("unexpected list: %+v err:%v", users, err)
	}
}

func TestDeleteUserCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	// create users
	createUserCmd.RunE(createUserCmd, []string{"john", "j@c.com", "pwd"})
	createUserCmd.RunE(createUserCmd, []string{"kate", "k@c.com", "pwd"})

	if err := deleteUserCmd.RunE(deleteUserCmd, []string{"john"}); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()
	users, _ := auth.List(db)
	if len(users) != 1 || users[0].Username != "kate" {
		t.Fatalf("unexpected users: %+v", users)
	}
}

func TestDeleteAllUsersCommand(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	createUserCmd.RunE(createUserCmd, []string{"john", "j@c.com", "pwd"})
	createUserCmd.RunE(createUserCmd, []string{"kate", "k@c.com", "pwd"})

	deleteUserCmd.Flags().Set("all", "true")
	err := deleteUserCmd.RunE(deleteUserCmd, []string{})
	deleteUserCmd.Flags().Set("all", "false")
	if err != nil {
		t.Fatalf("delete all error: %v", err)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()
	users, _ := auth.List(db)
	if len(users) != 0 {
		t.Fatalf("users not removed: %+v", users)
	}
}

package auth

import (
	"codex/src/memory"
	"os"
	"testing"
)

// TestSetAdminMarksVerified ensures that promoting a user also verifies them.
func TestSetAdminMarksVerified(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(cwd)

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer db.Close()

	if err := CreateUser(db, "bob", "b@c.com", "pwd"); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := SetAdmin(db, "bob", true); err != nil {
		t.Fatalf("set admin: %v", err)
	}

	u, err := GetByUsername(db, "bob")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !u.Admin || !u.Verified {
		t.Fatalf("user not admin/verified: %+v", u)
	}
}

package memory

// Package memory provides persistence for conversations and project metadata.
// It acts as the long term storage that the AI can query or append to when new
// messages are generated. The implementation uses SQLite for a lightweight
// local database.

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// MemoryEntry represents a single line of conversation stored in the database.
// Importance is an optional ranking that can be used by the AI to prioritise
// context when generating responses.
type MemoryEntry struct {
	ID         int
	Project    string
	Role       string
	Content    string
	Timestamp  time.Time
	Importance int
}

// InitDB opens the SQLite database stored in memory.db in the current working
// directory and ensures all required tables exist. It returns a handle to the
// database which callers must close. This function is used throughout the
// project whenever persistent storage is required.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "memory.db")
	if err != nil {
		return nil, err
	}
	createMemory := `CREATE TABLE IF NOT EXISTS memory (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        project TEXT,
        role TEXT,
        content TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        importance INTEGER DEFAULT 0
    );`
	if _, err := db.Exec(createMemory); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (name TEXT PRIMARY KEY);`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT);`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
               id INTEGER PRIMARY KEY AUTOINCREMENT,
               username TEXT UNIQUE,
               email TEXT UNIQUE,
               password TEXT,
               verified INTEGER DEFAULT 0,
               totp_secret TEXT,
               admin INTEGER DEFAULT 0
       );`); err != nil {
		db.Close()
		return nil, err
	}
	// attempt to add admin column if the table already existed without it
	db.Exec(`ALTER TABLE users ADD COLUMN admin INTEGER DEFAULT 0`)
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
               token TEXT PRIMARY KEY,
               user_id INTEGER,
               type TEXT,
               expires DATETIME
       );`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS model_cache (
               id TEXT PRIMARY KEY,
               pipeline TEXT,
               last_modified TEXT,
               downloads INTEGER,
               tags TEXT,
               sha TEXT,
               files TEXT,
               llama_compatible INTEGER DEFAULT 0,
               model_type TEXT,
               hidden_size INTEGER,
               n_layer INTEGER,
               num_attention_heads INTEGER,
               quantized INTEGER DEFAULT 0,
               gguf INTEGER DEFAULT 0,
               safetensors INTEGER DEFAULT 0,
               compatible_backends TEXT,
               license TEXT,
               model_card TEXT,
               download_size INTEGER
       );`); err != nil {
		db.Close()
		return nil, err
	}
	// attempt to add columns for new metadata if the table existed without them
	db.Exec(`ALTER TABLE model_cache ADD COLUMN llama_compatible INTEGER DEFAULT 0`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN model_type TEXT`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN hidden_size INTEGER`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN n_layer INTEGER`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN num_attention_heads INTEGER`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN quantized INTEGER DEFAULT 0`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN gguf INTEGER DEFAULT 0`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN safetensors INTEGER DEFAULT 0`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN compatible_backends TEXT`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN license TEXT`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN model_card TEXT`)
	db.Exec(`ALTER TABLE model_cache ADD COLUMN download_size INTEGER`)
	return db, nil
}

// AddEntry inserts a new memory record into the provided database connection.
// Importance is optional and defaults to zero. Higher importance can be used by
// the AI to prioritise which memories to surface during context gathering.
// AddEntry inserts a single memory row. Extension Point: additional columns such
// as embeddings could be stored here for advanced retrieval strategies.
func AddEntry(db *sql.DB, project, role, content string, importance ...int) error {
	imp := 0
	if len(importance) > 0 {
		imp = importance[0]
	}
	stmt, err := db.Prepare("INSERT INTO memory(project, role, content, importance) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, role, content, imp)
	return err
}

// LastNEntries retrieves the most recent `n` memories for the given project
// ordered by descending timestamp. This is used when the AI wants to recall the
// latest conversation context. Extension Point: filtering by role or date range
// could be added here.
func LastNEntries(db *sql.DB, project string, n int) ([]MemoryEntry, error) {
	stmt, err := db.Prepare(`SELECT id, project, role, content, timestamp, importance FROM memory WHERE project = ? ORDER BY timestamp DESC LIMIT ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []MemoryEntry
	for rows.Next() {
		var e MemoryEntry
		if err := rows.Scan(&e.ID, &e.Project, &e.Role, &e.Content, &e.Timestamp, &e.Importance); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// TopImportantEntries fetches the `n` highest ranked memories for a project.
// Entries are sorted primarily by Importance. Extension Point: this function
// could incorporate vector similarity metrics for more intelligent recall.
func TopImportantEntries(db *sql.DB, project string, n int) ([]MemoryEntry, error) {
	stmt, err := db.Prepare(`SELECT id, project, role, content, timestamp, importance FROM memory WHERE project = ? ORDER BY importance DESC, timestamp DESC LIMIT ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []MemoryEntry
	for rows.Next() {
		var e MemoryEntry
		if err := rows.Scan(&e.ID, &e.Project, &e.Role, &e.Content, &e.Timestamp, &e.Importance); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// AddProject inserts a project name into the projects table if it does not
// already exist. Projects segment stored conversations so multiple contexts can
// be maintained. Extension Point: project-level metadata could be stored in a
// dedicated table.
func AddProject(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO projects(name) VALUES(?)`, name)
	return err
}

// DeleteProject removes a project from the projects table and clears related
// settings such as the active project. Memory records themselves are not
// removed which allows historical data inspection if needed. Extension Point:
// consider cascading deletes if permanent removal is desired.
func DeleteProject(db *sql.DB, name string) error {
	if _, err := db.Exec(`DELETE FROM projects WHERE name = ?`, name); err != nil {
		return err
	}
	// clear active project if it was deleted
	active, err := GetActiveProject(db)
	if err != nil {
		return err
	}
	if active == name {
		_, err = db.Exec(`DELETE FROM settings WHERE key = 'active_project'`)
	}
	return err
}

// RenameProject updates all references when a project changes name. Both the
// projects table and existing memory entries are updated in a single
// transaction. The active project setting is also adjusted if needed.
func RenameProject(db *sql.DB, oldName, newName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE projects SET name = ? WHERE name = ?`, newName, oldName); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`UPDATE memory SET project = ? WHERE project = ?`, newName, oldName); err != nil {
		tx.Rollback()
		return err
	}
	var active string
	tx.QueryRow(`SELECT value FROM settings WHERE key = 'active_project'`).Scan(&active)
	if active == oldName {
		if _, err := tx.Exec(`UPDATE settings SET value = ? WHERE key = 'active_project'`, newName); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ListProjects returns all project names sorted alphabetically. Extension
// Point: metadata such as project creation time could be returned to aid UI
// clients
func ListProjects(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM projects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// SetActiveProject updates the settings table so that the supplied project is
// considered the current context. Subsequent memory operations will default to
// this project. Extension Point: switching projects could trigger hooks to
// reload cached context or embeddings.
func SetActiveProject(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO settings(key, value) VALUES('active_project', ?)`, name)
	return err
}

// GetActiveProject returns the project currently marked as active. An empty
// string indicates that no active project is set. AI Awareness: the active
// project controls which memory entries are surfaced for context.
func GetActiveProject(db *sql.DB) (string, error) {
	row := db.QueryRow(`SELECT value FROM settings WHERE key = 'active_project'`)
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return name, err
}

// AddMemory is a convenience wrapper that opens the default database file,
// records a memory entry and closes the connection. It is primarily used by
// higher level interfaces such as the CLI. Extension Point: callers could
// supply their own *sql.DB to reuse connections or implement transactional
// batching.
func AddMemory(project, role, content string, importance ...int) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return AddEntry(db, project, role, content, importance...)
}

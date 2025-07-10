package memory

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type MemoryEntry struct {
	ID         int
	Project    string
	Role       string
	Content    string
	Timestamp  time.Time
	Importance int
}

// InitDB opens the SQLite database stored in memory.db in the current working directory
// and creates the memory table if it does not already exist.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "memory.db")
	if err != nil {
		return nil, err
	}
	createTable := `CREATE TABLE IF NOT EXISTS memory (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        project TEXT,
        role TEXT,
        content TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        importance INTEGER DEFAULT 0
    );`
	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// AddEntry inserts a new memory record into the database.
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

// LastNEntries retrieves the last n memory records for a project ordered by timestamp descending.
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

// TopImportantEntries retrieves the top n most important records for a project ordered by importance descending.
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

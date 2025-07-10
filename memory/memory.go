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

// AddProject inserts a project name if it doesn't already exist.
func AddProject(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO projects(name) VALUES(?)`, name)
	return err
}

// DeleteProject removes a project from the table and associated settings.
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

// ListProjects returns all project names sorted alphabetically.
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

// SetActiveProject marks the given project as the active one.
func SetActiveProject(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO settings(key, value) VALUES('active_project', ?)`, name)
	return err
}

// GetActiveProject retrieves the current active project name.
func GetActiveProject(db *sql.DB) (string, error) {
	row := db.QueryRow(`SELECT value FROM settings WHERE key = 'active_project'`)
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return name, err
}

// AddMemory opens the default DB and stores a memory entry.
func AddMemory(project, role, content string) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return AddEntry(db, project, role, content)
}

//// AddMemory is a helper that opens the database, adds the entry, and closes the
//// connection. It is convenient for simple use cases like the CLI command.
//func AddMemory(project, role, content string, importance ...int) error {
//	db, err := InitDB()
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//	return AddEntry(db, project, role, content, importance...)
//}

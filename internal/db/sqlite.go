package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and initializes schema
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{conn: conn}

	// Initialize schema
	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// migrate creates the necessary tables
func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS processed_commits (
		hash TEXT PRIMARY KEY,
		processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS analysis_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		commit_hash TEXT NOT NULL,
		date TIMESTAMP NOT NULL,
		project TEXT NOT NULL,
		category TEXT NOT NULL,
		impact_summary TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (commit_hash) REFERENCES processed_commits(hash)
	);

	CREATE TABLE IF NOT EXISTS token_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		batch_id TEXT NOT NULL,
		input_tokens INTEGER NOT NULL,
		output_tokens INTEGER NOT NULL,
		cost_estimate REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_results_date ON analysis_results(date);
	CREATE INDEX IF NOT EXISTS idx_results_project ON analysis_results(project);
	CREATE INDEX IF NOT EXISTS idx_results_category ON analysis_results(category);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

package store

import (
	"fmt"
	"sync"

	log "github.com/go-pkgz/lgr"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // sqlite driver
)

// DB implements Store interface using SQLite
type DB struct {
	db *sqlx.DB
	mu sync.RWMutex
}

// New creates a new DB store with the given database path
func New(dbPath string) (*DB, error) {
	db, err := connectSQLite(dbPath)
	if err != nil {
		return nil, err
	}

	store := &DB{db: db}

	if err := store.createSchema(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	log.Printf("[DEBUG] initialized sqlite store at %s", dbPath)
	return store, nil
}

// connectSQLite establishes SQLite connection with pragmas
func connectSQLite(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	// set pragmas for performance and reliability
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=1000",
		"PRAGMA foreign_keys=ON",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}

	// limit connections for SQLite (single writer)
	db.SetMaxOpenConns(1)

	return db, nil
}

// createSchema creates the database tables if they don't exist
func (s *DB) createSchema() error {
	schema := `
		-- Providers
		CREATE TABLE IF NOT EXISTS providers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Accounts
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_id INTEGER NOT NULL REFERENCES providers(id),
			name TEXT NOT NULL,
			login TEXT,
			api_key TEXT,
			account_type TEXT NOT NULL DEFAULT 'cloud',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider_id, name)
		);

		-- Servers
		CREATE TABLE IF NOT EXISTS servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL REFERENCES accounts(id),
			name TEXT NOT NULL,
			ip TEXT,
			description TEXT,
			responsible TEXT,
			approximate_cost REAL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			server_type TEXT NOT NULL DEFAULT 'cloud',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Server Logs
		CREATE TABLE IF NOT EXISTS server_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
			action TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Indexes
		CREATE INDEX IF NOT EXISTS idx_accounts_provider ON accounts(provider_id);
		CREATE INDEX IF NOT EXISTS idx_servers_account ON servers(account_id);
		CREATE INDEX IF NOT EXISTS idx_servers_status ON servers(status);
		CREATE INDEX IF NOT EXISTS idx_server_logs_server ON server_logs(server_id);
		CREATE INDEX IF NOT EXISTS idx_server_logs_created ON server_logs(created_at);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

// Close closes the database connection
func (s *DB) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

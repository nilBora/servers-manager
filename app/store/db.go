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

	if err := store.seedDefaultProviders(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to seed providers: %w", err)
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

// seedDefaultProviders adds default providers if none exist
func (s *DB) seedDefaultProviders() error {
	// Check if providers already exist
	var count int
	if err := s.db.Get(&count, "SELECT COUNT(*) FROM providers"); err != nil {
		return fmt.Errorf("failed to count providers: %w", err)
	}

	if count > 0 {
		return nil // Already seeded
	}

	// Default providers
	providers := []struct {
		Name        string
		Description string
	}{
		{"Hetzner", "Hetzner Cloud and Dedicated Servers"},
		{"AWS", "Amazon Web Services"},
		{"Scaleway", "Scaleway Cloud Platform"},
		{"Vsys Host", "Vsys Hosting Services"},
	}

	for _, p := range providers {
		_, err := s.db.Exec(
			"INSERT INTO providers (name, description) VALUES (?, ?)",
			p.Name, p.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert provider %s: %w", p.Name, err)
		}
	}

	log.Printf("[INFO] seeded %d default providers", len(providers))
	return nil
}

// Close closes the database connection
func (s *DB) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

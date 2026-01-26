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

	if err := store.runMigrations(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
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
			ident TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Accounts
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_id INTEGER NOT NULL REFERENCES providers(id),
			group_name TEXT DEFAULT '',
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
			location TEXT DEFAULT '',
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

// runMigrations runs database migrations for schema updates
func (s *DB) runMigrations() error {
	// Migration: Add group_name column to accounts if it doesn't exist
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM pragma_table_info('accounts') WHERE name='group_name'`)
	if err != nil {
		return fmt.Errorf("failed to check accounts schema: %w", err)
	}
	if count == 0 {
		_, err := s.db.Exec(`ALTER TABLE accounts ADD COLUMN group_name TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add group_name column: %w", err)
		}
		log.Printf("[INFO] migration: added group_name column to accounts")
	}

	// Migration: Add location column to servers if it doesn't exist
	err = s.db.Get(&count, `SELECT COUNT(*) FROM pragma_table_info('servers') WHERE name='location'`)
	if err != nil {
		return fmt.Errorf("failed to check servers schema: %w", err)
	}
	if count == 0 {
		_, err := s.db.Exec(`ALTER TABLE servers ADD COLUMN location TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add location column: %w", err)
		}
		log.Printf("[INFO] migration: added location column to servers")
	}

	// Migration: Add ident column to providers if it doesn't exist
	err = s.db.Get(&count, `SELECT COUNT(*) FROM pragma_table_info('providers') WHERE name='ident'`)
	if err != nil {
		return fmt.Errorf("failed to check providers schema: %w", err)
	}
	if count == 0 {
		_, err := s.db.Exec(`ALTER TABLE providers ADD COLUMN ident TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add ident column: %w", err)
		}
		// Set default idents for existing providers
		_, _ = s.db.Exec(`UPDATE providers SET ident = 'hetzner_cloud' WHERE name = 'Hetzner'`)
		_, _ = s.db.Exec(`UPDATE providers SET ident = 'aws' WHERE name = 'AWS'`)
		_, _ = s.db.Exec(`UPDATE providers SET ident = 'scaleway' WHERE name = 'Scaleway'`)
		_, _ = s.db.Exec(`UPDATE providers SET ident = 'vsys_host' WHERE name = 'Vsys Host'`)
		log.Printf("[INFO] migration: added ident column to providers")

		// Add Hetzner Robot as new provider if Hetzner Cloud exists
		var hetznerExists int
		_ = s.db.Get(&hetznerExists, `SELECT COUNT(*) FROM providers WHERE ident = 'hetzner_cloud'`)
		if hetznerExists > 0 {
			_, _ = s.db.Exec(`INSERT INTO providers (ident, name, description) VALUES ('hetzner_robot', 'Hetzner Robot', 'Hetzner Dedicated Servers')`)
			log.Printf("[INFO] migration: added Hetzner Robot provider")
		}
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
		Ident       string
		Name        string
		Description string
	}{
		{"hetzner_cloud", "Hetzner Cloud", "Hetzner Cloud Servers"},
		{"hetzner_robot", "Hetzner Robot", "Hetzner Dedicated Servers"},
		{"aws", "AWS", "Amazon Web Services"},
		{"scaleway", "Scaleway", "Scaleway Cloud Platform"},
		{"vsys_host", "Vsys Host", "Vsys Hosting Services"},
	}

	for _, p := range providers {
		_, err := s.db.Exec(
			"INSERT INTO providers (ident, name, description) VALUES (?, ?, ?)",
			p.Ident, p.Name, p.Description,
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

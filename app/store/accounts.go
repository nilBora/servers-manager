package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CreateAccount creates a new account
func (s *DB) CreateAccount(ctx context.Context, a *Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	query := `INSERT INTO accounts (provider_id, name, login, api_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, a.ProviderID, a.Name, a.Login, a.ApiKey,
		a.CreatedAt, a.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: account with name %q already exists for this provider", ErrConflict, a.Name)
		}
		return fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	a.ID = id

	return nil
}

// GetAccount retrieves an account by ID
func (s *DB) GetAccount(ctx context.Context, id int64) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var a Account
	query := `SELECT id, provider_id, name, login, api_key, created_at, updated_at
		FROM accounts WHERE id = ?`
	if err := s.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &a, nil
}

// GetAccountWithProvider retrieves an account with provider info by ID
func (s *DB) GetAccountWithProvider(ctx context.Context, id int64) (*AccountWithProvider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var a AccountWithProvider
	query := `SELECT a.id, a.provider_id, a.name, a.login, a.api_key,
		a.created_at, a.updated_at,
		p.name as provider_name,
		(SELECT COUNT(*) FROM servers WHERE account_id = a.id) as server_count
		FROM accounts a
		JOIN providers p ON a.provider_id = p.id
		WHERE a.id = ?`
	if err := s.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &a, nil
}

// ListAccounts lists all accounts
func (s *DB) ListAccounts(ctx context.Context) ([]Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var accounts []Account
	query := `SELECT id, provider_id, name, login, api_key, created_at, updated_at
		FROM accounts ORDER BY name`
	if err := s.db.SelectContext(ctx, &accounts, query); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// ListAccountsWithProviders lists all accounts with provider info
func (s *DB) ListAccountsWithProviders(ctx context.Context) ([]AccountWithProvider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var accounts []AccountWithProvider
	query := `SELECT a.id, a.provider_id, a.name, a.login, a.api_key,
		a.created_at, a.updated_at,
		p.name as provider_name,
		(SELECT COUNT(*) FROM servers WHERE account_id = a.id) as server_count
		FROM accounts a
		JOIN providers p ON a.provider_id = p.id
		ORDER BY p.name, a.name`
	if err := s.db.SelectContext(ctx, &accounts, query); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// ListAccountsByProvider lists accounts by provider ID
func (s *DB) ListAccountsByProvider(ctx context.Context, providerID int64) ([]Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var accounts []Account
	query := `SELECT id, provider_id, name, login, api_key, created_at, updated_at
		FROM accounts WHERE provider_id = ? ORDER BY name`
	if err := s.db.SelectContext(ctx, &accounts, query, providerID); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// UpdateAccount updates an existing account
func (s *DB) UpdateAccount(ctx context.Context, a *Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	a.UpdatedAt = time.Now().UTC()

	query := `UPDATE accounts SET provider_id = ?, name = ?, login = ?, api_key = ?,
		updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, a.ProviderID, a.Name, a.Login, a.ApiKey,
		a.UpdatedAt, a.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: account with name %q already exists for this provider", ErrConflict, a.Name)
		}
		return fmt.Errorf("failed to update account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteAccount deletes an account by ID
func (s *DB) DeleteAccount(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM accounts WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("not found")

// ErrConflict is returned when a unique constraint is violated
var ErrConflict = errors.New("conflict")

// CreateProvider creates a new provider
func (s *DB) CreateProvider(ctx context.Context, p *Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	query := `INSERT INTO providers (name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Description, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: provider with name %q already exists", ErrConflict, p.Name)
		}
		return fmt.Errorf("failed to create provider: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	p.ID = id

	return nil
}

// GetProvider retrieves a provider by ID
func (s *DB) GetProvider(ctx context.Context, id int64) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var p Provider
	query := `SELECT id, name, description, created_at, updated_at FROM providers WHERE id = ?`
	if err := s.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &p, nil
}

// GetProviderByName retrieves a provider by name
func (s *DB) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var p Provider
	query := `SELECT id, name, description, created_at, updated_at FROM providers WHERE name = ?`
	if err := s.db.GetContext(ctx, &p, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &p, nil
}

// ListProviders lists all providers
func (s *DB) ListProviders(ctx context.Context) ([]Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var providers []Provider
	query := `SELECT id, name, description, created_at, updated_at FROM providers ORDER BY name`
	if err := s.db.SelectContext(ctx, &providers, query); err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	return providers, nil
}

// UpdateProvider updates an existing provider
func (s *DB) UpdateProvider(ctx context.Context, p *Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.UpdatedAt = time.Now().UTC()

	query := `UPDATE providers SET name = ?, description = ?, updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, p.Name, p.Description, p.UpdatedAt, p.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: provider with name %q already exists", ErrConflict, p.Name)
		}
		return fmt.Errorf("failed to update provider: %w", err)
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

// DeleteProvider deletes a provider by ID
func (s *DB) DeleteProvider(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM providers WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
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

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

	query := `INSERT INTO providers (name, type, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Type.String(), p.Description, p.CreatedAt, p.UpdatedAt)
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

	var p providerRow
	query := `SELECT id, name, type, description, created_at, updated_at FROM providers WHERE id = ?`
	if err := s.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return p.toProvider()
}

// GetProviderByName retrieves a provider by name
func (s *DB) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var p providerRow
	query := `SELECT id, name, type, description, created_at, updated_at FROM providers WHERE name = ?`
	if err := s.db.GetContext(ctx, &p, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return p.toProvider()
}

// ListProviders lists all providers
func (s *DB) ListProviders(ctx context.Context) ([]Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []providerRow
	query := `SELECT id, name, type, description, created_at, updated_at FROM providers ORDER BY name`
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	providers := make([]Provider, 0, len(rows))
	for _, r := range rows {
		p, err := r.toProvider()
		if err != nil {
			return nil, err
		}
		providers = append(providers, *p)
	}

	return providers, nil
}

// UpdateProvider updates an existing provider
func (s *DB) UpdateProvider(ctx context.Context, p *Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.UpdatedAt = time.Now().UTC()

	query := `UPDATE providers SET name = ?, type = ?, description = ?, updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, p.Name, p.Type.String(), p.Description, p.UpdatedAt, p.ID)
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

// providerRow is used for scanning database rows
type providerRow struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *providerRow) toProvider() (*Provider, error) {
	pt, err := parseProviderType(r.Type)
	if err != nil {
		return nil, err
	}
	return &Provider{
		ID:          r.ID,
		Name:        r.Name,
		Type:        pt,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}

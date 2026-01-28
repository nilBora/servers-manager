package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CreateUser creates a new user
func (s *DB) CreateUser(ctx context.Context, u *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := `INSERT INTO users (username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, u.Username, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: user with username %q already exists", ErrConflict, u.Username)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	u.ID = id

	return nil
}

// GetUserByUsername retrieves a user by username
func (s *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var u User
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?`
	if err := s.db.GetContext(ctx, &u, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

// GetUserByID retrieves a user by ID
func (s *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var u User
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = ?`
	if err := s.db.GetContext(ctx, &u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

// UpdateUserPassword updates a user's password
func (s *DB) UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, passwordHash, now, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
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

// CountUsers returns the number of users
func (s *DB) CountUsers(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	if err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users"); err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CreateSession creates a new session
func (s *DB) CreateSession(ctx context.Context, sess *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess.CreatedAt = time.Now().UTC()

	query := `INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, sess.ID, sess.UserID, sess.ExpiresAt, sess.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSession retrieves a session by ID
func (s *DB) GetSession(ctx context.Context, id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sess Session
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ? AND expires_at > ?`
	if err := s.db.GetContext(ctx, &sess, query, id, time.Now().UTC()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &sess, nil
}

// DeleteSession deletes a session
func (s *DB) DeleteSession(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteExpiredSessions removes all expired sessions
func (s *DB) DeleteExpiredSessions(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at <= ?", time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return nil
}

// DeleteUserSessions removes all sessions for a user
func (s *DB) DeleteUserSessions(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

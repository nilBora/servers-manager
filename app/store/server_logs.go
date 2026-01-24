package store

import (
	"context"
	"fmt"
	"time"

	"github.com/nilBora/servers-manager/app/enum"
)

// CreateLog creates a new server log entry
func (s *DB) CreateLog(ctx context.Context, l *ServerLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	l.CreatedAt = time.Now().UTC()

	query := `INSERT INTO server_logs (server_id, action, description, created_at)
		VALUES (?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, l.ServerID, l.Action.String(), l.Description, l.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	l.ID = id

	return nil
}

// ListLogs lists all logs with server info, limited by count
func (s *DB) ListLogs(ctx context.Context, limit int) ([]ServerLogWithServer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverLogWithServerRow
	query := `SELECT l.id, l.server_id, l.action, l.description, l.created_at,
		s.name as server_name, s.ip as server_ip
		FROM server_logs l
		JOIN servers s ON l.server_id = s.id
		ORDER BY l.created_at DESC
		LIMIT ?`
	if err := s.db.SelectContext(ctx, &rows, query, limit); err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	logs := make([]ServerLogWithServer, 0, len(rows))
	for _, r := range rows {
		l, err := r.toServerLogWithServer()
		if err != nil {
			return nil, err
		}
		logs = append(logs, *l)
	}

	return logs, nil
}

// ListLogsByServer lists logs for a specific server
func (s *DB) ListLogsByServer(ctx context.Context, serverID int64, limit int) ([]ServerLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverLogRow
	query := `SELECT id, server_id, action, description, created_at
		FROM server_logs WHERE server_id = ?
		ORDER BY created_at DESC
		LIMIT ?`
	if err := s.db.SelectContext(ctx, &rows, query, serverID, limit); err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	logs := make([]ServerLog, 0, len(rows))
	for _, r := range rows {
		l, err := r.toServerLog()
		if err != nil {
			return nil, err
		}
		logs = append(logs, *l)
	}

	return logs, nil
}

// ListLogsByAction lists logs filtered by action
func (s *DB) ListLogsByAction(ctx context.Context, action enum.LogAction, limit int) ([]ServerLogWithServer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverLogWithServerRow
	query := `SELECT l.id, l.server_id, l.action, l.description, l.created_at,
		s.name as server_name, s.ip as server_ip
		FROM server_logs l
		JOIN servers s ON l.server_id = s.id
		WHERE l.action = ?
		ORDER BY l.created_at DESC
		LIMIT ?`
	if err := s.db.SelectContext(ctx, &rows, query, action.String(), limit); err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	logs := make([]ServerLogWithServer, 0, len(rows))
	for _, r := range rows {
		l, err := r.toServerLogWithServer()
		if err != nil {
			return nil, err
		}
		logs = append(logs, *l)
	}

	return logs, nil
}

// serverLogRow is used for scanning database rows
type serverLogRow struct {
	ID          int64     `db:"id"`
	ServerID    int64     `db:"server_id"`
	Action      string    `db:"action"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r *serverLogRow) toServerLog() (*ServerLog, error) {
	action, err := enum.ParseLogAction(r.Action)
	if err != nil {
		return nil, err
	}
	return &ServerLog{
		ID:          r.ID,
		ServerID:    r.ServerID,
		Action:      action,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
	}, nil
}

type serverLogWithServerRow struct {
	serverLogRow
	ServerName string `db:"server_name"`
	ServerIP   string `db:"server_ip"`
}

func (r *serverLogWithServerRow) toServerLogWithServer() (*ServerLogWithServer, error) {
	l, err := r.serverLogRow.toServerLog()
	if err != nil {
		return nil, err
	}
	return &ServerLogWithServer{
		ServerLog:  *l,
		ServerName: r.ServerName,
		ServerIP:   r.ServerIP,
	}, nil
}

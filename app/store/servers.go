package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nilBora/servers-manager/app/enum"
)

// CreateServer creates a new server
func (s *DB) CreateServer(ctx context.Context, srv *Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	srv.CreatedAt = now
	srv.UpdatedAt = now

	query := `INSERT INTO servers (account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query, srv.AccountID, srv.Name, srv.IP, srv.Location,
		srv.Description, srv.Responsible, srv.ApproximateCost, srv.Backups, srv.Status.String(),
		srv.CreatedAt, srv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	srv.ID = id

	return nil
}

// GetServer retrieves a server by ID
func (s *DB) GetServer(ctx context.Context, id int64) (*Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var r serverRow
	query := `SELECT id, account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at
		FROM servers WHERE id = ?`
	if err := s.db.GetContext(ctx, &r, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	return r.toServer()
}

// GetServerWithAccount retrieves a server with account info by ID
func (s *DB) GetServerWithAccount(ctx context.Context, id int64) (*ServerWithAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var r serverWithAccountRow
	query := `SELECT s.id, s.account_id, s.name, s.ip, s.location, s.description, s.responsible,
		s.approximate_cost, s.backups, s.status, s.created_at, s.updated_at,
		a.name as account_name, a.group_name as account_group_name, a.provider_id,
		p.name as provider_name
		FROM servers s
		JOIN accounts a ON s.account_id = a.id
		JOIN providers p ON a.provider_id = p.id
		WHERE s.id = ?`
	if err := s.db.GetContext(ctx, &r, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	return r.toServerWithAccount()
}

// ListServers lists all servers
func (s *DB) ListServers(ctx context.Context) ([]Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverRow
	query := `SELECT id, account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at
		FROM servers ORDER BY name`
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}

	servers := make([]Server, 0, len(rows))
	for _, r := range rows {
		srv, err := r.toServer()
		if err != nil {
			return nil, err
		}
		servers = append(servers, *srv)
	}

	return servers, nil
}

// ListServersWithAccounts lists all servers with account info
func (s *DB) ListServersWithAccounts(ctx context.Context) ([]ServerWithAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverWithAccountRow
	query := `SELECT s.id, s.account_id, s.name, s.ip, s.location, s.description, s.responsible,
		s.approximate_cost, s.backups, s.status, s.created_at, s.updated_at,
		a.name as account_name, a.group_name as account_group_name, a.provider_id,
		p.name as provider_name
		FROM servers s
		JOIN accounts a ON s.account_id = a.id
		JOIN providers p ON a.provider_id = p.id
		ORDER BY p.name, a.group_name, a.name, s.name`
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}

	servers := make([]ServerWithAccount, 0, len(rows))
	for _, r := range rows {
		srv, err := r.toServerWithAccount()
		if err != nil {
			return nil, err
		}
		servers = append(servers, *srv)
	}

	return servers, nil
}

// ListServersByAccount lists servers by account ID
func (s *DB) ListServersByAccount(ctx context.Context, accountID int64) ([]Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverRow
	query := `SELECT id, account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at
		FROM servers WHERE account_id = ? ORDER BY name`
	if err := s.db.SelectContext(ctx, &rows, query, accountID); err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}

	servers := make([]Server, 0, len(rows))
	for _, r := range rows {
		srv, err := r.toServer()
		if err != nil {
			return nil, err
		}
		servers = append(servers, *srv)
	}

	return servers, nil
}

// ListServersByStatus lists servers by status
func (s *DB) ListServersByStatus(ctx context.Context, status enum.ServerStatus) ([]ServerWithAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows []serverWithAccountRow
	query := `SELECT s.id, s.account_id, s.name, s.ip, s.location, s.description, s.responsible,
		s.approximate_cost, s.backups, s.status, s.created_at, s.updated_at,
		a.name as account_name, a.group_name as account_group_name, a.provider_id,
		p.name as provider_name
		FROM servers s
		JOIN accounts a ON s.account_id = a.id
		JOIN providers p ON a.provider_id = p.id
		WHERE s.status = ?
		ORDER BY p.name, a.group_name, a.name, s.name`
	if err := s.db.SelectContext(ctx, &rows, query, status.String()); err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}

	servers := make([]ServerWithAccount, 0, len(rows))
	for _, r := range rows {
		srv, err := r.toServerWithAccount()
		if err != nil {
			return nil, err
		}
		servers = append(servers, *srv)
	}

	return servers, nil
}

// UpdateServer updates an existing server
func (s *DB) UpdateServer(ctx context.Context, srv *Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	srv.UpdatedAt = time.Now().UTC()

	query := `UPDATE servers SET account_id = ?, name = ?, ip = ?, location = ?, description = ?,
		responsible = ?, approximate_cost = ?, backups = ?, status = ?, updated_at = ?
		WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, srv.AccountID, srv.Name, srv.IP, srv.Location,
		srv.Description, srv.Responsible, srv.ApproximateCost, srv.Backups, srv.Status.String(),
		srv.UpdatedAt, srv.ID)
	if err != nil {
		return fmt.Errorf("failed to update server: %w", err)
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

// UpdateServerStatus updates only the status of a server
func (s *DB) UpdateServerStatus(ctx context.Context, id int64, status enum.ServerStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()

	query := `UPDATE servers SET status = ?, updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, status.String(), now, id)
	if err != nil {
		return fmt.Errorf("failed to update server status: %w", err)
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

// FindServerByNameAndAccount finds a server by name and account ID
func (s *DB) FindServerByNameAndAccount(ctx context.Context, name string, accountID int64) (*Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var r serverRow
	query := `SELECT id, account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at
		FROM servers WHERE name = ? AND account_id = ?`
	if err := s.db.GetContext(ctx, &r, query, name, accountID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find server: %w", err)
	}

	return r.toServer()
}

// FindServerByIPAndAccount finds a server by IP and account ID
func (s *DB) FindServerByIPAndAccount(ctx context.Context, ip string, accountID int64) (*Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var r serverRow
	query := `SELECT id, account_id, name, ip, location, description, responsible,
		approximate_cost, backups, status, created_at, updated_at
		FROM servers WHERE ip = ? AND account_id = ?`
	if err := s.db.GetContext(ctx, &r, query, ip, accountID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find server by IP: %w", err)
	}

	return r.toServer()
}

// DeleteServer deletes a server by ID
func (s *DB) DeleteServer(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM servers WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
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

// GetDashboardStats returns dashboard statistics
func (s *DB) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats DashboardStats
	query := `SELECT
		COUNT(*) as total_servers,
		COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) as active_servers,
		COALESCE(SUM(CASE WHEN status = 'paused' THEN 1 ELSE 0 END), 0) as paused_servers,
		COALESCE(SUM(CASE WHEN status != 'deleted' THEN approximate_cost ELSE 0 END), 0) as total_cost
		FROM servers`
	if err := s.db.GetContext(ctx, &stats, query); err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return &stats, nil
}

// GetServersGroupedByAccount returns servers grouped by account for dashboard
func (s *DB) GetServersGroupedByAccount(ctx context.Context, status *enum.ServerStatus) ([]AccountGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// build query with optional status filter
	query := `SELECT s.id, s.account_id, s.name, s.ip, s.location, s.description, s.responsible,
		s.approximate_cost, s.backups, s.status, s.created_at, s.updated_at,
		a.name as account_name, a.group_name as account_group_name, a.provider_id,
		p.name as provider_name
		FROM servers s
		JOIN accounts a ON s.account_id = a.id
		JOIN providers p ON a.provider_id = p.id`

	var args []interface{}
	if status != nil {
		query += ` WHERE s.status = ?`
		args = append(args, status.String())
	}
	query += ` ORDER BY p.name, a.group_name, a.name, s.name`

	var rows []serverWithAccountRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get servers grouped: %w", err)
	}

	// group by account
	groupMap := make(map[int64]*AccountGroup)
	var groups []AccountGroup

	for _, r := range rows {
		srv, err := r.toServerWithAccount()
		if err != nil {
			return nil, err
		}

		group, exists := groupMap[srv.AccountID]
		if !exists {
			group = &AccountGroup{
				AccountID:        srv.AccountID,
				AccountName:      srv.AccountName,
				AccountGroupName: srv.AccountGroupName,
				ProviderID:       srv.ProviderID,
				ProviderName:     srv.ProviderName,
				Servers:          make([]ServerWithAccount, 0),
			}
			groupMap[srv.AccountID] = group
			groups = append(groups, *group)
		}

		// find and update group in slice
		for i := range groups {
			if groups[i].AccountID == srv.AccountID {
				groups[i].Servers = append(groups[i].Servers, *srv)
				if srv.Status != enum.ServerStatusDeleted {
					groups[i].TotalCost += srv.ApproximateCost
				}
				break
			}
		}
	}

	return groups, nil
}

// GetServersGroupedHierarchically returns servers in hierarchical structure:
// Provider+GroupName -> Accounts (Projects) -> Servers
func (s *DB) GetServersGroupedHierarchically(ctx context.Context, status *enum.ServerStatus) ([]ProviderAccountGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// build query with optional status filter
	query := `SELECT s.id, s.account_id, s.name, s.ip, s.location, s.description, s.responsible,
		s.approximate_cost, s.backups, s.status, s.created_at, s.updated_at,
		a.name as account_name, a.group_name as account_group_name, a.provider_id,
		p.name as provider_name
		FROM servers s
		JOIN accounts a ON s.account_id = a.id
		JOIN providers p ON a.provider_id = p.id`

	var args []interface{}
	if status != nil {
		query += ` WHERE s.status = ?`
		args = append(args, status.String())
	}
	query += ` ORDER BY p.name, a.group_name, a.name, s.name`

	var rows []serverWithAccountRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get servers grouped: %w", err)
	}

	// Build hierarchical structure: ProviderAccountGroup -> AccountGroup -> Servers
	providerGroupMap := make(map[string]*ProviderAccountGroup) // key: "provider_id:group_name"
	accountGroupMap := make(map[int64]*AccountGroup)           // key: account_id
	var providerGroups []ProviderAccountGroup

	for _, r := range rows {
		srv, err := r.toServerWithAccount()
		if err != nil {
			return nil, err
		}

		// Create group key from provider_id and group_name
		groupKey := fmt.Sprintf("%d:%s", srv.ProviderID, srv.AccountGroupName)

		// Find or create ProviderAccountGroup
		providerGroup, exists := providerGroupMap[groupKey]
		if !exists {
			providerGroup = &ProviderAccountGroup{
				ProviderID:   srv.ProviderID,
				ProviderName: srv.ProviderName,
				GroupName:    srv.AccountGroupName,
				GroupKey:     groupKey,
				Accounts:     make([]AccountGroup, 0),
			}
			providerGroupMap[groupKey] = providerGroup
			providerGroups = append(providerGroups, *providerGroup)
		}

		// Find or create AccountGroup within ProviderAccountGroup
		accountGroup, exists := accountGroupMap[srv.AccountID]
		if !exists {
			accountGroup = &AccountGroup{
				AccountID:        srv.AccountID,
				AccountName:      srv.AccountName,
				AccountGroupName: srv.AccountGroupName,
				ProviderID:       srv.ProviderID,
				ProviderName:     srv.ProviderName,
				Servers:          make([]ServerWithAccount, 0),
			}
			accountGroupMap[srv.AccountID] = accountGroup

			// Add to the appropriate ProviderAccountGroup
			for i := range providerGroups {
				if providerGroups[i].GroupKey == groupKey {
					providerGroups[i].Accounts = append(providerGroups[i].Accounts, *accountGroup)
					break
				}
			}
		}

		// Add server to account group (find it in providerGroups slice)
		for i := range providerGroups {
			if providerGroups[i].GroupKey == groupKey {
				for j := range providerGroups[i].Accounts {
					if providerGroups[i].Accounts[j].AccountID == srv.AccountID {
						providerGroups[i].Accounts[j].Servers = append(providerGroups[i].Accounts[j].Servers, *srv)
						if srv.Status != enum.ServerStatusDeleted {
							providerGroups[i].Accounts[j].TotalCost += srv.ApproximateCost
							providerGroups[i].TotalCost += srv.ApproximateCost
						}
						providerGroups[i].ServerCount++
						break
					}
				}
				break
			}
		}
	}

	return providerGroups, nil
}

// serverRow is used for scanning database rows
type serverRow struct {
	ID              int64     `db:"id"`
	AccountID       int64     `db:"account_id"`
	Name            string    `db:"name"`
	IP              string    `db:"ip"`
	Location        string    `db:"location"`
	Description     string    `db:"description"`
	Responsible     string    `db:"responsible"`
	ApproximateCost float64   `db:"approximate_cost"`
	Backups         bool      `db:"backups"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (r *serverRow) toServer() (*Server, error) {
	st, err := enum.ParseServerStatus(r.Status)
	if err != nil {
		return nil, err
	}
	return &Server{
		ID:              r.ID,
		AccountID:       r.AccountID,
		Name:            r.Name,
		IP:              r.IP,
		Location:        r.Location,
		Description:     r.Description,
		Responsible:     r.Responsible,
		ApproximateCost: r.ApproximateCost,
		Backups:         r.Backups,
		Status:          st,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}, nil
}

type serverWithAccountRow struct {
	serverRow
	AccountName      string `db:"account_name"`
	AccountGroupName string `db:"account_group_name"`
	ProviderID       int64  `db:"provider_id"`
	ProviderName     string `db:"provider_name"`
}

func (r *serverWithAccountRow) toServerWithAccount() (*ServerWithAccount, error) {
	srv, err := r.serverRow.toServer()
	if err != nil {
		return nil, err
	}
	return &ServerWithAccount{
		Server:           *srv,
		AccountName:      r.AccountName,
		AccountGroupName: r.AccountGroupName,
		ProviderID:       r.ProviderID,
		ProviderName:     r.ProviderName,
	}, nil
}

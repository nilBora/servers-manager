package store

import (
	"context"

	"github.com/nilBora/servers-manager/app/enum"
)

// ProviderStore defines operations for providers
type ProviderStore interface {
	CreateProvider(ctx context.Context, p *Provider) error
	GetProvider(ctx context.Context, id int64) (*Provider, error)
	GetProviderByName(ctx context.Context, name string) (*Provider, error)
	ListProviders(ctx context.Context) ([]Provider, error)
	UpdateProvider(ctx context.Context, p *Provider) error
	DeleteProvider(ctx context.Context, id int64) error
}

// AccountStore defines operations for accounts
type AccountStore interface {
	CreateAccount(ctx context.Context, a *Account) error
	GetAccount(ctx context.Context, id int64) (*Account, error)
	GetAccountWithProvider(ctx context.Context, id int64) (*AccountWithProvider, error)
	ListAccounts(ctx context.Context) ([]Account, error)
	ListAccountsWithProviders(ctx context.Context) ([]AccountWithProvider, error)
	ListAccountsByProvider(ctx context.Context, providerID int64) ([]Account, error)
	UpdateAccount(ctx context.Context, a *Account) error
	DeleteAccount(ctx context.Context, id int64) error
}

// ServerStore defines operations for servers
type ServerStore interface {
	CreateServer(ctx context.Context, s *Server) error
	GetServer(ctx context.Context, id int64) (*Server, error)
	GetServerWithAccount(ctx context.Context, id int64) (*ServerWithAccount, error)
	FindServerByNameAndAccount(ctx context.Context, name string, accountID int64) (*Server, error)
	ListServers(ctx context.Context) ([]Server, error)
	ListServersWithAccounts(ctx context.Context) ([]ServerWithAccount, error)
	ListServersByAccount(ctx context.Context, accountID int64) ([]Server, error)
	ListServersByStatus(ctx context.Context, status enum.ServerStatus) ([]ServerWithAccount, error)
	UpdateServer(ctx context.Context, s *Server) error
	UpdateServerStatus(ctx context.Context, id int64, status enum.ServerStatus) error
	DeleteServer(ctx context.Context, id int64) error
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetServersGroupedByAccount(ctx context.Context, status *enum.ServerStatus) ([]AccountGroup, error)
	GetServersGroupedHierarchically(ctx context.Context, status *enum.ServerStatus) ([]ProviderAccountGroup, error)
}

// ServerLogStore defines operations for server logs
type ServerLogStore interface {
	CreateLog(ctx context.Context, l *ServerLog) error
	ListLogs(ctx context.Context, limit int) ([]ServerLogWithServer, error)
	ListLogsByServer(ctx context.Context, serverID int64, limit int) ([]ServerLog, error)
	ListLogsByAction(ctx context.Context, action enum.LogAction, limit int) ([]ServerLogWithServer, error)
}

// UserStore defines operations for users
type UserStore interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error
	CountUsers(ctx context.Context) (int, error)
}

// SessionStore defines operations for sessions
type SessionStore interface {
	CreateSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, id string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteExpiredSessions(ctx context.Context) error
	DeleteUserSessions(ctx context.Context, userID int64) error
}

// Store combines all store interfaces
type Store interface {
	ProviderStore
	AccountStore
	ServerStore
	ServerLogStore
	UserStore
	SessionStore
	Close() error
}

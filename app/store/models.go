package store

import (
	"time"

	"github.com/nilBora/servers-manager/app/enum"
)

// Provider represents a cloud infrastructure provider
type Provider struct {
	ID          int64              `db:"id"`
	Name        string             `db:"name"`
	Type        enum.ProviderType  `db:"type"`
	Description string             `db:"description"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
}

// Account represents an account at a provider
type Account struct {
	ID          int64            `db:"id"`
	ProviderID  int64            `db:"provider_id"`
	Name        string           `db:"name"`
	Login       string           `db:"login"`
	ApiKey      string           `db:"api_key"`
	AccountType enum.AccountType `db:"account_type"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

// AccountWithProvider extends Account with provider info for display
type AccountWithProvider struct {
	Account
	ProviderName string            `db:"provider_name"`
	ProviderType enum.ProviderType `db:"provider_type"`
	ServerCount  int               `db:"server_count"`
}

// Server represents a server instance
type Server struct {
	ID              int64             `db:"id"`
	AccountID       int64             `db:"account_id"`
	Name            string            `db:"name"`
	IP              string            `db:"ip"`
	Description     string            `db:"description"`
	Responsible     string            `db:"responsible"`
	ApproximateCost float64           `db:"approximate_cost"`
	Status          enum.ServerStatus `db:"status"`
	ServerType      enum.ServerType   `db:"server_type"`
	CreatedAt       time.Time         `db:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"`
}

// ServerWithAccount extends Server with account and provider info for display
type ServerWithAccount struct {
	Server
	AccountName  string            `db:"account_name"`
	ProviderID   int64             `db:"provider_id"`
	ProviderName string            `db:"provider_name"`
	ProviderType enum.ProviderType `db:"provider_type"`
}

// ServerLog represents a server action log entry
type ServerLog struct {
	ID          int64          `db:"id"`
	ServerID    int64          `db:"server_id"`
	Action      enum.LogAction `db:"action"`
	Description string         `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
}

// ServerLogWithServer extends ServerLog with server info for display
type ServerLogWithServer struct {
	ServerLog
	ServerName string `db:"server_name"`
	ServerIP   string `db:"server_ip"`
}

// DashboardStats holds dashboard statistics
type DashboardStats struct {
	TotalServers  int     `db:"total_servers"`
	ActiveServers int     `db:"active_servers"`
	PausedServers int     `db:"paused_servers"`
	TotalCost     float64 `db:"total_cost"`
}

// AccountGroup groups servers by account for dashboard display
type AccountGroup struct {
	AccountID    int64
	AccountName  string
	AccountType  enum.AccountType
	ProviderID   int64
	ProviderName string
	ProviderType enum.ProviderType
	Servers      []ServerWithAccount
	TotalCost    float64
}

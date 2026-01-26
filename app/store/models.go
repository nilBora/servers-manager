package store

import (
	"time"

	"github.com/nilBora/servers-manager/app/enum"
)

// Provider represents a cloud infrastructure provider
type Provider struct {
	ID          int64     `db:"id"`
	Ident       string    `db:"ident"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Account represents an account at a provider
type Account struct {
	ID         int64     `db:"id"`
	ProviderID int64     `db:"provider_id"`
	GroupName  string    `db:"group_name"`
	Name       string    `db:"name"`
	Login      string    `db:"login"`
	ApiKey     string    `db:"api_key"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// AccountWithProvider extends Account with provider info for display
type AccountWithProvider struct {
	Account
	ProviderIdent string `db:"provider_ident"`
	ProviderName  string `db:"provider_name"`
	ServerCount   int    `db:"server_count"`
}

// AccountGroupSummary represents a group of accounts for dashboard display
type AccountGroupSummary struct {
	ProviderID   int64
	ProviderName string
	GroupName    string
	Accounts     []AccountWithServers
	TotalCost    float64
	ServerCount  int
}

// AccountWithServers extends account with its servers
type AccountWithServers struct {
	Account
	ProviderName string
	Servers      []ServerWithAccount
	TotalCost    float64
}

// Server represents a server instance
type Server struct {
	ID              int64             `db:"id"`
	AccountID       int64             `db:"account_id"`
	Name            string            `db:"name"`
	IP              string            `db:"ip"`
	Location        string            `db:"location"`
	Description     string            `db:"description"`
	Responsible     string            `db:"responsible"`
	ApproximateCost float64           `db:"approximate_cost"`
	Status          enum.ServerStatus `db:"status"`
	CreatedAt       time.Time         `db:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"`
}

// ServerWithAccount extends Server with account and provider info for display
type ServerWithAccount struct {
	Server
	AccountName      string `db:"account_name"`
	AccountGroupName string `db:"account_group_name"`
	ProviderID       int64  `db:"provider_id"`
	ProviderName     string `db:"provider_name"`
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
	AccountID        int64
	AccountName      string
	AccountGroupName string
	ProviderID       int64
	ProviderName     string
	Servers          []ServerWithAccount
	TotalCost        float64
}

// ProviderAccountGroup groups accounts by provider and group_name for hierarchical display
// Structure: Provider + GroupName -> Accounts (Projects) -> Servers
type ProviderAccountGroup struct {
	ProviderID   int64
	ProviderName string
	GroupName    string
	GroupKey     string // unique key: "provider_id:group_name"
	Accounts     []AccountGroup
	TotalCost    float64
	ServerCount  int
}

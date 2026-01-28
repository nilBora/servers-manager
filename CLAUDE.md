# Servers Manager

Go/HTMX application for managing servers across multiple cloud providers.

## Tech Stack

- **Backend**: Go with go-chi router
- **Frontend**: HTMX + vanilla JavaScript
- **Database**: SQLite (with sqlx)
- **Templates**: Go html/template (embedded)
- **Auth**: Session-based with bcrypt password hashing

## Project Structure

```
app/
├── main.go                    # Entry point, CLI flags
├── enum/                      # Type-safe enums
│   └── enum.go
├── store/                     # Database layer
│   ├── store.go              # Interfaces
│   ├── db.go                 # SQLite setup, schema, migrations
│   ├── models.go             # Data models
│   ├── providers.go          # Provider CRUD
│   ├── accounts.go           # Account CRUD
│   ├── servers.go            # Server CRUD + dashboard queries
│   ├── server_logs.go        # Log operations
│   └── users.go              # User & session operations
├── hetzner/                   # Hetzner API clients
│   ├── client.go             # Cloud API
│   └── robot.go              # Robot API (dedicated servers)
└── server/
    ├── server.go             # HTTP server setup
    └── web/
        ├── handler.go        # Template loading, routes
        ├── auth.go           # Authentication middleware & handlers
        ├── pages.go          # Page handlers
        ├── providers.go      # Provider handlers
        ├── accounts.go       # Account handlers
        ├── servers.go        # Server handlers
        ├── dashboard.go      # Dashboard handlers
        ├── sync.go           # Hetzner sync handlers
        ├── static/           # CSS, JS (embedded)
        └── templates/        # HTML templates (embedded)
```

## Key Concepts

### Providers
Cloud providers with unique `ident` for API integrations:
- `hetzner_cloud` - Hetzner Cloud (API token)
- `hetzner_robot` - Hetzner Robot (username:password)
- `aws` - Amazon Web Services
- `scaleway` - Scaleway
- `vsys_host` - Vsys Host

### Accounts
Credentials for provider access. Fields:
- `provider_id` - Link to provider
- `group_name` - Visual grouping (e.g., "Main Account")
- `name` - Account/project name
- `api_key` - API credentials

### Servers
Server instances. Fields:
- `account_id` - Link to account
- `name`, `ip`, `location`, `description`
- `responsible` - Person/team responsible
- `approximate_cost` - Monthly cost
- `status` - active/paused/deleted

### Server Logs
Action history for servers:
- Actions: added, updated, paused, deleted, synced

## Running

```bash
go run ./app --db=servers.db --port=8080
```

First run opens `/setup` to create admin account.

## API Sync

### Hetzner Cloud
- Provider ident: `hetzner_cloud`
- API Key format: Bearer token
- Syncs: name, IP, location, description, monthly price

### Hetzner Robot
- Provider ident: `hetzner_robot`
- API Key format: `username:password`
- Syncs: name, IP, datacenter, product info

## Database

SQLite with WAL mode. Migrations run automatically on startup.

Tables: `providers`, `accounts`, `servers`, `server_logs`, `users`, `sessions`

## Authentication

- Session-based auth with HttpOnly cookies
- Sessions stored in database (7-day expiry)
- Passwords hashed with bcrypt
- All routes protected except `/login`, `/setup`, `/logout`

## Frontend

- HTMX for dynamic updates without page reloads
- Dashboard blocks are drag-and-drop sortable (order saved in localStorage)
- Light/dark theme toggle (saved in cookie)

# Servers Manager

A lightweight web application for managing servers across multiple cloud providers. Built with Go and HTMX — no JavaScript framework required.

## Features

- Manage servers from multiple cloud providers in one place
- Automatic daily sync with provider APIs
- Manual sync on demand
- Server status tracking (active / paused / deleted)
- Action logs for every server change
- Light / dark theme toggle
- Session-based authentication

## Supported Providers

| Provider | Ident | Auth |
|---|---|---|
| Hetzner Cloud | `hetzner_cloud` | API token |
| Hetzner Robot | `hetzner_robot` | `username:password` |
| AWS | `aws` | — |
| Scaleway | `scaleway` | — |
| Vsys Host | `vsys_host` | — |

## Tech Stack

- **Backend**: Go, [chi](https://github.com/go-chi/chi) router
- **Frontend**: [HTMX](https://htmx.org), vanilla JS
- **Database**: SQLite (via [sqlx](https://github.com/jmoiron/sqlx) + [modernc sqlite](https://gitlab.com/cznic/sqlite))
- **Templates**: Go `html/template` (embedded)
- **Auth**: Session-based, bcrypt password hashing

## Getting Started

### Requirements

- Go 1.24+

### Run

```bash
go run ./app --db=servers.db --address=:8080
```

On first run, open `http://localhost:8080/setup` to create an admin account.

### Flags

| Flag | Env | Default | Description |
|---|---|---|---|
| `--db` | `DB` | `servers.db` | SQLite database file path |
| `--address` | `ADDRESS` | `:8080` | HTTP listen address |
| `--debug` | `DEBUG` | `false` | Enable debug logging |

## Project Structure

```
app/
├── main.go                 # Entry point, CLI flags
├── enum/                   # Type-safe enums (theme, status, log action)
├── store/                  # Database layer (SQLite, migrations, CRUD)
├── hetzner/                # Hetzner Cloud and Robot API clients
└── server/
    ├── server.go           # HTTP server, graceful shutdown
    └── web/
        ├── handler.go      # Routes, template loading
        ├── auth.go         # Auth middleware and handlers
        ├── pages.go        # Page handlers
        ├── accounts.go     # Account CRUD handlers
        ├── providers.go    # Provider CRUD handlers
        ├── servers.go      # Server CRUD handlers
        ├── dashboard.go    # Dashboard and theme handlers
        ├── sync.go         # Provider sync logic
        ├── static/         # Embedded CSS and JS
        └── templates/      # Embedded HTML templates
```

## Database

SQLite with WAL mode. Schema migrations run automatically on startup.

Tables: `providers`, `accounts`, `servers`, `server_logs`, `users`, `sessions`

## Sync

Servers are synced automatically once per day in the background. A manual sync can be triggered from the Servers page.

Synced fields per provider:

**Hetzner Cloud** — name, IP, location, server type, monthly price, backup status
**Hetzner Robot** — name, IP, datacenter, product info

Servers no longer present in the API are automatically marked as deleted.

## License

MIT

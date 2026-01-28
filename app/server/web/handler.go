// Package web provides HTTP handlers for the web UI
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/nilBora/servers-manager/app/enum"
	"github.com/nilBora/servers-manager/app/store"
)

//go:embed static
var staticFS embed.FS

//go:embed templates
var templatesFS embed.FS

// StaticFS returns the embedded static filesystem for external use
func StaticFS() (fs.FS, error) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to get static sub-filesystem: %w", err)
	}
	return sub, nil
}

// Handler handles web UI requests
type Handler struct {
	store store.Store
	tmpl  *template.Template
}

// New creates a new web handler
func New(st store.Store) (*Handler, error) {
	tmpl, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Handler{
		store: st,
		tmpl:  tmpl,
	}, nil
}

// Register registers web UI routes on the given router
func (h *Handler) Register(r chi.Router) {
	// Public routes (no auth required)
	r.Get("/login", h.handleLogin)
	r.Post("/login", h.handleLoginPost)
	r.Get("/setup", h.handleSetup)
	r.Post("/setup", h.handleSetupPost)
	r.Get("/logout", h.handleLogout)

	// Protected routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		// pages
		r.Get("/", h.handleDashboard)
		r.Get("/providers", h.handleProviders)
		r.Get("/accounts", h.handleAccounts)
		r.Get("/servers", h.handleServers)
		r.Get("/logs", h.handleLogs)

		// provider CRUD
		r.Get("/web/providers", h.handleProviderTable)
		r.Get("/web/providers/new", h.handleProviderForm)
		r.Get("/web/providers/{id}/edit", h.handleProviderEditForm)
		r.Post("/web/providers", h.handleProviderCreate)
		r.Put("/web/providers/{id}", h.handleProviderUpdate)
		r.Delete("/web/providers/{id}", h.handleProviderDelete)

		// account CRUD
		r.Get("/web/accounts", h.handleAccountTable)
		r.Get("/web/accounts/new", h.handleAccountForm)
		r.Get("/web/accounts/{id}/edit", h.handleAccountEditForm)
		r.Post("/web/accounts", h.handleAccountCreate)
		r.Put("/web/accounts/{id}", h.handleAccountUpdate)
		r.Delete("/web/accounts/{id}", h.handleAccountDelete)

		// server CRUD
		r.Get("/web/servers", h.handleServerTable)
		r.Get("/web/servers/new", h.handleServerForm)
		r.Get("/web/servers/{id}/edit", h.handleServerEditForm)
		r.Get("/web/servers/{id}/view", h.handleServerView)
		r.Post("/web/servers", h.handleServerCreate)
		r.Put("/web/servers/{id}", h.handleServerUpdate)
		r.Put("/web/servers/{id}/status", h.handleServerStatusUpdate)
		r.Delete("/web/servers/{id}", h.handleServerDelete)

		// logs
		r.Get("/web/logs", h.handleLogTable)

		// sync
		r.Post("/web/sync/hetzner", h.handleHetznerSync)

		// dashboard
		r.Get("/web/dashboard", h.handleDashboardContent)
		r.Get("/web/dashboard/stats", h.handleDashboardStats)

		// settings
		r.Post("/web/theme", h.handleThemeToggle)
	})
}

// templateFuncs returns custom template functions
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatCost": func(cost float64) string {
			return fmt.Sprintf("$%.2f", cost)
		},
		"maskApiKey": func(key string) string {
			if len(key) <= 8 {
				return "****"
			}
			return key[:4] + "****" + key[len(key)-4:]
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"eq":  func(a, b interface{}) bool { return a == b },
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"statusClass": func(status enum.ServerStatus) string {
			switch status {
			case enum.ServerStatusActive:
				return "status-active"
			case enum.ServerStatusPaused:
				return "status-paused"
			case enum.ServerStatusDeleted:
				return "status-deleted"
			}
			return ""
		},
		"actionClass": func(action enum.LogAction) string {
			switch action {
			case enum.LogActionAdded:
				return "action-added"
			case enum.LogActionDeleted:
				return "action-deleted"
			case enum.LogActionPaused:
				return "action-paused"
			case enum.LogActionUpdated:
				return "action-updated"
			case enum.LogActionSynced:
				return "action-synced"
			}
			return ""
		},
	}
}

// parseTemplates parses all templates from embedded filesystem
func parseTemplates() (*template.Template, error) {
	tmpl := template.New("").Funcs(templateFuncs())

	// parse partials first (they are used by main templates)
	partials := []string{
		"provider-table",
		"provider-form",
		"account-table",
		"account-form",
		"server-table",
		"server-form",
		"server-card",
		"server-logs",
		"dashboard-stats",
		"dashboard-accounts",
		"status-badge",
		"nav",
	}

	for _, name := range partials {
		content, err := templatesFS.ReadFile("templates/partials/" + name + ".html")
		if err != nil {
			return nil, fmt.Errorf("read partial %s: %w", name, err)
		}
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse partial %s: %w", name, err)
		}
	}

	// parse page templates (each is self-contained)
	pages := []string{
		"dashboard.html",
		"providers.html",
		"accounts.html",
		"servers.html",
		"logs.html",
		"login.html",
		"setup.html",
	}

	for _, name := range pages {
		content, err := templatesFS.ReadFile("templates/" + name)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
	}

	return tmpl, nil
}

// templateData holds common data passed to templates
type templateData struct {
	Theme        enum.Theme
	ActivePage   string
	Error        string
	Success      string

	// dashboard data
	Stats          *store.DashboardStats
	Groups         []store.AccountGroup
	ProviderGroups []store.ProviderAccountGroup
	StatusFilter   string

	// providers data
	Providers []store.Provider
	Provider  *store.Provider

	// accounts data
	Accounts []store.AccountWithProvider
	Account  *store.AccountWithProvider

	// servers data
	Servers  []store.ServerWithAccount
	Server   *store.ServerWithAccount
	Statuses []enum.ServerStatus

	// logs data
	Logs         []store.ServerLogWithServer
	Actions      []enum.LogAction
	ActionFilter string
}

// getTheme returns the current theme from cookie
func (h *Handler) getTheme(r *http.Request) enum.Theme {
	if c, err := r.Cookie("theme"); err == nil {
		if theme, err := enum.ParseTheme(c.Value); err == nil {
			return theme
		}
	}
	return enum.ThemeSystem
}

// parseID parses an ID from the URL parameter
func parseID(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %s", idStr)
	}
	return id, nil
}

// renderError renders an error response
func (h *Handler) renderError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`<div class="error">%s</div>`, message)))
}

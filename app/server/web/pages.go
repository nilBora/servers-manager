package web

import (
	"net/http"

	"github.com/nilBora/servers-manager/app/enum"
)

// handleDashboard renders the dashboard page
func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDashboardStats(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load stats")
		return
	}

	groups, err := h.store.GetServersGroupedByAccount(r.Context(), nil)
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load servers")
		return
	}

	data := templateData{
		Theme:      h.getTheme(r),
		ActivePage: "dashboard",
		Stats:      stats,
		Groups:     groups,
		Statuses:   enum.AllServerStatuses(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProviders renders the providers page
func (h *Handler) handleProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.store.ListProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load providers")
		return
	}

	data := templateData{
		Theme:         h.getTheme(r),
		ActivePage:    "providers",
		Providers:     providers,
		ProviderTypes: enum.AllProviderTypes(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "providers.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAccounts renders the accounts page
func (h *Handler) handleAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.store.ListAccountsWithProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	providers, err := h.store.ListProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load providers")
		return
	}

	data := templateData{
		Theme:        h.getTheme(r),
		ActivePage:   "accounts",
		Accounts:     accounts,
		Providers:    providers,
		AccountTypes: enum.AllAccountTypes(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "accounts.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleServers renders the servers page
func (h *Handler) handleServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.store.ListServersWithAccounts(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load servers")
		return
	}

	accounts, err := h.store.ListAccountsWithProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	data := templateData{
		Theme:       h.getTheme(r),
		ActivePage:  "servers",
		Servers:     servers,
		Accounts:    accounts,
		ServerTypes: enum.AllServerTypes(),
		Statuses:    enum.AllServerStatuses(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "servers.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleLogs renders the logs page
func (h *Handler) handleLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.store.ListLogs(r.Context(), 100)
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load logs")
		return
	}

	data := templateData{
		Theme:      h.getTheme(r),
		ActivePage: "logs",
		Logs:       logs,
		Actions:    enum.AllLogActions(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "logs.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

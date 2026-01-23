package web

import (
	"net/http"

	"github.com/nilBora/servers-manager/app/enum"
	"github.com/nilBora/servers-manager/app/store"
)

// handleDashboardContent renders the dashboard content partial
func (h *Handler) handleDashboardContent(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	var status *enum.ServerStatus
	if statusFilter != "" {
		s, err := enum.ParseServerStatus(statusFilter)
		if err == nil {
			status = &s
		}
	}

	stats, err := h.store.GetDashboardStats(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load stats")
		return
	}

	groups, err := h.store.GetServersGroupedByAccount(r.Context(), status)
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load servers")
		return
	}

	data := templateData{
		Stats:        stats,
		Groups:       groups,
		Statuses:     enum.AllServerStatuses(),
		StatusFilter: statusFilter,
	}

	if err := h.tmpl.ExecuteTemplate(w, "dashboard-accounts", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleDashboardStats renders the dashboard stats partial
func (h *Handler) handleDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDashboardStats(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load stats")
		return
	}

	data := templateData{
		Stats: stats,
	}

	if err := h.tmpl.ExecuteTemplate(w, "dashboard-stats", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleThemeToggle toggles between light and dark theme
func (h *Handler) handleThemeToggle(w http.ResponseWriter, r *http.Request) {
	current := h.getTheme(r)

	var next enum.Theme
	switch current {
	case enum.ThemeLight:
		next = enum.ThemeDark
	case enum.ThemeDark:
		next = enum.ThemeSystem
	default:
		next = enum.ThemeLight
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    next.String(),
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// trigger full page reload via HX-Refresh
	w.Header().Set("HX-Refresh", "true")
}

// handleLogTable renders the log table partial
func (h *Handler) handleLogTable(w http.ResponseWriter, r *http.Request) {
	actionFilter := r.URL.Query().Get("action")

	var logs []store.ServerLogWithServer
	var err error

	if actionFilter != "" {
		action, parseErr := enum.ParseLogAction(actionFilter)
		if parseErr == nil {
			logs, err = h.store.ListLogsByAction(r.Context(), action, 100)
		}
	}
	if logs == nil {
		logs, err = h.store.ListLogs(r.Context(), 100)
	}

	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load logs")
		return
	}

	data := templateData{
		Logs:         logs,
		Actions:      enum.AllLogActions(),
		ActionFilter: actionFilter,
	}

	if err := h.tmpl.ExecuteTemplate(w, "server-logs", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

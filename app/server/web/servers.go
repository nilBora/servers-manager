package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nilBora/servers-manager/app/enum"
	"github.com/nilBora/servers-manager/app/store"
)

// handleServerTable renders the server table partial
func (h *Handler) handleServerTable(w http.ResponseWriter, r *http.Request) {
	servers, err := h.store.ListServersWithAccounts(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load servers")
		return
	}

	data := templateData{
		Servers:  servers,
		Statuses: enum.AllServerStatuses(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "server-table", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleServerForm renders the new server form
func (h *Handler) handleServerForm(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.store.ListAccountsWithProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	data := templateData{
		Accounts: accounts,
		Statuses: enum.AllServerStatuses(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "server-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleServerEditForm renders the edit server form
func (h *Handler) handleServerEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	server, err := h.store.GetServerWithAccount(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Server not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to load server")
		return
	}

	accounts, err := h.store.ListAccountsWithProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	data := templateData{
		Server:   server,
		Accounts: accounts,
		Statuses: enum.AllServerStatuses(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "server-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleServerView renders the server view modal
func (h *Handler) handleServerView(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	server, err := h.store.GetServerWithAccount(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Server not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to load server")
		return
	}

	logs, err := h.store.ListLogsByServer(r.Context(), id, 10)
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load logs")
		return
	}

	data := struct {
		Server *store.ServerWithAccount
		Logs   []store.ServerLog
	}{
		Server: server,
		Logs:   logs,
	}

	if err := h.tmpl.ExecuteTemplate(w, "server-card", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleServerCreate handles creating a new server
func (h *Handler) handleServerCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	accountID, err := strconv.ParseInt(r.FormValue("account_id"), 10, 64)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid account")
		return
	}

	status, err := enum.ParseServerStatus(r.FormValue("status"))
	if err != nil {
		status = enum.ServerStatusActive
	}

	cost, _ := strconv.ParseFloat(r.FormValue("approximate_cost"), 64)

	server := &store.Server{
		AccountID:       accountID,
		Name:            r.FormValue("name"),
		IP:              r.FormValue("ip"),
		Location:        r.FormValue("location"),
		Description:     r.FormValue("description"),
		Responsible:     r.FormValue("responsible"),
		ApproximateCost: cost,
		Backups:         r.FormValue("backups") == "on",
		Status:          status,
	}

	if server.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.CreateServer(r.Context(), server); err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to create server")
		return
	}

	// create log entry
	logEntry := &store.ServerLog{
		ServerID:    server.ID,
		Action:      enum.LogActionAdded,
		Description: "Server added",
	}
	_ = h.store.CreateLog(r.Context(), logEntry)

	// return updated table
	h.handleServerTable(w, r)
}

// handleServerUpdate handles updating an existing server
func (h *Handler) handleServerUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	accountID, err := strconv.ParseInt(r.FormValue("account_id"), 10, 64)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid account")
		return
	}

	status, err := enum.ParseServerStatus(r.FormValue("status"))
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	cost, _ := strconv.ParseFloat(r.FormValue("approximate_cost"), 64)

	server := &store.Server{
		ID:              id,
		AccountID:       accountID,
		Name:            r.FormValue("name"),
		IP:              r.FormValue("ip"),
		Location:        r.FormValue("location"),
		Description:     r.FormValue("description"),
		Responsible:     r.FormValue("responsible"),
		ApproximateCost: cost,
		Backups:         r.FormValue("backups") == "on",
		Status:          status,
	}

	if server.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.UpdateServer(r.Context(), server); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Server not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to update server")
		return
	}

	// create log entry
	logEntry := &store.ServerLog{
		ServerID:    id,
		Action:      enum.LogActionUpdated,
		Description: "Server updated",
	}
	_ = h.store.CreateLog(r.Context(), logEntry)

	// return updated table
	h.handleServerTable(w, r)
}

// handleServerStatusUpdate handles updating only server status
func (h *Handler) handleServerStatusUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	status, err := enum.ParseServerStatus(r.FormValue("status"))
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if err := h.store.UpdateServerStatus(r.Context(), id, status); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Server not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to update server status")
		return
	}

	// create log entry
	var action enum.LogAction
	switch status {
	case enum.ServerStatusPaused:
		action = enum.LogActionPaused
	case enum.ServerStatusDeleted:
		action = enum.LogActionDeleted
	default:
		action = enum.LogActionUpdated
	}
	logEntry := &store.ServerLog{
		ServerID:    id,
		Action:      action,
		Description: "Status changed to " + status.String(),
	}
	_ = h.store.CreateLog(r.Context(), logEntry)

	// return updated table
	h.handleServerTable(w, r)
}

// handleServerDelete handles deleting a server
func (h *Handler) handleServerDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.DeleteServer(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Server not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to delete server")
		return
	}

	// return updated table
	h.handleServerTable(w, r)
}

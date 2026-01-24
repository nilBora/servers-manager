package web

import (
	"errors"
	"net/http"

	"github.com/nilBora/servers-manager/app/store"
)

// handleProviderTable renders the provider table partial
func (h *Handler) handleProviderTable(w http.ResponseWriter, r *http.Request) {
	providers, err := h.store.ListProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load providers")
		return
	}

	data := templateData{
		Providers: providers,
	}

	if err := h.tmpl.ExecuteTemplate(w, "provider-table", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProviderForm renders the new provider form
func (h *Handler) handleProviderForm(w http.ResponseWriter, r *http.Request) {
	data := templateData{}

	if err := h.tmpl.ExecuteTemplate(w, "provider-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProviderEditForm renders the edit provider form
func (h *Handler) handleProviderEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	provider, err := h.store.GetProvider(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Provider not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to load provider")
		return
	}

	data := templateData{
		Provider: provider,
	}

	if err := h.tmpl.ExecuteTemplate(w, "provider-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProviderCreate handles creating a new provider
func (h *Handler) handleProviderCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	provider := &store.Provider{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	if provider.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.CreateProvider(r.Context(), provider); err != nil {
		if errors.Is(err, store.ErrConflict) {
			h.renderError(w, http.StatusConflict, "Provider with this name already exists")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to create provider")
		return
	}

	// return updated table
	h.handleProviderTable(w, r)
}

// handleProviderUpdate handles updating an existing provider
func (h *Handler) handleProviderUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	provider := &store.Provider{
		ID:          id,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	if provider.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.UpdateProvider(r.Context(), provider); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Provider not found")
			return
		}
		if errors.Is(err, store.ErrConflict) {
			h.renderError(w, http.StatusConflict, "Provider with this name already exists")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to update provider")
		return
	}

	// return updated table
	h.handleProviderTable(w, r)
}

// handleProviderDelete handles deleting a provider
func (h *Handler) handleProviderDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.DeleteProvider(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Provider not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to delete provider")
		return
	}

	// return updated table
	h.handleProviderTable(w, r)
}

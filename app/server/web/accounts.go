package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nilBora/servers-manager/app/store"
)

// handleAccountTable renders the account table partial
func (h *Handler) handleAccountTable(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.store.ListAccountsWithProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	data := templateData{
		Accounts: accounts,
	}

	if err := h.tmpl.ExecuteTemplate(w, "account-table", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAccountForm renders the new account form
func (h *Handler) handleAccountForm(w http.ResponseWriter, r *http.Request) {
	providers, err := h.store.ListProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load providers")
		return
	}

	data := templateData{
		Providers: providers,
	}

	if err := h.tmpl.ExecuteTemplate(w, "account-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAccountEditForm renders the edit account form
func (h *Handler) handleAccountEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.store.GetAccountWithProvider(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Account not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to load account")
		return
	}

	providers, err := h.store.ListProviders(r.Context())
	if err != nil {
		h.renderError(w, http.StatusInternalServerError, "Failed to load providers")
		return
	}

	data := templateData{
		Account:   account,
		Providers: providers,
	}

	if err := h.tmpl.ExecuteTemplate(w, "account-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAccountCreate handles creating a new account
func (h *Handler) handleAccountCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	providerID, err := strconv.ParseInt(r.FormValue("provider_id"), 10, 64)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid provider")
		return
	}

	account := &store.Account{
		ProviderID: providerID,
		Name:       r.FormValue("name"),
		Login:      r.FormValue("login"),
		ApiKey:     r.FormValue("api_key"),
	}

	if account.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.CreateAccount(r.Context(), account); err != nil {
		if errors.Is(err, store.ErrConflict) {
			h.renderError(w, http.StatusConflict, "Account with this name already exists for this provider")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	// return updated table
	h.handleAccountTable(w, r)
}

// handleAccountUpdate handles updating an existing account
func (h *Handler) handleAccountUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	providerID, err := strconv.ParseInt(r.FormValue("provider_id"), 10, 64)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "Invalid provider")
		return
	}

	account := &store.Account{
		ID:         id,
		ProviderID: providerID,
		Name:       r.FormValue("name"),
		Login:      r.FormValue("login"),
		ApiKey:     r.FormValue("api_key"),
	}

	if account.Name == "" {
		h.renderError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.store.UpdateAccount(r.Context(), account); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Account not found")
			return
		}
		if errors.Is(err, store.ErrConflict) {
			h.renderError(w, http.StatusConflict, "Account with this name already exists for this provider")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to update account")
		return
	}

	// return updated table
	h.handleAccountTable(w, r)
}

// handleAccountDelete handles deleting an account
func (h *Handler) handleAccountDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		h.renderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.DeleteAccount(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.renderError(w, http.StatusNotFound, "Account not found")
			return
		}
		h.renderError(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	// return updated table
	h.handleAccountTable(w, r)
}

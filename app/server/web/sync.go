package web

import (
	"net/http"

	log "github.com/go-pkgz/lgr"

	"github.com/nilBora/servers-manager/app/enum"
	"github.com/nilBora/servers-manager/app/hetzner"
	"github.com/nilBora/servers-manager/app/store"
)

// handleHetznerSync syncs servers from Hetzner Cloud accounts
func (h *Handler) handleHetznerSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all accounts with Hetzner provider
	accounts, err := h.store.ListAccountsWithProviders(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to list accounts: %v", err)
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	var synced int
	var errors []string

	for _, acc := range accounts {
		// Only sync Hetzner accounts with API keys
		if acc.ProviderName != "Hetzner" || acc.ApiKey == "" {
			continue
		}

		log.Printf("[INFO] syncing Hetzner account: %s", acc.Name)

		client := hetzner.NewClient(acc.ApiKey)
		servers, err := client.ListServers(ctx)
		if err != nil {
			log.Printf("[ERROR] failed to fetch servers from Hetzner account %s: %v", acc.Name, err)
			errors = append(errors, acc.Name+": "+err.Error())
			continue
		}

		for _, srv := range servers {
			// Map Hetzner status to our status
			status := mapHetznerStatus(srv.Status)

			// Get server details from API response
			ip := srv.GetServerIP()
			location := srv.GetServerLocationDescription()
			desc := srv.GetDescription()
			cost := srv.GetMonthlyPrice()

			// Check if server already exists (by name + account)
			existing, err := h.store.FindServerByNameAndAccount(ctx, srv.Name, acc.ID)
			if err != nil && err != store.ErrNotFound {
				log.Printf("[ERROR] failed to check existing server: %v", err)
				continue
			}

			if existing != nil {
				// Update existing server
				existing.IP = ip
				existing.Location = location
				existing.Description = desc
				existing.Status = status
				existing.ApproximateCost = cost

				if err := h.store.UpdateServer(ctx, existing); err != nil {
					log.Printf("[ERROR] failed to update server %s: %v", srv.Name, err)
					continue
				}

				// Log sync action
				logEntry := &store.ServerLog{
					ServerID:    existing.ID,
					Action:      enum.LogActionSynced,
					Description: "Synced from Hetzner Cloud",
				}
				_ = h.store.CreateLog(ctx, logEntry)
			} else {
				// Create new server
				newServer := &store.Server{
					AccountID:       acc.ID,
					Name:            srv.Name,
					IP:              ip,
					Location:        location,
					Description:     desc,
					Responsible:     "",
					ApproximateCost: cost,
					Status:          status,
				}

				if err := h.store.CreateServer(ctx, newServer); err != nil {
					log.Printf("[ERROR] failed to create server %s: %v", srv.Name, err)
					continue
				}

				// Log add action
				logEntry := &store.ServerLog{
					ServerID:    newServer.ID,
					Action:      enum.LogActionAdded,
					Description: "Added from Hetzner Cloud sync",
				}
				_ = h.store.CreateLog(ctx, logEntry)
			}

			synced++
		}
	}

	log.Printf("[INFO] Hetzner sync completed: %d servers synced", synced)

	// Return updated server table
	h.handleServerTable(w, r)
}

// mapHetznerStatus maps Hetzner server status to our status
func mapHetznerStatus(hetznerStatus string) enum.ServerStatus {
	switch hetznerStatus {
	case "running":
		return enum.ServerStatusActive
	case "off", "stopped":
		return enum.ServerStatusPaused
	case "deleting":
		return enum.ServerStatusDeleted
	default:
		return enum.ServerStatusActive
	}
}

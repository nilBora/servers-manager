package web

import (
	"context"
	"net/http"

	log "github.com/go-pkgz/lgr"

	"github.com/nilBora/servers-manager/app/enum"
	"github.com/nilBora/servers-manager/app/hetzner"
	"github.com/nilBora/servers-manager/app/store"
)

// Provider idents for sync
const (
	ProviderIdentHetznerCloud = "hetzner_cloud"
	ProviderIdentHetznerRobot = "hetzner_robot"
)

// handleHetznerSync syncs servers from all Hetzner accounts (both Cloud and Robot)
func (h *Handler) handleHetznerSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all accounts with provider info
	accounts, err := h.store.ListAccountsWithProviders(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to list accounts: %v", err)
		h.renderError(w, http.StatusInternalServerError, "Failed to load accounts")
		return
	}

	var synced int

	for _, acc := range accounts {
		// Skip accounts without API keys
		if acc.ApiKey == "" {
			continue
		}

		switch acc.ProviderIdent {
		case ProviderIdentHetznerCloud:
			count, err := h.syncHetznerCloud(ctx, &acc)
			if err != nil {
				log.Printf("[ERROR] failed to sync Hetzner Cloud account %s: %v", acc.Name, err)
				continue
			}
			synced += count

		case ProviderIdentHetznerRobot:
			count, err := h.syncHetznerRobot(ctx, &acc)
			if err != nil {
				log.Printf("[ERROR] failed to sync Hetzner Robot account %s: %v", acc.Name, err)
				continue
			}
			synced += count
		}
	}

	log.Printf("[INFO] Hetzner sync completed: %d servers synced", synced)

	// Return updated server table
	h.handleServerTable(w, r)
}

// syncHetznerCloud syncs servers from Hetzner Cloud API
func (h *Handler) syncHetznerCloud(ctx context.Context, acc *store.AccountWithProvider) (int, error) {
	log.Printf("[INFO] syncing Hetzner Cloud account: %s", acc.Name)

	client := hetzner.NewClient(acc.ApiKey)
	servers, err := client.ListServers(ctx)
	if err != nil {
		return 0, err
	}

	var synced int
	for _, srv := range servers {
		status := mapHetznerCloudStatus(srv.Status)
		ip := srv.GetServerIP()
		location := srv.GetServerLocationDescription()
		desc := srv.GetDescription()
		cost := srv.GetMonthlyPrice()

		existing, err := h.store.FindServerByNameAndAccount(ctx, srv.Name, acc.ID)
		if err != nil && err != store.ErrNotFound {
			log.Printf("[ERROR] failed to check existing server: %v", err)
			continue
		}

		if existing != nil {
			existing.IP = ip
			existing.Location = location
			existing.Description = desc
			existing.Status = status
			existing.ApproximateCost = cost

			if err := h.store.UpdateServer(ctx, existing); err != nil {
				log.Printf("[ERROR] failed to update server %s: %v", srv.Name, err)
				continue
			}

			logEntry := &store.ServerLog{
				ServerID:    existing.ID,
				Action:      enum.LogActionSynced,
				Description: "Synced from Hetzner Cloud",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		} else {
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

			logEntry := &store.ServerLog{
				ServerID:    newServer.ID,
				Action:      enum.LogActionAdded,
				Description: "Added from Hetzner Cloud sync",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		}

		synced++
	}

	return synced, nil
}

// syncHetznerRobot syncs servers from Hetzner Robot API (dedicated servers)
func (h *Handler) syncHetznerRobot(ctx context.Context, acc *store.AccountWithProvider) (int, error) {
	log.Printf("[INFO] syncing Hetzner Robot account: %s", acc.Name)

	client, err := hetzner.NewRobotClient(acc.ApiKey)
	if err != nil {
		return 0, err
	}

	servers, err := client.ListServers(ctx)
	if err != nil {
		return 0, err
	}

	var synced int
	for _, srv := range servers {
		status := mapHetznerRobotStatus(srv.Status, srv.Cancelled)
		ip := srv.GetServerIP()
		location := srv.GetServerLocation()
		desc := srv.GetDescription()
		name := srv.GetServerName()

		existing, err := h.store.FindServerByNameAndAccount(ctx, name, acc.ID)
		if err != nil && err != store.ErrNotFound {
			log.Printf("[ERROR] failed to check existing server: %v", err)
			continue
		}

		if existing != nil {
			existing.IP = ip
			existing.Location = location
			existing.Description = desc
			existing.Status = status

			if err := h.store.UpdateServer(ctx, existing); err != nil {
				log.Printf("[ERROR] failed to update server %s: %v", name, err)
				continue
			}

			logEntry := &store.ServerLog{
				ServerID:    existing.ID,
				Action:      enum.LogActionSynced,
				Description: "Synced from Hetzner Robot",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		} else {
			newServer := &store.Server{
				AccountID:       acc.ID,
				Name:            name,
				IP:              ip,
				Location:        location,
				Description:     desc,
				Responsible:     "",
				ApproximateCost: 0, // Robot API doesn't provide pricing
				Status:          status,
			}

			if err := h.store.CreateServer(ctx, newServer); err != nil {
				log.Printf("[ERROR] failed to create server %s: %v", name, err)
				continue
			}

			logEntry := &store.ServerLog{
				ServerID:    newServer.ID,
				Action:      enum.LogActionAdded,
				Description: "Added from Hetzner Robot sync",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		}

		synced++
	}

	return synced, nil
}

// mapHetznerCloudStatus maps Hetzner Cloud server status to our status
func mapHetznerCloudStatus(hetznerStatus string) enum.ServerStatus {
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

// mapHetznerRobotStatus maps Hetzner Robot server status to our status
func mapHetznerRobotStatus(status string, cancelled bool) enum.ServerStatus {
	if cancelled {
		return enum.ServerStatusDeleted
	}
	switch status {
	case "ready":
		return enum.ServerStatusActive
	case "in process":
		return enum.ServerStatusPaused
	default:
		return enum.ServerStatusActive
	}
}

package web

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"

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

// findExistingServer looks up an existing server by IP first, then by name within the same account.
// Returns nil if not found.
func (h *Handler) findExistingServer(ctx context.Context, ip, name string, accountID int64) *store.Server {
	// match by IP first (most reliable identifier)
	if ip != "" {
		existing, err := h.store.FindServerByIPAndAccount(ctx, ip, accountID)
		if err == nil {
			return existing
		}
	}

	// fallback to name
	existing, err := h.store.FindServerByNameAndAccount(ctx, name, accountID)
	if err == nil {
		return existing
	}

	return nil
}

// markDeletedServers marks servers that are in DB but not in API response as deleted.
// seenIDs contains the IDs of servers that were found in the API.
func (h *Handler) markDeletedServers(ctx context.Context, accountID int64, seenIDs map[int64]bool) {
	dbServers, err := h.store.ListServersByAccount(ctx, accountID)
	if err != nil {
		log.Printf("[ERROR] failed to list servers for account %d: %v", accountID, err)
		return
	}

	for _, srv := range dbServers {
		if seenIDs[srv.ID] {
			continue // server exists in API, skip
		}
		if srv.Status == enum.ServerStatusDeleted {
			continue // already deleted, skip
		}

		srv.Status = enum.ServerStatusDeleted
		if err := h.store.UpdateServer(ctx, &srv); err != nil {
			log.Printf("[ERROR] failed to mark server %s as deleted: %v", srv.Name, err)
			continue
		}

		logEntry := &store.ServerLog{
			ServerID:    srv.ID,
			Action:      enum.LogActionDeleted,
			Description: "Server no longer found in API, marked as deleted",
		}
		_ = h.store.CreateLog(ctx, logEntry)
		log.Printf("[INFO] marked server %s (IP: %s) as deleted — not found in API", srv.Name, srv.IP)
	}
}

// serverChanges compares existing server with new values and returns a human-readable diff.
// Returns empty string if nothing changed.
func serverChanges(existing *store.Server, name, ip, location, desc string, cost float64, backups bool, status enum.ServerStatus) string {
	var changes []string

	if existing.Name != name {
		changes = append(changes, fmt.Sprintf("Name: %s → %s", existing.Name, name))
	}
	if existing.IP != ip {
		changes = append(changes, fmt.Sprintf("IP: %s → %s", existing.IP, ip))
	}
	if existing.Location != location {
		changes = append(changes, fmt.Sprintf("Location: %s → %s", existing.Location, location))
	}
	if existing.Status != status {
		changes = append(changes, fmt.Sprintf("Status: %s → %s", existing.Status.String(), status.String()))
	}
	if math.Abs(existing.ApproximateCost-cost) > 0.01 {
		changes = append(changes, fmt.Sprintf("Cost: $%.2f → $%.2f", existing.ApproximateCost, cost))
	}
	if existing.Backups != backups {
		changes = append(changes, fmt.Sprintf("Backups: %v → %v", existing.Backups, backups))
	}
	if existing.Description != desc {
		changes = append(changes, "Description updated")
	}

	return strings.Join(changes, ", ")
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
	seenIDs := make(map[int64]bool)

	for _, srv := range servers {
		status := mapHetznerCloudStatus(srv.Status)
		ip := srv.GetServerIP()
		location := srv.GetServerLocationDescription()
		desc := srv.GetDescription()
		cost := srv.GetMonthlyPrice()
		backups := srv.HasBackups()
		if backups {
			cost *= 1.2 // Hetzner charges +20% for backups
		}

		existing := h.findExistingServer(ctx, ip, srv.Name, acc.ID)

		if existing != nil {
			seenIDs[existing.ID] = true

			// detect changes before overwriting
			diff := serverChanges(existing, srv.Name, ip, location, desc, cost, backups, status)

			existing.Name = srv.Name
			existing.IP = ip
			existing.Location = location
			existing.Description = desc
			existing.Status = status
			existing.ApproximateCost = cost
			existing.Backups = backups

			if err := h.store.UpdateServer(ctx, existing); err != nil {
				log.Printf("[ERROR] failed to update server %s: %v", srv.Name, err)
				continue
			}

			// only log if something actually changed
			if diff != "" {
				logEntry := &store.ServerLog{
					ServerID:    existing.ID,
					Action:      enum.LogActionSynced,
					Description: diff,
				}
				_ = h.store.CreateLog(ctx, logEntry)
			}
		} else {
			newServer := &store.Server{
				AccountID:       acc.ID,
				Name:            srv.Name,
				IP:              ip,
				Location:        location,
				Description:     desc,
				Responsible:     "",
				ApproximateCost: cost,
				Backups:         backups,
				Status:          status,
			}

			if err := h.store.CreateServer(ctx, newServer); err != nil {
				log.Printf("[ERROR] failed to create server %s: %v", srv.Name, err)
				continue
			}

			seenIDs[newServer.ID] = true

			logEntry := &store.ServerLog{
				ServerID:    newServer.ID,
				Action:      enum.LogActionAdded,
				Description: "Added from Hetzner Cloud sync",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		}

		synced++
	}

	// Mark servers not found in API as deleted
	h.markDeletedServers(ctx, acc.ID, seenIDs)

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
	seenIDs := make(map[int64]bool)

	for _, srv := range servers {
		status := mapHetznerRobotStatus(srv.Status, srv.Cancelled)
		ip := srv.GetServerIP()
		location := srv.GetServerLocation()
		desc := srv.GetDescription()
		name := srv.GetServerName()

		existing := h.findExistingServer(ctx, ip, name, acc.ID)

		if existing != nil {
			seenIDs[existing.ID] = true

			diff := serverChanges(existing, name, ip, location, desc, existing.ApproximateCost, existing.Backups, status)

			existing.Name = name
			existing.IP = ip
			existing.Location = location
			existing.Description = desc
			existing.Status = status

			if err := h.store.UpdateServer(ctx, existing); err != nil {
				log.Printf("[ERROR] failed to update server %s: %v", name, err)
				continue
			}

			if diff != "" {
				logEntry := &store.ServerLog{
					ServerID:    existing.ID,
					Action:      enum.LogActionSynced,
					Description: diff,
				}
				_ = h.store.CreateLog(ctx, logEntry)
			}
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

			seenIDs[newServer.ID] = true

			logEntry := &store.ServerLog{
				ServerID:    newServer.ID,
				Action:      enum.LogActionAdded,
				Description: "Added from Hetzner Robot sync",
			}
			_ = h.store.CreateLog(ctx, logEntry)
		}

		synced++
	}

	// Mark servers not found in API as deleted
	h.markDeletedServers(ctx, acc.ID, seenIDs)

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

package hetzner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const baseURL = "https://api.hetzner.cloud/v1"

// Client is a Hetzner Cloud API client
type Client struct {
	httpClient *http.Client
	apiToken   string
}

// NewClient creates a new Hetzner API client
func NewClient(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiToken: apiToken,
	}
}

// Server represents a Hetzner server
type Server struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	PublicNet  PublicNet  `json:"public_net"`
	ServerType ServerType `json:"server_type"`
	Datacenter Datacenter `json:"datacenter"`
	Location   Location   `json:"location"`
	Image      *Image     `json:"image"`
}

// PublicNet contains public network information
type PublicNet struct {
	IPv4 IPv4 `json:"ipv4"`
	IPv6 IPv6 `json:"ipv6"`
}

// IPv4 contains IPv4 address info
type IPv4 struct {
	IP string `json:"ip"`
}

// IPv6 contains IPv6 address info
type IPv6 struct {
	IP string `json:"ip"`
}

// ServerType contains server type info
type ServerType struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cores       int           `json:"cores"`
	Memory      float64       `json:"memory"`
	Disk        int           `json:"disk"`
	Prices      []ServerPrice `json:"prices"`
}

// ServerPrice contains pricing info for a location
type ServerPrice struct {
	Location     string      `json:"location"`
	PriceHourly  PriceAmount `json:"price_hourly"`
	PriceMonthly PriceAmount `json:"price_monthly"`
}

// PriceAmount contains gross and net prices
type PriceAmount struct {
	Gross string `json:"gross"`
	Net   string `json:"net"`
}

// Datacenter contains datacenter info
type Datacenter struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
}

// Location contains location info
type Location struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	NetworkZone string  `json:"network_zone"`
}

// Image contains OS image info
type Image struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OSFlavor    string `json:"os_flavor"`
	OSVersion   string `json:"os_version"`
}

// ListServersResponse is the API response for listing servers
type ListServersResponse struct {
	Servers []Server `json:"servers"`
	Meta    Meta     `json:"meta"`
}

// Meta contains pagination info
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination contains pagination details
type Pagination struct {
	Page         int  `json:"page"`
	PerPage      int  `json:"per_page"`
	PreviousPage *int `json:"previous_page"`
	NextPage     *int `json:"next_page"`
	LastPage     int  `json:"last_page"`
	TotalEntries int  `json:"total_entries"`
}

// ListServers fetches all servers from Hetzner Cloud
func (c *Client) ListServers(ctx context.Context) ([]Server, error) {
	var allServers []Server
	page := 1

	for {
		url := fmt.Sprintf("%s/servers?page=%d&per_page=50", baseURL, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
		}

		var result ListServersResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}

		allServers = append(allServers, result.Servers...)

		// Check if there are more pages
		if result.Meta.Pagination.NextPage == nil {
			break
		}
		page = *result.Meta.Pagination.NextPage
	}

	return allServers, nil
}

// GetServerIP returns the primary IP address (IPv4 preferred)
func (s *Server) GetServerIP() string {
	if s.PublicNet.IPv4.IP != "" {
		return s.PublicNet.IPv4.IP
	}
	return s.PublicNet.IPv6.IP
}

// GetServerLocation returns the location name (e.g., "nbg1", "fsn1")
func (s *Server) GetServerLocation() string {
	// Prefer datacenter location, fallback to server location
	if s.Datacenter.Location.Name != "" {
		return s.Datacenter.Location.Name
	}
	return s.Location.Name
}

// GetServerLocationDescription returns human-readable location description
func (s *Server) GetServerLocationDescription() string {
	loc := s.Datacenter.Location
	if loc.City != "" && loc.Country != "" {
		return fmt.Sprintf("%s, %s", loc.City, loc.Country)
	}
	if loc.Description != "" {
		return loc.Description
	}
	return s.GetServerLocation()
}

// GetMonthlyPrice returns the monthly price for this server's location
func (s *Server) GetMonthlyPrice() float64 {
	location := s.GetServerLocation()

	for _, price := range s.ServerType.Prices {
		if price.Location == location {
			if gross, err := strconv.ParseFloat(price.PriceMonthly.Gross, 64); err == nil {
				return gross
			}
		}
	}

	// Fallback: return first available price
	if len(s.ServerType.Prices) > 0 {
		if gross, err := strconv.ParseFloat(s.ServerType.Prices[0].PriceMonthly.Gross, 64); err == nil {
			return gross
		}
	}

	return 0
}

// GetDescription returns a description combining server type and OS info
func (s *Server) GetDescription() string {
	desc := s.ServerType.Description
	if s.Image != nil && s.Image.Description != "" {
		desc += " | " + s.Image.Description
	}
	return desc
}

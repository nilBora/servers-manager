package hetzner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const robotBaseURL = "https://robot-ws.your-server.de"

// RobotClient is a Hetzner Robot API client for dedicated servers
type RobotClient struct {
	httpClient *http.Client
	username   string
	password   string
}

// NewRobotClient creates a new Hetzner Robot API client
// credentials should be in format "username:password"
func NewRobotClient(credentials string) (*RobotClient, error) {
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid credentials format, expected 'username:password'")
	}

	return &RobotClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: parts[0],
		password: parts[1],
	}, nil
}

// RobotServer represents a Hetzner Robot dedicated server
type RobotServer struct {
	ServerIP      string   `json:"server_ip"`
	ServerIPv6Net string   `json:"server_ipv6_net"`
	ServerNumber  int64    `json:"server_number"`
	ServerName    string   `json:"server_name"`
	Product       string   `json:"product"`
	DC            string   `json:"dc"`
	Traffic       string   `json:"traffic"`
	Status        string   `json:"status"`
	Cancelled     bool     `json:"cancelled"`
	PaidUntil     string   `json:"paid_until"`
	IP            []string `json:"ip"`
	Subnet        []Subnet `json:"subnet"`
}

// Subnet represents an IP subnet
type Subnet struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

// robotServerWrapper wraps the server response
type robotServerWrapper struct {
	Server RobotServer `json:"server"`
}

// ListServers fetches all dedicated servers from Hetzner Robot
func (c *RobotClient) ListServers(ctx context.Context) ([]RobotServer, error) {
	url := robotBaseURL + "/server"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: invalid username or password")
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("access forbidden: rate limit exceeded or IP blocked")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var wrappers []robotServerWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrappers); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	servers := make([]RobotServer, len(wrappers))
	for i, w := range wrappers {
		servers[i] = w.Server
	}

	return servers, nil
}

// GetServerIP returns the primary IP address
func (s *RobotServer) GetServerIP() string {
	return s.ServerIP
}

// GetServerName returns the server name, falling back to product if empty
func (s *RobotServer) GetServerName() string {
	if s.ServerName != "" {
		return s.ServerName
	}
	return fmt.Sprintf("Server #%d", s.ServerNumber)
}

// GetServerLocation returns the datacenter location
func (s *RobotServer) GetServerLocation() string {
	return s.DC
}

// GetDescription returns a description with product info
func (s *RobotServer) GetDescription() string {
	desc := s.Product
	if s.Traffic != "" {
		desc += " | Traffic: " + s.Traffic
	}
	return desc
}

// IsActive returns true if server is ready and not cancelled
func (s *RobotServer) IsActive() bool {
	return s.Status == "ready" && !s.Cancelled
}

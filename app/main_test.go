package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/nilBora/servers-manager/app/server"
	"github.com/nilBora/servers-manager/app/store"
)

func TestSetupLog(t *testing.T) {
	setupLog(false)
	setupLog(true)
}

func TestServerLifecycle(t *testing.T) {
	st, addr := newTestStore(t)
	defer st.Close()

	srv, err := server.New(st, testConfig(addr))
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()

	baseURL := "http://" + addr
	if err := waitReady(baseURL+"/setup", 2*time.Second); err != nil {
		t.Fatalf("server not ready: %v", err)
	}

	// /setup is public — must return 200
	assertStatus(t, baseURL+"/setup", http.StatusOK)

	// /login is public — must return 200
	assertStatus(t, baseURL+"/login", http.StatusOK)

	// protected route without session follows redirect to /setup or /login → 200
	assertStatus(t, baseURL+"/", http.StatusOK)

	// graceful shutdown
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("srv.Run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("server did not shut down within timeout")
	}
}

func TestServerShutdownSpeed(t *testing.T) {
	st, addr := newTestStore(t)
	defer st.Close()

	cfg := testConfig(addr)
	cfg.ShutdownTimeout = 500 * time.Millisecond

	srv, err := server.New(st, cfg)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()

	if err := waitReady("http://"+addr+"/login", 2*time.Second); err != nil {
		t.Fatalf("server not ready: %v", err)
	}

	start := time.Now()
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("srv.Run: %v", err)
		}
		if elapsed := time.Since(start); elapsed > 3*time.Second {
			t.Errorf("shutdown too slow: %s", elapsed)
		}
	case <-time.After(5 * time.Second):
		t.Error("server did not shut down")
	}
}

// testConfig returns a base server config for the given address.
func testConfig(addr string) server.Config {
	return server.Config{
		Address:         addr,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 1 * time.Second,
	}
}

// newTestStore creates a store backed by a temp DB and returns a free TCP address.
func newTestStore(t *testing.T) (*store.DB, string) {
	t.Helper()

	st, err := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	return st, addr
}

// waitReady polls url until it responds or timeout is reached.
func waitReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:gosec
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(20 * time.Millisecond)
	}
	return fmt.Errorf("not ready after %s: %s", timeout, url)
}

// assertStatus makes a GET request and checks the response status code.
func assertStatus(t *testing.T, url string, want int) {
	t.Helper()
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	resp.Body.Close()
	if resp.StatusCode != want {
		t.Errorf("GET %s: want %d, got %d", url, want, resp.StatusCode)
	}
}

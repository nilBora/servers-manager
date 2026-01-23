// Package server provides HTTP server for the servers manager web UI
package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/go-pkgz/lgr"

	"github.com/nilBora/servers-manager/app/server/web"
	"github.com/nilBora/servers-manager/app/store"
)

// Server represents the HTTP server
type Server struct {
	Config
	store      store.Store
	webHandler *web.Handler
	staticFS   fs.FS
}

// Config holds server configuration
type Config struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	Version         string
}

// New creates a new Server instance
func New(st store.Store, cfg Config) (*Server, error) {
	staticContent, err := web.StaticFS()
	if err != nil {
		return nil, fmt.Errorf("failed to load static files: %w", err)
	}

	webHandler, err := web.New(st)
	if err != nil {
		return nil, fmt.Errorf("failed to create web handler: %w", err)
	}

	return &Server{
		Config:     cfg,
		store:      st,
		webHandler: webHandler,
		staticFS:   staticContent,
	}, nil
}

// Run starts the HTTP server and blocks until context is canceled
func (s *Server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:              s.Address,
		Handler:           s.routes(),
		ReadHeaderTimeout: s.ReadTimeout,
		WriteTimeout:      s.WriteTimeout,
		IdleTimeout:       s.IdleTimeout,
	}

	// graceful shutdown
	go func() {
		<-ctx.Done()
		log.Printf("[INFO] shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("[WARN] shutdown error: %v", err)
		}
	}()

	log.Printf("[INFO] started server on %s", s.Address)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// routes configures and returns the HTTP handler with all routes
func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(s.staticFS))))

	// web UI routes
	s.webHandler.Register(r)

	return r
}

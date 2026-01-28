package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"

	"github.com/nilBora/servers-manager/app/server"
	"github.com/nilBora/servers-manager/app/store"
)

// opts holds command line options
var opts struct {
	DB      string `long:"db" env:"DB" default:"servers.db" description:"database file path"`
	Address string `long:"address" env:"ADDRESS" default:":8080" description:"server address"`
	Debug   bool   `long:"debug" env:"DEBUG" description:"enable debug mode"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	setupLog(opts.Debug)
	log.Printf("[INFO] servers-manager starting")
	// initialize store
	st, err := store.New(opts.DB)
	if err != nil {
		log.Fatalf("[ERROR] failed to initialize store: %v", err)
	}
	defer st.Close()

	// create server
	srv, err := server.New(st, server.Config{
		Address:         opts.Address,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("[ERROR] failed to create server: %v", err)
	}

	// setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Printf("[INFO] received shutdown signal")
		cancel()
	}()

	// run server
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("[ERROR] server error: %v", err)
	}

	log.Printf("[INFO] servers-manager stopped")
}

func setupLog(debug bool) {
	if debug {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
	} else {
		log.Setup(log.Msec, log.LevelBraces)
	}
}

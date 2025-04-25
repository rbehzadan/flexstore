package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rbehzadan/schemaless-api/internal/app"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

// Run starts the server with the provided configuration
func Run(cfg *config.Config) {
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create server with reasonable timeouts
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      app.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	fmt.Printf("Starting Schemaless API Server v%s on %s\n", cfg.Version, cfg.Addr)
	log.Fatal(server.ListenAndServe())
}

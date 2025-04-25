package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rbehzadan/schemaless-api/cmd/server"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

//go:embed VERSION
var versionFile embed.FS

func main() {
	// Define command line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		addr        = flag.String("addr", ":8080", "HTTP service address")
	)

	flag.Parse()

	// Read version from embedded file
	versionBytes, err := versionFile.ReadFile("VERSION")
	if err != nil {
		panic(err)
	}
	version := strings.TrimSpace(string(versionBytes))

	// Handle version flag
	if *showVersion {
		fmt.Println("Schemaless API Server Version:", version)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Initialize configuration
	cfg := config.NewConfig()
	cfg.Version = version
	cfg.StartTime = time.Now()
	cfg.Addr = *addr

	// Start the server with the initialized config
	server.Run(cfg)
}

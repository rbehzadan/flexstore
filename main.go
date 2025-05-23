package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rbehzadan/flexstore/cmd/server"
	"github.com/rbehzadan/flexstore/pkg/config"
)

//go:embed VERSION
var versionFile embed.FS

func main() {
	// Define command line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		addr        = flag.String("addr", ":8080", "HTTP service address")
		enableAuth  = flag.Bool("auth", true, "Enable HTTP Basic Authentication")
		username    = flag.String("username", "admin", "Username for HTTP Basic Authentication")
		password    = flag.String("password", "password", "Password for HTTP Basic Authentication")
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
		fmt.Println("FlexStore API Server Version:", version)
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
	cfg.EnableBasicAuth = *enableAuth
	cfg.AuthUsername = *username
	cfg.AuthPassword = *password

	// Start the server with the initialized config
	server.Run(cfg)
}

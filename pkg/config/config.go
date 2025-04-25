package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	Version         string
	StartTime       time.Time
	Addr            string
	SqlitePath      string
	AuthUsername    string
	AuthPassword    string
	EnableBasicAuth bool
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Version:         "0.0.0",
		StartTime:       time.Now(),
		Addr:            ":8080",
		SqlitePath:      "data/db.sqlite",
		AuthUsername:    "admin",
		AuthPassword:    "password",
		EnableBasicAuth: false,
	}
}

// GetUptime returns the uptime of the server as a string
func (c *Config) GetUptime() string {
	return time.Since(c.StartTime).Round(time.Second).String()
}

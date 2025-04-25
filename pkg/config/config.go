package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	Version   string
	StartTime time.Time
	Addr      string
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Version:   "0.0.0",
		StartTime: time.Now(),
		Addr:      ":8080",
	}
}

// GetUptime returns the uptime of the server as a string
func (c *Config) GetUptime() string {
	return time.Since(c.StartTime).Round(time.Second).String()
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents a database connection
type DB struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Path string
}

// NewConfig creates a default database configuration
func NewConfig(sqlitePath string) *Config {
	return &Config{
		Path: sqlitePath,
	}
}

// New creates a new database connection
func New(config *Config) (*DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(config.Path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	sqlDB, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection parameters
	sqlDB.SetMaxOpenConns(1) // SQLite only supports one writer at a time

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{sqlDB}

	// Initialize database schema
	if err := db.Initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return db, nil
}

// Initialize creates the necessary database tables if they don't exist
func (db *DB) Initialize() error {
	// Schema for collections
	collections := `
	CREATE TABLE IF NOT EXISTS collections (
		name TEXT PRIMARY KEY,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	// Schema for documents
	documents := `
	CREATE TABLE IF NOT EXISTS documents (
		id TEXT,
		collection_name TEXT,
		data TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id, collection_name),
		FOREIGN KEY (collection_name) REFERENCES collections(name) ON DELETE CASCADE
	);`

	// Create collections table
	if _, err := db.Exec(collections); err != nil {
		return fmt.Errorf("failed to create collections table: %w", err)
	}

	// Create documents table
	if _, err := db.Exec(documents); err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

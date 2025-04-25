package models

import (
	"time"
)

// Collection represents a document collection
type Collection struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CollectionList represents a list of collections with metadata
type CollectionList struct {
	Total       int          `json:"total"`
	Collections []Collection `json:"collections"`
}

// NewCollection creates a new collection
func NewCollection(name string) *Collection {
	now := time.Now()
	return &Collection{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Document represents a schemaless document
type Document struct {
	ID             string          `json:"id"`
	CollectionName string          `json:"collection_name"`
	Data           json.RawMessage `json:"data"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// DocumentList represents a list of documents with metadata
type DocumentList struct {
	Total     int        `json:"total"`
	Offset    int        `json:"offset"`
	Limit     int        `json:"limit"`
	Documents []Document `json:"documents"`
}

// NewDocument creates a new document
func NewDocument(collectionName string, data json.RawMessage) *Document {
	now := time.Now()

	// Generate ID from hash of data + timestamp if not provided
	hash := sha256.Sum256(append(data, []byte(now.String())...))
	id := hex.EncodeToString(hash[:8]) // Using first 8 bytes of hash for ID

	return &Document{
		ID:             id,
		CollectionName: collectionName,
		Data:           data,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewDocumentWithID creates a new document with a specified ID
func NewDocumentWithID(id string, collectionName string, data json.RawMessage) *Document {
	now := time.Now()
	return &Document{
		ID:             id,
		CollectionName: collectionName,
		Data:           data,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// ValidateJSON ensures the provided data is valid JSON
func ValidateJSON(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON data")
	}
	return nil
}

// DocumentQuery represents query parameters for retrieving documents
type DocumentQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

// NewDocumentQuery creates a new document query with default values
func NewDocumentQuery() *DocumentQuery {
	return &DocumentQuery{
		Limit:  100,
		Offset: 0,
	}
}

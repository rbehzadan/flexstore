package service

import (
	"encoding/json"
	"io"

	"github.com/rbehzadan/schemaless-api/internal/db"
	"github.com/rbehzadan/schemaless-api/internal/models"
)

// DocumentService handles document operations
type DocumentService struct {
	repo *db.DocumentRepository
}

// NewDocumentService creates a new document service
func NewDocumentService(repo *db.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

// Create creates a new document
func (s *DocumentService) Create(collectionName string, data json.RawMessage) (*models.Document, error) {
	return s.repo.Create(collectionName, data)
}

// CreateWithID creates a new document with the specified ID
func (s *DocumentService) CreateWithID(id, collectionName string, data json.RawMessage) (*models.Document, error) {
	return s.repo.CreateWithID(id, collectionName, data)
}

// GetByID retrieves a document by ID
func (s *DocumentService) GetByID(id, collectionName string) (*models.Document, error) {
	return s.repo.GetByID(id, collectionName)
}

// Update updates a document
func (s *DocumentService) Update(id, collectionName string, data json.RawMessage) (*models.Document, error) {
	return s.repo.Update(id, collectionName, data)
}

// Delete deletes a document
func (s *DocumentService) Delete(id, collectionName string) error {
	return s.repo.Delete(id, collectionName)
}

// List retrieves documents from a collection with pagination
func (s *DocumentService) List(collectionName string, queryParams *models.DocumentQuery) (*models.DocumentList, error) {
	return s.repo.List(collectionName, queryParams)
}

// BulkCreate creates multiple documents in a collection
func (s *DocumentService) BulkCreate(collectionName string, dataItems []json.RawMessage) ([]models.Document, error) {
	return s.repo.BulkCreate(collectionName, dataItems)
}

// ProcessJSONFile processes a JSON file for bulk insertion
func (s *DocumentService) ProcessJSONFile(collectionName string, r io.Reader) ([]models.Document, error) {
	// Read the file content
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Check if data is a JSON array
	var jsonArray []json.RawMessage
	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		// Try as a single JSON object
		var jsonObj map[string]interface{}
		err = json.Unmarshal(data, &jsonObj)
		if err != nil {
			return nil, err
		}

		// Convert single object to array
		jsonData, err := json.Marshal(jsonObj)
		if err != nil {
			return nil, err
		}
		jsonArray = []json.RawMessage{jsonData}
	}

	// Bulk create documents
	return s.BulkCreate(collectionName, jsonArray)
}

package service

import (
	"github.com/rbehzadan/flexstore/internal/db"
	"github.com/rbehzadan/flexstore/internal/models"
)

// CollectionService handles collection operations
type CollectionService struct {
	repo *db.CollectionRepository
}

// NewCollectionService creates a new collection service
func NewCollectionService(repo *db.CollectionRepository) *CollectionService {
	return &CollectionService{repo: repo}
}

// Create creates a new collection
func (s *CollectionService) Create(name string) (*models.Collection, error) {
	return s.repo.Create(name)
}

// GetByName retrieves a collection by name
func (s *CollectionService) GetByName(name string) (*models.Collection, error) {
	return s.repo.GetByName(name)
}

// Delete deletes a collection
func (s *CollectionService) Delete(name string) error {
	return s.repo.Delete(name)
}

// List retrieves all collections
func (s *CollectionService) List() (*models.CollectionList, error) {
	return s.repo.List()
}

// Exists checks if a collection exists
func (s *CollectionService) Exists(name string) (bool, error) {
	return s.repo.Exists(name)
}

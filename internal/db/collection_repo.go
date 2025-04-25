package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rbehzadan/schemaless-api/internal/models"
)

// CollectionRepository handles collection operations
type CollectionRepository struct {
	db *DB
}

// NewCollectionRepository creates a new collection repository
func NewCollectionRepository(db *DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create creates a new collection
func (r *CollectionRepository) Create(name string) (*models.Collection, error) {
	// Check if collection already exists
	exists, err := r.Exists(name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("collection '%s' already exists", name)
	}

	// Create collection
	collection := models.NewCollection(name)

	// Insert collection into database
	query := `INSERT INTO collections (name, created_at, updated_at) VALUES (?, ?, ?)`
	_, err = r.db.Exec(query, collection.Name, collection.CreatedAt, collection.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

// GetByName retrieves a collection by name
func (r *CollectionRepository) GetByName(name string) (*models.Collection, error) {
	query := `SELECT name, created_at, updated_at FROM collections WHERE name = ?`
	row := r.db.QueryRow(query, name)

	var collection models.Collection
	err := row.Scan(&collection.Name, &collection.CreatedAt, &collection.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("collection '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &collection, nil
}

// Exists checks if a collection exists
func (r *CollectionRepository) Exists(name string) (bool, error) {
	query := `SELECT 1 FROM collections WHERE name = ?`
	row := r.db.QueryRow(query, name)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if collection exists: %w", err)
	}

	return true, nil
}

// Delete deletes a collection
func (r *CollectionRepository) Delete(name string) error {
	// Check if collection exists
	exists, err := r.Exists(name)
	if err != nil {
		return fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("collection '%s' not found", name)
	}

	// Delete collection
	query := `DELETE FROM collections WHERE name = ?`
	_, err = r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// List retrieves all collections
func (r *CollectionRepository) List() (*models.CollectionList, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM collections`
	var total int
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count collections: %w", err)
	}

	// Get collections
	query := `SELECT name, created_at, updated_at FROM collections ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	defer rows.Close()

	collections := make([]models.Collection, 0)
	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(&collection.Name, &collection.CreatedAt, &collection.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over collections: %w", err)
	}

	return &models.CollectionList{
		Total:       total,
		Collections: collections,
	}, nil
}

// Update updates a collection
func (r *CollectionRepository) Update(name string) (*models.Collection, error) {
	// Check if collection exists
	collection, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	// Update timestamp
	collection.UpdatedAt = time.Now()

	// Update collection
	query := `UPDATE collections SET updated_at = ? WHERE name = ?`
	_, err = r.db.Exec(query, collection.UpdatedAt, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	return collection, nil
}

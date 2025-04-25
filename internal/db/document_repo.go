package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rbehzadan/schemaless-api/internal/models"
)

// DocumentRepository handles document operations
type DocumentRepository struct {
	db             *DB
	collectionRepo *CollectionRepository
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *DB, collectionRepo *CollectionRepository) *DocumentRepository {
	return &DocumentRepository{
		db:             db,
		collectionRepo: collectionRepo,
	}
}

// Create creates a new document
func (r *DocumentRepository) Create(collectionName string, data json.RawMessage) (*models.Document, error) {
	// Check if collection exists
	exists, err := r.collectionRepo.Exists(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection exists: %w", err)
	}

	// If collection doesn't exist, create it first
	if !exists {
		_, err := r.collectionRepo.Create(collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	// Validate JSON data
	if err := models.ValidateJSON(data); err != nil {
		return nil, err
	}

	// Create document
	document := models.NewDocument(collectionName, data)

	// Insert document into database
	query := `INSERT INTO documents (id, collection_name, data, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err = r.db.Exec(
		query,
		document.ID,
		document.CollectionName,
		document.Data,
		document.CreatedAt,
		document.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Update collection timestamp
	_, err = r.collectionRepo.Update(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection timestamp: %w", err)
	}

	return document, nil
}

// CreateWithID creates a new document with the specified ID
func (r *DocumentRepository) CreateWithID(id, collectionName string, data json.RawMessage) (*models.Document, error) {
	// Check if collection exists
	exists, err := r.collectionRepo.Exists(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", collectionName)
	}

	// Check if document with this ID already exists
	exists, err = r.Exists(id, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if document exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("document with ID '%s' already exists in collection '%s'", id, collectionName)
	}

	// Validate JSON data
	if err := models.ValidateJSON(data); err != nil {
		return nil, err
	}

	// Create document
	document := models.NewDocumentWithID(id, collectionName, data)

	// Insert document into database
	query := `INSERT INTO documents (id, collection_name, data, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err = r.db.Exec(
		query,
		document.ID,
		document.CollectionName,
		document.Data,
		document.CreatedAt,
		document.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Update collection timestamp
	_, err = r.collectionRepo.Update(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection timestamp: %w", err)
	}

	return document, nil
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(id, collectionName string) (*models.Document, error) {
	query := `SELECT id, collection_name, data, created_at, updated_at 
			  FROM documents 
			  WHERE id = ? AND collection_name = ?`
	row := r.db.QueryRow(query, id, collectionName)

	var document models.Document
	var dataBytes []byte
	err := row.Scan(
		&document.ID,
		&document.CollectionName,
		&dataBytes,
		&document.CreatedAt,
		&document.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("document with ID '%s' not found in collection '%s'", id, collectionName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	document.Data = json.RawMessage(dataBytes)
	return &document, nil
}

// Exists checks if a document exists
func (r *DocumentRepository) Exists(id, collectionName string) (bool, error) {
	query := `SELECT 1 FROM documents WHERE id = ? AND collection_name = ?`
	row := r.db.QueryRow(query, id, collectionName)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if document exists: %w", err)
	}

	return true, nil
}

// Update updates a document
func (r *DocumentRepository) Update(id, collectionName string, data json.RawMessage) (*models.Document, error) {
	// Check if document exists
	exists, err := r.Exists(id, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if document exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("document with ID '%s' not found in collection '%s'", id, collectionName)
	}

	// Validate JSON data
	if err := models.ValidateJSON(data); err != nil {
		return nil, err
	}

	// Update document
	now := time.Now()
	query := `UPDATE documents 
			  SET data = ?, updated_at = ? 
			  WHERE id = ? AND collection_name = ?`
	_, err = r.db.Exec(query, data, now, id, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Update collection timestamp
	_, err = r.collectionRepo.Update(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection timestamp: %w", err)
	}

	// Get updated document
	return r.GetByID(id, collectionName)
}

// Delete deletes a document
func (r *DocumentRepository) Delete(id, collectionName string) error {
	// Check if document exists
	exists, err := r.Exists(id, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check if document exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("document with ID '%s' not found in collection '%s'", id, collectionName)
	}

	// Delete document
	query := `DELETE FROM documents WHERE id = ? AND collection_name = ?`
	_, err = r.db.Exec(query, id, collectionName)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Update collection timestamp
	_, err = r.collectionRepo.Update(collectionName)
	if err != nil {
		return fmt.Errorf("failed to update collection timestamp: %w", err)
	}

	return nil
}

// List retrieves documents from a collection with pagination
func (r *DocumentRepository) List(collectionName string, queryParams *models.DocumentQuery) (*models.DocumentList, error) {
	// Check if collection exists
	exists, err := r.collectionRepo.Exists(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", collectionName)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM documents WHERE collection_name = ?`
	var total int
	err = r.db.QueryRow(countQuery, collectionName).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Get documents with pagination
	query := `SELECT id, collection_name, data, created_at, updated_at 
			  FROM documents 
			  WHERE collection_name = ? 
			  ORDER BY created_at DESC 
			  LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, collectionName, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var document models.Document
		var dataBytes []byte
		err := rows.Scan(
			&document.ID,
			&document.CollectionName,
			&dataBytes,
			&document.CreatedAt,
			&document.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		document.Data = json.RawMessage(dataBytes)
		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over documents: %w", err)
	}

	return &models.DocumentList{
		Total:     total,
		Offset:    queryParams.Offset,
		Limit:     queryParams.Limit,
		Documents: documents,
	}, nil
}

// BulkCreate creates multiple documents in a collection
func (r *DocumentRepository) BulkCreate(collectionName string, dataItems []json.RawMessage) ([]models.Document, error) {
	// Check if collection exists
	exists, err := r.collectionRepo.Exists(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", collectionName)
	}

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare statement for inserting documents
	stmt, err := tx.Prepare(`INSERT INTO documents (id, collection_name, data, created_at, updated_at) 
							 VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert each document
	documents := make([]models.Document, 0, len(dataItems))
	for _, data := range dataItems {
		// Validate JSON
		if err := models.ValidateJSON(data); err != nil {
			return nil, err
		}

		// Create document
		document := models.NewDocument(collectionName, data)
		documents = append(documents, *document)

		// Insert document
		_, err = stmt.Exec(
			document.ID,
			document.CollectionName,
			document.Data,
			document.CreatedAt,
			document.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert document: %w", err)
		}
	}

	// Update collection timestamp
	now := time.Now()
	_, err = tx.Exec(`UPDATE collections SET updated_at = ? WHERE name = ?`, now, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection timestamp: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return documents, nil
}

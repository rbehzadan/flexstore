package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rbehzadan/flexstore/internal/api"
	"github.com/rbehzadan/flexstore/internal/models"
	"github.com/rbehzadan/flexstore/internal/service"
)

// DocumentHandlers contains handlers for document operations
type DocumentHandlers struct {
	documentService *service.DocumentService
}

// NewDocumentHandlers creates new document handlers
func NewDocumentHandlers(documentService *service.DocumentService) *DocumentHandlers {
	return &DocumentHandlers{
		documentService: documentService,
	}
}

// CreateDocument creates a new document
func (h *DocumentHandlers) CreateDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]

		// Read request body
		var data json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON data")
			return
		}

		// Create document
		document, err := h.documentService.Create(collectionName, data)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "CREATE_DOCUMENT_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusCreated, document)
	}
}

// GetDocument gets a document by ID
func (h *DocumentHandlers) GetDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name and document ID from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]
		id := vars["id"]

		// Get document
		document, err := h.documentService.GetByID(id, collectionName)
		if err != nil {
			api.RespondWithError(w, http.StatusNotFound, "DOCUMENT_NOT_FOUND", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, document)
	}
}

// UpdateDocument updates a document
func (h *DocumentHandlers) UpdateDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name and document ID from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]
		id := vars["id"]

		// Read request body
		var data json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON data")
			return
		}

		// Update document
		document, err := h.documentService.Update(id, collectionName, data)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "UPDATE_DOCUMENT_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, document)
	}
}

// DeleteDocument deletes a document
func (h *DocumentHandlers) DeleteDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name and document ID from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]
		id := vars["id"]

		// Delete document
		err := h.documentService.Delete(id, collectionName)
		if err != nil {
			api.RespondWithError(w, http.StatusNotFound, "DOCUMENT_NOT_FOUND", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
	}
}

// ListDocuments lists documents in a collection
func (h *DocumentHandlers) ListDocuments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]

		// Parse query parameters
		query := models.NewDocumentQuery()

		// Get limit parameter
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err == nil && limit > 0 {
				query.Limit = limit
			}
		}

		// Get offset parameter
		offsetStr := r.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err == nil && offset >= 0 {
				query.Offset = offset
			}
		}

		// Get documents
		documents, err := h.documentService.List(collectionName, query)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "LIST_DOCUMENTS_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, documents)
	}
}

// BulkCreateDocuments creates multiple documents in a collection
func (h *DocumentHandlers) BulkCreateDocuments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]

		// Read request body
		var dataItems []json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&dataItems); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON array")
			return
		}

		// Create documents
		documents, err := h.documentService.BulkCreate(collectionName, dataItems)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "BULK_CREATE_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message":   "Documents created successfully",
			"count":     len(documents),
			"documents": documents,
		})
	}
}

// UploadJSONFile uploads and processes a JSON file
func (h *DocumentHandlers) UploadJSONFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		collectionName := vars["name"]

		// Parse form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_FORM", "Could not parse form")
			return
		}

		// Get file
		file, _, err := r.FormFile("file")
		if err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_FILE", "Could not get file from form")
			return
		}
		defer file.Close()

		// Process file
		documents, err := h.documentService.ProcessJSONFile(collectionName, file)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "PROCESS_FILE_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message":   "File processed successfully",
			"count":     len(documents),
			"documents": documents,
		})
	}
}

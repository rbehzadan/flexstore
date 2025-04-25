package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rbehzadan/flexstore/internal/api"
	"github.com/rbehzadan/flexstore/internal/service"
)

// CollectionHandlers contains handlers for collection operations
type CollectionHandlers struct {
	collectionService *service.CollectionService
}

// NewCollectionHandlers creates new collection handlers
func NewCollectionHandlers(collectionService *service.CollectionService) *CollectionHandlers {
	return &CollectionHandlers{
		collectionService: collectionService,
	}
}

// CreateCollection creates a new collection
func (h *CollectionHandlers) CreateCollection() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		// Validate name
		if req.Name == "" {
			api.RespondWithError(w, http.StatusBadRequest, "INVALID_NAME", "Collection name cannot be empty")
			return
		}

		// Create collection
		collection, err := h.collectionService.Create(req.Name)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "CREATE_COLLECTION_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusCreated, collection)
	}
}

// GetCollection gets a collection by name
func (h *CollectionHandlers) GetCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		name := vars["name"]

		// Get collection
		collection, err := h.collectionService.GetByName(name)
		if err != nil {
			api.RespondWithError(w, http.StatusNotFound, "COLLECTION_NOT_FOUND", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, collection)
	}
}

// DeleteCollection deletes a collection
func (h *CollectionHandlers) DeleteCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collection name from URL
		vars := mux.Vars(r)
		name := vars["name"]

		// Delete collection
		err := h.collectionService.Delete(name)
		if err != nil {
			api.RespondWithError(w, http.StatusNotFound, "COLLECTION_NOT_FOUND", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Collection deleted successfully"})
	}
}

// ListCollections lists all collections
func (h *CollectionHandlers) ListCollections() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get collections
		collections, err := h.collectionService.List()
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "LIST_COLLECTIONS_ERROR", err.Error())
			return
		}

		// Respond
		api.RespondWithJSON(w, http.StatusOK, collections)
	}
}

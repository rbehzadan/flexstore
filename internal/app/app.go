package app

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rbehzadan/schemaless-api/internal/api/handlers"
	"github.com/rbehzadan/schemaless-api/internal/api/middleware"
	"github.com/rbehzadan/schemaless-api/internal/db"
	"github.com/rbehzadan/schemaless-api/internal/service"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

// App represents the application
type App struct {
	Router            *mux.Router
	DB                *db.DB
	CollectionService *service.CollectionService
	DocumentService   *service.DocumentService
	Config            *config.Config
}

// NewApp initializes the application
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize database
	dbConfig := db.NewConfig(cfg.SqlitePath)
	database, err := db.New(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	collectionRepo := db.NewCollectionRepository(database)
	documentRepo := db.NewDocumentRepository(database, collectionRepo)

	// Initialize services
	collectionService := service.NewCollectionService(collectionRepo)
	documentService := service.NewDocumentService(documentRepo)

	// Initialize router
	router := mux.NewRouter()

	app := &App{
		Router:            router,
		DB:                database,
		CollectionService: collectionService,
		DocumentService:   documentService,
		Config:            cfg,
	}

	app.setupRoutes()
	return app, nil
}

// setupRoutes configures the routes
func (a *App) setupRoutes() {
	// Register middleware
	a.Router.Use(middleware.LoggingMiddleware)
	a.Router.Use(middleware.RecoveryMiddleware)

	// Initialize handlers
	healthHandler := handlers.HealthHandler(a.Config)
	collectionHandlers := handlers.NewCollectionHandlers(a.CollectionService)
	documentHandlers := handlers.NewDocumentHandlers(a.DocumentService)
	protectedHandler := handlers.ProtectedHandler(a.Config)

	// Register health endpoint
	a.Router.HandleFunc("/health", healthHandler).Methods("GET")

	// Collection routes
	a.Router.HandleFunc("/api/collections", collectionHandlers.ListCollections()).Methods("GET")
	a.Router.HandleFunc("/api/collections", collectionHandlers.CreateCollection()).Methods("POST")
	a.Router.HandleFunc("/api/collections/{name}", collectionHandlers.GetCollection()).Methods("GET")
	a.Router.HandleFunc("/api/collections/{name}", collectionHandlers.DeleteCollection()).Methods("DELETE")

	// Document routes
	a.Router.HandleFunc("/api/collections/{name}/documents", documentHandlers.ListDocuments()).Methods("GET")
	a.Router.HandleFunc("/api/collections/{name}/documents", documentHandlers.CreateDocument()).Methods("POST")
	a.Router.HandleFunc("/api/collections/{name}/documents/{id}", documentHandlers.GetDocument()).Methods("GET")
	a.Router.HandleFunc("/api/collections/{name}/documents/{id}", documentHandlers.UpdateDocument()).Methods("PUT")
	a.Router.HandleFunc("/api/collections/{name}/documents/{id}", documentHandlers.DeleteDocument()).Methods("DELETE")

	// Bulk operations
	a.Router.HandleFunc("/api/collections/{name}/bulk", documentHandlers.BulkCreateDocuments()).Methods("POST")
	a.Router.HandleFunc("/api/upload/{name}", documentHandlers.UploadJSONFile()).Methods("POST")

	// Create a subrouter for protected routes
	protectedRouter := a.Router.PathPrefix("/api/protected").Subrouter()

	// Apply authentication middleware if enabled
	if a.Config.EnableBasicAuth {
		authMiddleware := middleware.BasicAuthMiddleware(a.Config.AuthUsername, a.Config.AuthPassword)
		protectedRouter.Use(authMiddleware)
	}

	// Add protected routes
	protectedRouter.HandleFunc("/info", protectedHandler).Methods("GET")

}

// SetupRouter configures and returns the router with all routes and middleware for testing
func SetupRouter(cfg *config.Config) *mux.Router {
	app, err := NewApp(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize application: %v", err))
	}
	return app.Router
}

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rbehzadan/schemaless-api/internal/api"
	"github.com/rbehzadan/schemaless-api/internal/api/handlers"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

func TestHealthHandler(t *testing.T) {
	// Set up test configuration
	cfg := config.NewConfig()
	cfg.Version = "1.0.0-test"
	cfg.StartTime = time.Now().Add(-1 * time.Minute) // Set start time to 1 minute ago

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthHandler(cfg))

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}

	// Unmarshal the response body
	var response api.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Check the response structure
	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v",
			response.Status, "success")
	}

	// Check the health info
	healthInfo, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Data is not of expected type: %T", response.Data)
		return
	}

	// Check status
	if healthInfo["status"] != "ok" {
		t.Errorf("health status is wrong: got %v want %v",
			healthInfo["status"], "ok")
	}

	// Check version
	if healthInfo["version"] != "1.0.0-test" {
		t.Errorf("version is wrong: got %v want %v",
			healthInfo["version"], "1.0.0-test")
	}

	// Check uptime (should be around 1 minute, but we'll just check if it exists)
	if uptime, exists := healthInfo["uptime"]; !exists || uptime == "" {
		t.Errorf("uptime is missing or empty: %v", uptime)
	}
}

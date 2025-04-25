package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbehzadan/schemaless-api/internal/api"
	"github.com/rbehzadan/schemaless-api/internal/api/middleware"
)

// panicHandler returns a handler that panics for testing recovery middleware
func panicHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create handler with middleware
	handler := middleware.RecoveryMiddleware(panicHandler())

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
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
	if response.Status != "error" {
		t.Errorf("handler returned wrong status: got %v want %v",
			response.Status, "error")
	}

	// Check the error info
	if response.Error == nil {
		t.Errorf("Error info is nil")
		return
	}

	if response.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("wrong error code: got %v want %v",
			response.Error.Code, "INTERNAL_SERVER_ERROR")
	}
}

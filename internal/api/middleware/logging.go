package middleware

import (
	"log"
	"net/http"
	"time"
)

// ResponseWriterWrapper is a wrapper for http.ResponseWriter that captures status code and response size
type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

// NewResponseWriterWrapper creates a new ResponseWriterWrapper
func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

// WriteHeader captures the status code
func (rww *ResponseWriterWrapper) WriteHeader(code int) {
	rww.StatusCode = code
	rww.ResponseWriter.WriteHeader(code)
}

// Write captures the response size
func (rww *ResponseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rww.ResponseWriter.Write(b)
	rww.Size += size
	return size, err
}

// LoggingMiddleware logs HTTP requests with standard details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status and size
		rww := NewResponseWriterWrapper(w)

		// Process request
		next.ServeHTTP(rww, r)

		// Calculate duration
		duration := time.Since(start)

		// Log request details in industry standard format
		log.Printf(
			"%s | %s | %s | %d | %d bytes | %s",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			rww.StatusCode,
			rww.Size,
			duration,
		)
	})
}

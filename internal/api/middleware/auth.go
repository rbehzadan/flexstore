package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/rbehzadan/schemaless-api/internal/api"
)

// BasicAuthMiddleware provides HTTP Basic Authentication
func BasicAuthMiddleware(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorized(w)
				return
			}

			// Check if it's Basic auth
			if !strings.HasPrefix(authHeader, "Basic ") {
				unauthorized(w)
				return
			}

			// Decode credentials
			payload, err := base64.StdEncoding.DecodeString(authHeader[6:])
			if err != nil {
				unauthorized(w)
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				unauthorized(w)
				return
			}

			providedUser := pair[0]
			providedPass := pair[1]

			// Check credentials using constant-time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(providedUser), []byte(username)) != 1 ||
				subtle.ConstantTimeCompare([]byte(providedPass), []byte(password)) != 1 {
				unauthorized(w)
				return
			}

			// Authentication successful, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// unauthorized sends a 401 Unauthorized response
func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	api.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
}

package handlers

import (
	"net/http"
	"time"

	"github.com/rbehzadan/schemaless-api/internal/api"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

// ProtectedInfoResponse contains protected information
type ProtectedInfoResponse struct {
	Message    string    `json:"message"`
	ServerTime time.Time `json:"server_time"`
	Version    string    `json:"version"`
	Uptime     string    `json:"uptime"`
}

// ProtectedHandler returns a handler function for the protected endpoint
func ProtectedHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := ProtectedInfoResponse{
			Message:    "This is a protected endpoint that requires authentication",
			ServerTime: time.Now(),
			Version:    cfg.Version,
			Uptime:     cfg.GetUptime(),
		}

		resp := api.Response{
			Status: "success",
			Data:   info,
		}

		api.RespondWithJSON(w, http.StatusOK, resp)
	}
}

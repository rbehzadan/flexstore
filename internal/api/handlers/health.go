package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rbehzadan/schemaless-api/internal/api"
	"github.com/rbehzadan/schemaless-api/pkg/config"
)

// HealthHandler returns a handler function for health checks
func HealthHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		healthInfo := api.HealthResponse{
			Status:  "ok",
			Version: cfg.Version,
			Uptime:  cfg.GetUptime(),
		}

		resp := api.Response{
			Status: "success",
			Data:   healthInfo,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

package api

import (
	"encoding/json"
	"net/http"
)

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, statusCode int, errorCode, message string) {
	response := Response{
		Status: "error",
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: message,
		},
	}
	RespondWithJSON(w, statusCode, response)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// If data is already a Response, use it directly
	var response Response
	if resp, ok := data.(Response); ok {
		response = resp
	} else {
		response = Response{
			Status: "success",
			Data:   data,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

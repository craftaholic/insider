package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/craftaholic/insider/internal/domain/dto"
	"github.com/craftaholic/insider/internal/shared/log"
)

// sendJSONResponse for response handling.
func (mc *MessageController) sendJSONResponse(c context.Context, w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	logger := log.FromCtx(c)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed writing response", "error", err)
		// Try to send a simple error response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (mc *MessageController) sendErrorResponse(
	c context.Context,
	w http.ResponseWriter,
	message string,
	statusCode int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	logger := log.FromCtx(c)

	errorResp := dto.CreateErrorResponse(message)
	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		logger.Error("Failed writing error response", "error", err)
		// Fallback to standard http.Error
		http.Error(w, message, statusCode)
	}

	logger.Error("Request handled failed", "error", message)
}

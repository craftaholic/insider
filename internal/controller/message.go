package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/craftaholic/insider/internal/utils"
)

type MessageController struct {
	MessageUsecase interfaces.MessageUsecase
}

func NewMessageController(messageUsecase interfaces.MessageUsecase) *MessageController {
	return &MessageController{
		MessageUsecase: messageUsecase,
	}
}

// Start handles starting the automated sending notification
// swagger:route GET /message/start message start
//
// # Start Automated Message Sending
//
// This endpoint starts the automated message sending process.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: startResponse
//	500: errorResponse
func (mc *MessageController) Start(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Starting automated sending message")

	// Use background context because the logic will run in another thread
	// in the background.
	err := mc.MessageUsecase.StartAutomatedSending(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Request handled failed", "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return proper JSON response
	response := map[string]string{"status": "OK", "message": "Automated sending started successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}
	logger.Info("Finished start automated sending message request")
}

// Stop handles stopping the automated sending notification
// swagger:route POST /message/stop message stop
//
// # Stop Automated Message Sending
//
// This endpoint stops the automated message sending process.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: stopResponse
//	500: errorResponse
func (mc *MessageController) Stop(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Stopping automated sending message")
	ctx := logger.WithCtx(r.Context())

	err := mc.MessageUsecase.StopAutomatedSending(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Request handled failed", "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return proper JSON response
	response := map[string]string{"status": "OK", "message": "Automated sending stopped successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}
	logger.Info("Finished stop automated sending message request")
}

// GetSentMessagesWithPagination retrieves sent messages with pagination
// swagger:route GET /message/sent message getSentMessages
//
// # Get Sent Messages with Pagination
//
// Retrieves a paginated list of sent messages.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: messagesResponse
//	400: errorResponse
//	500: errorResponse
func (mc *MessageController) GetSentMessagesWithPagination(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Getting sent messages with pagination")
	ctx := logger.WithCtx(r.Context())

	page := r.URL.Query().Get("page")
	// Default value if page not declared
	pageInt := 1

	if page != "" {
		var err error
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			logger.Error("Request handled failed", "error", err.Error())
			return
		}

		// Validate page number is positive
		if pageInt < 1 {
			http.Error(w, "Page number must be greater than 0", http.StatusBadRequest)
			logger.Error("Request handled failed", "error", "invalid page number")
			return
		}
	}

	messages, err := mc.MessageUsecase.GetSentMessagesWithPagination(ctx, pageInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Request handled failed", "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		logger.Error("Failed to marshal response", "error", err.Error())
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}
	logger.Info("Finished getting sent messages with pagination request")
}

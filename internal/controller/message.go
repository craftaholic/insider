package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/craftaholic/insider/internal/domain/dto"
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
// swagger:route POST /service/start message start
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
		mc.sendErrorResponse(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateStandardResponse("OK", "Automated sending started successfully")
	mc.sendJSONResponse(r.Context(), w, response, http.StatusOK)
	logger.Info("Finished start automated sending message request")
}

// Stop handles stopping the automated sending notification
// swagger:route POST /service/stop message stop
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
		mc.sendErrorResponse(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateStandardResponse("OK", "Automated sending stopped successfully")
	mc.sendJSONResponse(r.Context(), w, response, http.StatusOK)
	logger.Info("Finished stop automated sending message request")
}

// Status method get the status of automated sending notification
// swagger:route GET /service/status message stop
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
//	200: statusResponse
//	500: errorResponse
func (mc *MessageController) Status(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Getting automated sending message service status")
	ctx := logger.WithCtx(r.Context())

	status, err := mc.MessageUsecase.GetAutomatedSendingStatus(ctx)
	if err != nil {
		mc.sendErrorResponse(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Automated sending service is running"
	if !status {
		message = "Automated sending service is stopped"
	}

	response := dto.CreateStandardResponse("OK", message)
	mc.sendJSONResponse(r.Context(), w, response, http.StatusOK)
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
			mc.sendErrorResponse(ctx, w, "Invalid page number", http.StatusBadRequest)
			return
		}

		// Validate page number is positive
		if pageInt < 1 {
			mc.sendErrorResponse(ctx, w, "Page number must be greater than 0", http.StatusBadRequest)
			return
		}
	}

	// Get domain entities from usecase
	messages, err := mc.MessageUsecase.GetSentMessagesWithPagination(ctx, pageInt)
	if err != nil {
		mc.sendErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert domain entities to DTOs
	messageDTOs := dto.ConvertMessagesToDTO(messages)

	mc.sendJSONResponse(ctx, w, messageDTOs, http.StatusOK)
	logger.Info("Finished getting sent messages with pagination request")
}

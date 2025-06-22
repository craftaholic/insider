package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/craftaholic/insider/internal/utils"
)

type MessageController struct {
	MessageUsecase domain.MessageUsecase
}

func NewMessageController(messageUsecase domain.MessageUsecase) *MessageController {
	return &MessageController{
		MessageUsecase: messageUsecase,
	}
}

func (mc *MessageController) Start(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Starting automated sending message")

	err := mc.MessageUsecase.StartAutomatedSending(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Request handled failed", "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}

	logger.Info("Finished start automated sending message request")
}

func (mc *MessageController) Stop(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Starting automated sending message")
	ctx := logger.WithCtx(r.Context())

	err := mc.MessageUsecase.StopAutomatedSending(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Request handled failed", "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}

	logger.Info("Finished start automated sending message request")
}

func (mc *MessageController) GetSentMessagesWithPagination(w http.ResponseWriter, r *http.Request) {
	logger := log.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(mc))
	logger.Info("Starting automated sending message")
	ctx := logger.WithCtx(r.Context())

	page := r.URL.Query().Get("page")
	// Default value
	pageInt := 0
	if page != nil {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			logger.Error("Request handled failed", "error", err.Error())
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

	// TODO: handle this return
	_, err = w.Write([]byte(messages[0].Content))
	if err != nil {
		logger.Error("Failed writing response", "error", err)
	}

	logger.Info("Finished start automated sending message request")
}

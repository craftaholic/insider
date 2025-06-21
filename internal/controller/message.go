package controller

import (
	"net/http"

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
	ctx := logger.WithCtx(r.Context())

	err := mc.MessageUsecase.StartAutomatedSending(ctx)
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

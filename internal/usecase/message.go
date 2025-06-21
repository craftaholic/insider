package usecase

import (
	"context"

	"github.com/craftaholic/insider/internal/domain"
)

type MessageUsecase struct {
	messageRepository   domain.MessageRepository
	cacheRepository     domain.CacheRepository
	notificationService domain.NotificationService
}

func NewMessageUsecase(messageRepository domain.MessageRepository, cacheRepository domain.CacheRepository, notificationService domain.NotificationService) domain.MessageUsecase {
	return &MessageUsecase{
		messageRepository:   messageRepository,
		cacheRepository:     cacheRepository,
		notificationService: notificationService,
	}
}

func (mu *MessageUsecase) StartAutomatedSending(c context.Context) error {
	return nil
}

func (mu *MessageUsecase) StopAutomatedSending(c context.Context) error {
	return nil
}

func (mu *MessageUsecase) GetMessagesWithPagination(c context.Context, page int) ([]domain.Message, error) {
	return []domain.Message{}, nil
}

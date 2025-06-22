package repository

import (
	"context"

	"github.com/craftaholic/insider/internal/domain"
)

type NotificationService struct {
	apiKey   string
	endPoint string
}

func NewNotificationService(apiKey string, endPoint string) domain.NotificationService {
	return &NotificationService{
		apiKey:   apiKey,
		endPoint: endPoint,
	}
}

func (ns *NotificationService) SendNotification(c context.Context, message domain.Message) (string, error) {
	return "", nil
}

package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/craftaholic/insider/internal/utils"
	"github.com/go-resty/resty/v2"
)

// NotificationRequest structs is for parsing request.
type NotificationRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// NotificationResponse structs is for parsing response.
type NotificationResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type NotificationService struct {
	apiKey   string
	endPoint string
	client   *resty.Client
}

func NewNotificationService(client *resty.Client, apiKey string, endPoint string) domain.NotificationService {
	return &NotificationService{
		apiKey:   apiKey,
		endPoint: endPoint,
		client:   client,
	}
}

func (ns *NotificationService) SendNotification(c context.Context, message domain.Message) (string, error) {
	logger := log.FromCtx(c).WithFields("action", "Sending notification")

	notificationRequest := NotificationRequest{
		To:      message.PhoneNumber,
		Content: message.Content,
	}

	body, err := json.Marshal(notificationRequest)
	if err != nil {
		logger.Error("Can't marshal notification request object to byte")
		return "", err
	}

	response, webhookErr := ns.client.R().
		SetHeader("Authorization", "Bearer "+ns.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(ns.endPoint)
	if webhookErr != nil {
		logger.Error("Error sending notification", "error", err)
		return "", webhookErr
	}

	var notificationResponse NotificationResponse
	err = json.Unmarshal(response.Body(), &notificationResponse)
	if err != nil {
		logger.Error("Error sending notification", "error", err)
	}

	if notificationResponse.Message != "Accepted" {
		logger.Error("Error sending notification", "response for webhook", notificationResponse.Message)
		return "", errors.New("Error sending notification " + "response for webhook " + notificationResponse.Message)
	}

	return notificationResponse.MessageID, nil
}

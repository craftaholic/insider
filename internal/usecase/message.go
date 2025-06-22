package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/shared/log"
)

type MessageUsecase struct {
	messageRepository   domain.MessageRepository
	cacheRepository     domain.CacheRepository
	notificationService domain.NotificationService

	workerPool *WorkerPool
	cancel     context.CancelFunc
	isRunning  bool
	mu         sync.RWMutex
}

func NewMessageUsecase(
	messageRepository domain.MessageRepository,
	cacheRepository domain.CacheRepository,
	notificationService domain.NotificationService,
) domain.MessageUsecase {
	return &MessageUsecase{
		messageRepository:   messageRepository,
		cacheRepository:     cacheRepository,
		notificationService: notificationService,
	}
}

func (mu *MessageUsecase) StartAutomatedSending(c context.Context) error {
	mu.mu.Lock()
	defer mu.mu.Unlock()

	if mu.isRunning {
		return errors.New("automated sending is already running")
	}

	serviceCtx, cancel := context.WithCancel(c)
	mu.cancel = cancel

	// Create worker pool
	mu.workerPool = newWorkerPool(c, 5) // 5 concurrent workers
	mu.workerPool.Start(mu.processSingleMessage)

	// Start message fetcher
	go mu.messageFetcher(serviceCtx)

	mu.isRunning = true
	return nil
}

func (mu *MessageUsecase) messageFetcher(c context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Done():
			return
		case <-ticker.C:
			if mu.isRunning {
				messages, err := mu.messageRepository.GetPending(c, 2)
				if err != nil {
					continue
				}

				for _, message := range messages {
					log.FromCtx(c).Info("Fetching", "message", message.ID)
					succeed := mu.workerPool.AddJob(message)
					// If can't add job to queue -> convert the status
					if !succeed {
						message.Status = "pending"
						err = mu.messageRepository.Update(c, message.ID, message)
						if err != nil {
							log.FromCtx(c).Error("Error changing status of message back to pending", "message", message.ID)
						}
					}
				}
			}
		}
	}
}

func (mu *MessageUsecase) StopAutomatedSending(c context.Context) error {
	mu.mu.Lock()
	defer mu.mu.Unlock()

	if !mu.isRunning {
		return errors.New("automated sending is not running")
	}

	if mu.workerPool != nil {
		mu.workerPool.Stop()
	}

	mu.isRunning = false
	mu.cancel()
	return nil
}

func (mu *MessageUsecase) GetSentMessagesWithPagination(c context.Context, page int) ([]domain.Message, error) {
	return nil, nil
}

func (mu *MessageUsecase) processSingleMessage(c context.Context, message domain.Message) error {
	logger := log.FromCtx(c).WithFields("message", message.ID)

	logger.Info("Processing message")

	logger.Info("Sending notification")
	messageUUID, err := mu.notificationService.SendNotification(c, message)
	if err != nil {
		logger.Error("Failed sending notification", "error", err)
		return err
	}

	if messageUUID == "" {
		logger.Error("MessageUUID return from send notification is empty")
		return errors.New("the message uuid from send notification webhook is empty")
	}

	timestamp := time.Now()
	updates := map[string]interface{}{
		"status":       "sent",
		"sent_at":      timestamp,
		"message_uuid": messageUUID, // Store the UUID from notification service
		"updated_at":   timestamp,
	}
	err = mu.messageRepository.UpdateSelective(c, message.ID, updates)
	if err != nil {
		logger.Error("Failed updating message in the db", "error", err)
		return err
	}

	timestampBytes, err := timestamp.MarshalBinary()
	if err != nil {
		logger.Error("Failed converting timestamp to binary to store into cache", "error", err)
		return err
	}

	err = mu.cacheRepository.Set(messageUUID, timestampBytes, 0)
	if err != nil {
		logger.Error("Failed storing messageID into cache", "error", err)
		return err
	}

	logger.Info("Message handled successfully")

	return nil
}

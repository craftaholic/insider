package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/craftaholic/insider/internal/domain/entity"
	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/craftaholic/insider/internal/shared/log"
)

type MessageUsecase struct {
	messageRepository   interfaces.MessageRepository
	cacheRepository     interfaces.CacheRepository
	notificationService interfaces.NotificationService

	jobBuffer            int
	workerCount          int
	producerCronDuration int
	producerBatchNumber  int

	workerPool *WorkerPool
	cancel     context.CancelFunc
	isRunning  bool
	mu         sync.RWMutex
}

func NewMessageUsecase(
	messageRepository interfaces.MessageRepository,
	cacheRepository interfaces.CacheRepository,
	notificationService interfaces.NotificationService,
	jobBuffer int,
	workerCount int,
	producerCronDuration int,
	producerBatchNumber int,
) interfaces.MessageUsecase {
	return &MessageUsecase{
		messageRepository:    messageRepository,
		cacheRepository:      cacheRepository,
		notificationService:  notificationService,
		workerCount:          workerCount,
		jobBuffer:            jobBuffer,
		producerCronDuration: producerCronDuration,
		producerBatchNumber:  producerBatchNumber,
	}
}

func (mu *MessageUsecase) StartAutomatedSending(c context.Context) error {
	mu.mu.Lock()
	defer mu.mu.Unlock()

	logger := log.FromCtx(c)

	if mu.isRunning {
		logger.Info("Automated sending already started don't need to do anything")	
		return nil
	}

	serviceCtx, cancel := context.WithCancel(c)
	mu.cancel = cancel

	// Create worker pool
	mu.workerPool = newWorkerPool(c, mu.workerCount, mu.jobBuffer) // 5 concurrent workers
	mu.workerPool.Start(mu.processSingleMessage)

	// Start message fetcher
	go mu.messageFetcher(serviceCtx)

	mu.isRunning = true
	return nil
}

func (mu *MessageUsecase) messageFetcher(c context.Context) {
	ticker := time.NewTicker(time.Duration(mu.producerCronDuration) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Done():
			return
		case <-ticker.C:
			if mu.isRunning {
				messages, err := mu.messageRepository.GetPending(c, mu.producerBatchNumber)
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
							log.FromCtx(c).
								Error("Error changing status of message back to pending", "message", message.ID)
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

	logger := log.FromCtx(c)
	logger.Info("Stopping automated sending notification...")

	if !mu.isRunning {
		logger.Info("Automated sending service already stopped")
	}

	if mu.workerPool != nil {
		mu.workerPool.Stop()
	}

	mu.isRunning = false
	mu.cancel()
	logger.Info("Stopping automated sending notification successfully")
	return nil
}

func (mu *MessageUsecase) GetSentMessagesWithPagination(c context.Context, page int) ([]entity.Message, error) {
	logger := log.FromCtx(c).WithFields("action", "Get sent message with pagination", "page", page)
	logger.Info("Getting all sent message of this page")

	return mu.messageRepository.GetSentWithPagination(c, page)
}

// This function will provide at-least 1 notification sent but it will make sure
// there are no cases where notification never sent.
func (mu *MessageUsecase) processSingleMessage(ctx context.Context, message entity.Message) error {
	logger := log.FromCtx(ctx).WithFields("message_id", message.ID)
	logger.Info("Processing message")

	// 1. Send notification
	logger.Info("Sending notification")
	messageUUID, err := mu.notificationService.SendNotification(ctx, message)
	if err != nil {
		// Update status to failed before returning
		mu.handleMessageFailure(ctx, message.ID, "notification_failed", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if messageUUID == "" {
		logger.Error("Notification sent but returned empty message UUID")
	}

	// 2. Update message status (only update what changed)
	timestamp := time.Now()
	updates := map[string]any{
		"status":     "sent",
		"sent_at":    timestamp,
		"message_id": messageUUID, // Store the UUID from notification service
		"updated_at": timestamp,
	}

	err = mu.messageRepository.UpdateSelective(ctx, message.ID, updates)
	if err != nil {
		// This error won't return cause message already sent
		logger.Error("Failed to update message status", "error", err)
	}

	// 3. Cache the result
	if err = mu.cacheMessageResult(messageUUID, timestamp); err != nil {
		// This error won't return cause message already sent
		logger.Warn("Failed to cache message result", "error", err)
	}

	logger.Info("Message processed successfully", "message_uuid", messageUUID)
	return nil
}

func (mu *MessageUsecase) handleMessageFailure(
	ctx context.Context,
	messageID uint64,
	reason string,
	originalErr error,
) {
	logger := log.FromCtx(ctx).WithFields("message_id", messageID)

	updates := map[string]any{
		"status":        "failed",
		"error_message": fmt.Sprintf("%s: %v", reason, originalErr),
		"updated_at":    time.Now(),
	}

	if err := mu.messageRepository.UpdateSelective(ctx, messageID, updates); err != nil {
		logger.Error("Failed to set message status to failed", "error", err)
	}
}

// Helper function for caching.
func (mu *MessageUsecase) cacheMessageResult(messageUUID string, timestamp time.Time) error {
	timestampBytes, err := timestamp.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal timestamp: %w", err)
	}

	// Use a reasonable TTL instead of 0 (forever)
	ttl := 0 * time.Hour // Cache for 24 hours
	return mu.cacheRepository.Set(messageUUID, timestampBytes, ttl)
}

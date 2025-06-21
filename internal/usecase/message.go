package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/craftaholic/insider/internal/domain"
)

type MessageUsecase struct {
	messageRepository   domain.MessageRepository
	cacheRepository     domain.CacheRepository
	notificationService domain.NotificationService

	workerPool *WorkerPool
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

		mu: sync.RWMutex{},
	}
}

func (mu *MessageUsecase) StartAutomatedSending(ctx context.Context) error {
	mu.mu.Lock()
	defer mu.mu.Unlock()

	if mu.isRunning {
		return errors.New("automated sending is already running")
	}

	// Create worker pool
	mu.workerPool = newWorkerPool(context.Background(), 5) // 5 concurrent workers
	mu.workerPool.Start(mu.processSingleMessage)

	// Start message fetcher
	go mu.messageFetcher(ctx)

	mu.isRunning = true
	return nil
}

func (mu *MessageUsecase) messageFetcher(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			messages, err := mu.messageRepository.GetPending(ctx, 2)
			if err != nil {
				continue
			}

			for _, message := range messages {
				mu.workerPool.AddJob(message)
			}
		}
	}
}

func (mu *MessageUsecase) StopAutomatedSending(ctx context.Context) error {
	mu.mu.Lock()
	defer mu.mu.Unlock()

	if !mu.isRunning {
		return errors.New("automated sending is not running")
	}

	if mu.workerPool != nil {
		mu.workerPool.Stop()
	}

	mu.isRunning = false
	return nil
}

func (mu *MessageUsecase) GetSentMessagesWithPagination(c context.Context, page int) ([]Message, error) {
	return nil, nil
}

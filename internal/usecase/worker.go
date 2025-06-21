package usecase

import (
	"context"

	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/shared/log"
)

type WorkerPool struct {
	ctx         context.Context
	cancel      context.CancelFunc
	jobChan     chan domain.Message
	workerCount int
}

func newWorkerPool(ctx context.Context, workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		ctx:         ctx,
		cancel:      cancel,
		jobChan:     make(chan domain.Message, 100),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) Start(processor func(context.Context, domain.Message) error) {
	for i := 0; i < wp.workerCount; i++ {
		go func(workerID int) {
			for {
				select {
				case <-wp.ctx.Done():
					return
				case message := <-wp.jobChan:
					if err := processor(wp.ctx, message); err != nil {
						// Log error
						log.FromCtx(wp.ctx).Error("Worker failed to process message",
							"workerID", workerID, "messageID", message.ID, "error", err)
					}
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) AddJob(message domain.Message) {
	select {
	case wp.jobChan <- message:
	case <-wp.ctx.Done():
		return
	}
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobChan)
}

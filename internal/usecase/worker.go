package usecase

import (
	"context"
	"sync"

	"github.com/craftaholic/insider/internal/domain"
	"github.com/craftaholic/insider/internal/shared/log"
)

type WorkerPool struct {
	ctx         context.Context
	cancel      context.CancelFunc
	jobChan     chan domain.Message
	workerCount int
	wg          sync.WaitGroup
}

func newWorkerPool(ctx context.Context, workerCount int, buffer int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		ctx:         ctx,
		cancel:      cancel,
		jobChan:     make(chan domain.Message, buffer),
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
	}
}

// Start will create multiple workers each runs in
// 1 go routines. Each of these workers will handle
// 1 message from the db (every 2 mins there will
// be new messages sent into the channel).
func (wp *WorkerPool) Start(processor func(context.Context, domain.Message) error) {
	for i := range wp.workerCount {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()

			for {
				select {
				case message := <-wp.jobChan:
					// This will always be handled first if there are still
					// messages in the channel this will execute all of it first
					// before checking the condition of the context
					if err := processor(wp.ctx, message); err != nil {
						// Log error
						log.FromCtx(wp.ctx).Error("Worker failed to process message",
							"workerID", workerID, "messageID", message.ID, "error", err)
					}
				case <-wp.ctx.Done():
					// If all messages in the channel is handled then it will check
					// the ctx.Done condition to make sure no messages droped while
					// there is a stop signal
					return
				}
			}
		}(i)
	}
}

// AddJob will continue add job to the jobChan buffer
// if there is a cancel signal event -> stop receiving
// new message.
func (wp *WorkerPool) AddJob(message domain.Message) bool {
	select {
	// Always check the context first to
	case <-wp.ctx.Done():
		return false
	case wp.jobChan <- message:
		return true
	default:
		return false
	}
}

// Stop function will send a signal event to the context
// that's being used to stop receiving all new messages.
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait()
}

// Close function close the channel, this is to spit it
// away from the Stop function so it will avoid having a
// potential race condition.
func (wp *WorkerPool) Close() {
	wp.wg.Wait()
	close(wp.jobChan)
}

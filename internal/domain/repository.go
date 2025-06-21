package domain

import (
	"context"
	"time"
)

type MessageRepository interface {
	Update(c context.Context, id uint64, message Message) error
	GetPending(c context.Context, batch int) ([]Message, error)
	GetSentWithPagination(c context.Context, page int) ([]Message, error)
}

type CacheRepository interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}

type NotificationService interface {
	SendNotification(c context.Context, message Message) error
}

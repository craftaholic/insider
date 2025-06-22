package interfaces

import (
	"context"
	"time"

	"github.com/craftaholic/insider/internal/domain/entity"
)

type MessageRepository interface {
	Update(c context.Context, id uint64, message entity.Message) error
	UpdateSelective(ctx context.Context, id uint64, updates map[string]any) error
	GetPending(c context.Context, batch int) ([]entity.Message, error)
	GetSentWithPagination(c context.Context, page int) ([]entity.Message, error)
}

type CacheRepository interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}

type NotificationService interface {
	SendNotification(c context.Context, message entity.Message) (string, error)
}

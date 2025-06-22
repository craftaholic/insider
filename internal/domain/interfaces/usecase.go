package interfaces

import (
	"context"

	"github.com/craftaholic/insider/internal/domain/entity"
)

type MessageUsecase interface {
	StartAutomatedSending(c context.Context) error
	StopAutomatedSending(c context.Context) error
	GetAutomatedSendingStatus(c context.Context) (bool, error)
	GetSentMessagesWithPagination(c context.Context, page int) ([]entity.Message, error)
}

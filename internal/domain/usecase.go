package domain

import (
	"context"
)

type MessageUsecase interface {
	StartAutomatedSending(c context.Context) error
	StopAutomatedSending(c context.Context) error
	GetSentMessagesWithPagination(c context.Context, page int) ([]Message, error)
}

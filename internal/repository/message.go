package repository

import (
	"context"
	"errors"

	"github.com/craftaholic/insider/internal/domain"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) Send(ctx context.Context, message domain.Message) error {
	// TODO: update this to interact with infra layer for sending instead
	// then update to db
	if message.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if message.Content == "" {
		return errors.New("message content is required")
	}

	// Set default status if not provided
	if message.Status == "" {
		message.Status = domain.StatusPending
	}

	result := r.db.WithContext(ctx).Create(&message)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *messageRepository) GetPending(ctx context.Context, batch int) ([]domain.Message, error) {
	if batch <= 0 {
		return nil, errors.New("batch size must be greater than 0")
	}

	var messages []domain.Message

	err := r.db.WithContext(ctx).
		Raw("SELECT * FROM get_unsent_messages(?)", batch).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetSentWithPagination(ctx context.Context, page int) ([]domain.Message, error) {
	if page <= 0 {
		return nil, errors.New("page must be greater than 0")
	}

	const pageSize = 20
	offset := (page - 1) * pageSize

	var messages []domain.Message

	err := r.db.WithContext(ctx).
		Where("status = ?", domain.StatusSent).
		Offset(offset).
		Limit(pageSize).
		Order("sent_at DESC").
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/craftaholic/insider/internal/domain/entity"
	"github.com/craftaholic/insider/internal/domain/interfaces"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) interfaces.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) UpdateSelective(ctx context.Context, id uint64, updates map[string]any) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update message with id %d: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message with id %d not found", id)
	}

	return nil
}

func (r *messageRepository) Update(ctx context.Context, id uint64, message entity.Message) error {
	// Update the message with the given ID
	result := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("id = ?", id).
		Updates(message)

	// Check for database errors
	if result.Error != nil {
		return fmt.Errorf("failed to update message with id %d: %w", id, result.Error)
	}

	// Check if any rows were affected (message exists)
	if result.RowsAffected == 0 {
		return fmt.Errorf("message with id %d not found", id)
	}

	return nil
}

func (r *messageRepository) GetPending(ctx context.Context, batch int) ([]entity.Message, error) {
	if batch <= 0 {
		return nil, errors.New("batch size must be greater than 0")
	}

	var messages []entity.Message

	err := r.db.WithContext(ctx).
		Raw("SELECT * FROM get_unsent_messages(?)", batch).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetSentWithPagination(ctx context.Context, page int) ([]entity.Message, error) {
	if page <= 0 {
		return nil, errors.New("page must be greater than 0")
	}

	const pageSize = 20
	offset := (page - 1) * pageSize

	var messages []entity.Message

	err := r.db.WithContext(ctx).
		Where("status = ?", entity.StatusSent).
		Offset(offset).
		Limit(pageSize).
		Order("sent_at DESC").
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

package dto

import (
	"time"

	"github.com/craftaholic/insider/internal/domain/entity"
)

// MessageDTO represents a message for API responses
// swagger:model
type MessageDTO struct {
	// Message ID
	// example: 123
	ID uint64 `json:"id"`

	// Phone Number
	// example: +84338252331
	PhoneNumber string `json:"phone_number"`

	// Message content
	// example: Hello, this is a test message
	Content string `json:"content"`

	// Timestamp when message was created (ISO 8601 string)
	// example: 2025-06-22T10:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when message was sent (ISO 8601 string, nullable)
	// example: 2025-06-22T10:35:00Z
	SentAt *time.Time `json:"sent_at"`

	// Message status
	// example: sent
	Status entity.MessageStatus `json:"status"`

	// Message ID of the notification sent
	// example: e975f171-3ce5-4ea4-bf03-ae5b8849d2cb
	MessageID *string `json:"message_id"`

	// Error message if sending failed
	// example: null
	ErrorMessage *string `json:"error_message,omitempty"`

	// Updated At
	// example: 2025-06-22T10:35:00Z
	UpdatedAt *time.Time `json:"updated_at"`
}

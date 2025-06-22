package entity

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// MessageStatus represents the status enum.
type MessageStatus string

const (
	StatusPending    MessageStatus = "pending"
	StatusProcessing MessageStatus = "processing"
	StatusSent       MessageStatus = "sent"
	StatusFailed     MessageStatus = "failed"
)

// Scan implements the Scanner interface for database reads.
func (ms *MessageStatus) Scan(value any) error {
	if value == nil {
		*ms = StatusPending
		return nil
	}
	if str, ok := value.(string); ok {
		*ms = MessageStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into MessageStatus", value)
}

// Value implements the driver.Valuer interface for database writes.
func (ms *MessageStatus) Value() (driver.Value, error) {
	return ms, nil
}

type Message struct {
	ID           uint64        `json:"id"            gorm:"primaryKey;column:id"`
	PhoneNumber  string        `json:"phone_number"  gorm:"column:phone_number;type:varchar(20);not null"`
	Content      string        `json:"content"       gorm:"column:content;type:text;not null"`
	Status       MessageStatus `json:"status"        gorm:"column:status;type:varchar(20);default:pending;check:status IN ('pending', 'processing', 'sent', 'failed')"`
	CreatedAt    time.Time     `json:"created_at"    gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP"`
	SentAt       *time.Time    `json:"sent_at"       gorm:"column:sent_at;type:timestamptz"`
	MessageID    *string       `json:"message_id"    gorm:"column:message_id;type:varchar(255)"`
	ErrorMessage *string       `json:"error_message" gorm:"column:error_message;type:text"`
	UpdatedAt    *time.Time    `json:"updated_at"    gorm:"column:updated_at;type:timestamptz"`
}

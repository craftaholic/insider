package dto

// MessageDTO represents a message for API responses
// swagger:model
type MessageDTO struct {
	// Message ID
	// example: 123
	ID int `json:"id"`

	// Message content
	// example: Hello, this is a test message
	Content string `json:"content"`

	// Recipient information
	// example: user@example.com
	Recipient string `json:"recipient"`

	// Sender information
	// example: sender@example.com
	Sender string `json:"sender"`

	// Subject of the message
	// example: Important notification
	Subject string `json:"subject"`

	// Message type
	// example: email
	Type string `json:"type"`

	// Timestamp when message was created
	// example: 2025-06-22T10:30:00Z
	CreatedAt string `json:"created_at"`

	// Timestamp when message was sent
	// example: 2025-06-22T10:35:00Z
	SentAt string `json:"sent_at"`

	// Message status
	// example: sent
	Status string `json:"status"`

	// Number of retry attempts
	// example: 0
	RetryCount int `json:"retry_count"`

	// Error message if sending failed
	// example: null
	ErrorMessage *string `json:"error_message,omitempty"`
}

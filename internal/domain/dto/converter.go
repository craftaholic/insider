package dto

import (
	"github.com/craftaholic/insider/internal/domain/entity"
)

// ConvertMessageToDTO converts a domain entity to DTO.
func ConvertMessageToDTO(msg entity.Message) MessageDTO {
	dto := MessageDTO{
		ID:           msg.ID,
		Content:      msg.Content,
		PhoneNumber:  msg.PhoneNumber,
		CreatedAt:    msg.CreatedAt,
		UpdatedAt:    msg.UpdatedAt,
		Status:       msg.Status,
		ErrorMessage: msg.ErrorMessage,
		MessageID:    msg.MessageID,
	}

	// Handle nullable SentAt
	if msg.SentAt != nil {
		dto.SentAt = msg.SentAt
	}

	// Handle nullable MessageId
	if msg.MessageID != nil {
		dto.MessageID = msg.MessageID
	}

	return dto
}

// ConvertMessagesToDTO converts a slice of domain entities to DTOs.
func ConvertMessagesToDTO(messages []entity.Message) []MessageDTO {
	dtos := make([]MessageDTO, len(messages))
	for i, msg := range messages {
		dtos[i] = ConvertMessageToDTO(msg)
	}
	return dtos
}

// CreateStandardResponse creates a standard success response.
func CreateStandardResponse(status, message string) StandardResponse {
	return StandardResponse{
		Status:  status,
		Message: message,
	}
}

// CreateErrorResponse creates an error response.
func CreateErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Error: err,
	}
}

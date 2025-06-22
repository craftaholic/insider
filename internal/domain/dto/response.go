package dto

// StandardResponse represents a standard API response.
type StandardResponse struct {
	// Status of the operation
	// example: OK
	Status string `json:"status"`

	// Descriptive message
	// example: Operation completed successfully
	Message string `json:"message"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	// Error message
	// example: Internal server error occurred
	Error string `json:"error"`
}

// swagger:response startResponse
type StartResponse struct {
	// Success response for start operation
	// in: body
	Body StandardResponse `json:"body"`
}

// swagger:response stopResponse
type StopResponse struct {
	// Success response for stop operation
	// in: body
	Body StandardResponse `json:"body"`
}

// swagger:response healthResponse
type HealthResponse struct {
	// Success response for stop operation
	// in: body
	Body StandardResponse `json:"body"`
}

// swagger:response messagesResponse
type MessagesResponse struct {
	// List of sent messages
	// in: messages
	Body []MessageDTO `json:"messages"`
}

// swagger:response errorResponse
type ErrorResponseWrapper struct {
	// Error response
	// in: body
	Body ErrorResponse `json:"body"`
}

// PaginatedMessagesResponse for future pagination metadata.
type PaginatedMessagesResponse struct {
	// List of messages
	Messages []MessageDTO `json:"messages"`

	// Pagination metadata
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination information.
type PaginationMeta struct {
	// Current page number
	// example: 1
	Page int `json:"page"`

	// Number of items per page
	// example: 10
	PerPage int `json:"per_page"`

	// Total number of items
	// example: 100
	Total int `json:"total"`

	// Total number of pages
	// example: 10
	TotalPages int `json:"total_pages"`

	// Whether there is a next page
	// example: true
	HasNext bool `json:"has_next"`

	// Whether there is a previous page
	// example: false
	HasPrev bool `json:"has_prev"`
}

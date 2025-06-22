package dto

// swagger:response startResponse
type StartResponse struct {
	// Success response for start operation
	// in: body
	Body struct {
		// Status of the operation
		// example: OK
		Status string `json:"status"`
		// Descriptive message
		// example: Automated sending started successfully
		Message string `json:"message"`
	}
}

// swagger:response stopResponse
type StopResponse struct {
	// Success response for stop operation
	// in: body
	Body struct {
		// Status of the operation
		// example: OK
		Status string `json:"status"`
		// Descriptive message
		// example: Automated sending stopped successfully
		Message string `json:"message"`
	}
}

// swagger:response messagesResponse
type MessagesResponse struct {
	// List of sent messages
	// in: body
	Body []MessageDTO `json:"body"`
}

// swagger:response errorResponse
type ErrorResponse struct {
	// Error response
	// in: body
	Body struct {
		// Error message
		// example: Internal server error occurred
		Error string `json:"error"`
	}
}

// Standard API response wrapper
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

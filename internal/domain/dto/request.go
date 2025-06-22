package dto

// swagger:parameters getSentMessages
type GetSentMessagesParams struct {
	// Page number for pagination (default: 1)
	// in: query
	// required: false
	// minimum: 1
	// example: 1
	Page int `json:"page"`
}

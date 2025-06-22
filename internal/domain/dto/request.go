package dto

// swagger:parameters getSentMessages
type GetSentMessagesParams struct {
	// Page number for pagination
	// in: query
	// required: false
	// minimum: 1
	// example: 1
	Page int `json:"page"`
}

// swagger:parameters start
type StartParams struct {
	// No parameters required for this endpoint
}

// swagger:parameters stop
type StopParams struct {
	// No parameters required for this endpoint
}

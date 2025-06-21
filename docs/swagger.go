// Package classification User API
//
// Documentation for User API
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /
//	Version: 0.0.1
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package apidocs

// Health Check Response
// swagger:response healthResponse
type HealthResponse struct {
	// in: body
	Body struct {
		// The status of the health check
		Status string `json:"status" example:"ok"`
	}
}

// Error Response
// swagger:response errorResponse
type ErrorResponse struct {
	// in: body
	Body struct {
		// Error message
		Error string `json:"error" example:"Internal server error"`
	}
}

// User Response
// swagger:response userResponse
type UserResponse struct {
	// in: body
	Body struct {
		// User ID
		ID int `json:"id" example:"1"`
		// User name
		Name string `json:"name" example:"John Doe"`
		// User email
		Email string `json:"email" example:"john@example.com"`
	}
}

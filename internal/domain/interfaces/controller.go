package interfaces

import "net/http"

type MessageController interface {
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	GetSentMessagesWithPagination(w http.ResponseWriter, r *http.Request)
}

type HealthController interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

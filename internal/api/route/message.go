package route

import (
	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewMessageRouter(router chi.Router, mc interfaces.MessageController) {
	router.Post("/service/start", mc.Start)
	router.Post("/service/stop", mc.Stop)
	router.Get("/service/status", mc.Status)
	router.Get("/message/sent", mc.GetSentMessagesWithPagination)
}

package route

import (
	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewMessageRouter(router chi.Router, mc interfaces.MessageController) {
	router.Post("/message/start", mc.Start)
	router.Post("/message/stop", mc.Stop)
	router.Get("/message/sent", mc.GetSentMessagesWithPagination)
}

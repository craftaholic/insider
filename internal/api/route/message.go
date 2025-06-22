package route

import (
	"github.com/craftaholic/insider/internal/controller"
	"github.com/go-chi/chi/v5"
)

func NewMessageRouter(router chi.Router, mc *controller.MessageController) {
	router.Post("/message/start", mc.Start)
	router.Post("/message/stop", mc.Stop)
	router.Get("/message/sent", mc.GetSentMessagesWithPagination)
}

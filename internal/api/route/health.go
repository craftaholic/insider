package route

import (
	"github.com/craftaholic/insider/internal/controller"
	"github.com/go-chi/chi/v5"
)

func NewHealthRouter(router chi.Router, hc *controller.HealthController) {
	router.Get("/health", hc.HealthCheck)
}

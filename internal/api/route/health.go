package route

import (
	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewHealthRouter(router chi.Router, hc interfaces.HealthController) {
	router.Get("/health", hc.HealthCheck)
}

package controller

import (
	"net/http"

	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/craftaholic/insider/internal/utils"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger := log.BaseLogger.FromCtx(r.Context()).WithFields("controller", utils.GetStructName(hc))
	logger.Info("Processing health check request")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	logger.Info("Finished health check request")
}

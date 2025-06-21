package custommiddleware

import (
	"net/http"

	"github.com/craftaholic/insider/internal/shared/log"
	"github.com/craftaholic/insider/internal/utils"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := utils.GenerateUUIDv7()
		logger := log.BaseLogger.WithFields(
			"request_id",
			requestID,
			"remote_addr",
			r.RemoteAddr,
			"method",
			r.Method,
			"url",
			r.URL.String(),
		)
		ctx := logger.WithCtx(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

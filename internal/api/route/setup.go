package route

import (
	"net/http"

	custommiddleware "github.com/craftaholic/insider/internal/api/middleware"
	"github.com/craftaholic/insider/internal/bootstrap"
	"github.com/craftaholic/insider/internal/shared/constant"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	swagger "github.com/go-openapi/runtime/middleware"
)

func SetupRoute(app bootstrap.Application) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, // Use your allowed origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Origin",
			"X-Requested-With",
		},
		ExposedHeaders:   []string{"Link", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           constant.CorsMaxAge, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// Define middleware
	r.Use(middleware.RealIP)
	r.Use(custommiddleware.LoggingMiddleware)
	r.Use(middleware.Recoverer)

	// Swagger UI
	opts := swagger.SwaggerUIOpts{SpecURL: "/swagger.json"}
	sh := swagger.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))

	// Public APIs
	r.Group(func(r chi.Router) {
		NewHealthRouter(r, app.HealthController)
		NewMessageRouter(r, app.MessageController)
	})

	return r
}

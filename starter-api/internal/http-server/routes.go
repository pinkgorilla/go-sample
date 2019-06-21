package httpserver

import (
	"time"

	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/pinkgorilla/go-sample/starter-api/internal/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func buildRoutes(a *app.App) chi.Router {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// A good base middleware stack
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(app.InjectorMiddleware(a))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chimw.Timeout(60 * time.Second))

	// prometheus handler
	r.Handle("/metrics", promhttp.Handler())

	// callback method

	// Public API
	r.Route("/v1", func(r chi.Router) {
		r.Get("/cashout", handlers.CreateCashOutHandler())
	})

	return r
}

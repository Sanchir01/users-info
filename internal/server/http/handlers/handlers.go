package httphandlers

import (
	"net/http"

	_ "github.com/Sanchir01/users-info/docs"
	"github.com/Sanchir01/users-info/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(PrometheusMiddleware)
	router.Route("/apiv1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", handlers.UserHandler.GetAllUsers)
			r.Delete("/{id}", handlers.UserHandler.DeleteUser)
			r.Post("/create", handlers.UserHandler.CreateUser)
			r.Patch("/{id}", handlers.UserHandler.UpdateUser)
		})
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}

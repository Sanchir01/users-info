package httphandlers

import (
	_ "github.com/Sanchir01/users-info/docs"
	"github.com/Sanchir01/users-info/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)

	router.Route("/apiv1", func(r chi.Router) {
		r.Group(func(r chi.Router) {

			r.Get("/info", func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte("Hello, World!")); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			})

		})

	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	return router
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}

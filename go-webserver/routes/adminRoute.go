package routes

import (
	"github.com/go-chi/chi/v5"
)

// ----- Create Router ---------------------------

func CreateAdminRouter(metricsConfig MetricsConfig) chi.Router {
	r := chi.NewRouter()

	// Route:	/metrics
	r.Get("/metrics", metricsConfig.GetMetricsHandler)
	// Route:	/reset
	r.Get("/reset", metricsConfig.ResetMetricsHandler)

	return r
}

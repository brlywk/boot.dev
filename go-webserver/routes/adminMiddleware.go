package routes

import (
	"fmt"
	"net/http"
)

// ----- Types -----------------------------------

type MetricsConfig struct {
	FileserverHits int
}

func (cfg *MetricsConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

// ----- Handlers --------------------------------

func (cfg *MetricsConfig) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	template := `
	<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>

	</html>
	`

	w.Write([]byte(fmt.Sprintf(template, cfg.FileserverHits)))
}

func (cfg *MetricsConfig) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.FileserverHits = 0
}

package httpserver

import (
	"net/http"

	"github.com/woragis/management/backend/server/internal/middleware"
)

func NewHandler(app *App, cfg middleware.Config) http.Handler {
	mux := http.NewServeMux()
	Mount(mux, app)
	if cfg.MetricsEnabled {
		mux.HandleFunc("GET /metrics", middleware.PrometheusHandler())
	}
	return middleware.Chain(cfg, mux)
}

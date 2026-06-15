package httpserver

import "net/http"

func Mount(mux *http.ServeMux, app *App) {
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))
}

package httpserver

import (
	"net/http"

	"github.com/woragis/management/backend/server/internal/middleware"
)

func Mount(mux *http.ServeMux, app *App) {
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))

	if app.DevProjects != nil {
		dh := newDevprojectHandler(app.DevProjects)
		admin := func(h http.HandlerFunc) http.Handler {
			return middleware.AdminAuth(app.AdminAPIKey, h)
		}
		mux.Handle("GET /v1/admin/projects", admin(dh.list))
		mux.Handle("POST /v1/admin/projects", admin(dh.create))
		mux.Handle("GET /v1/admin/projects/{id}", admin(dh.get))
		mux.Handle("PATCH /v1/admin/projects/{id}", admin(dh.update))
		mux.Handle("DELETE /v1/admin/projects/{id}", admin(dh.delete))
		mux.Handle("POST /v1/admin/projects/{id}/links", admin(dh.createLink))
		mux.Handle("DELETE /v1/admin/projects/{id}/links/{linkId}", admin(dh.deleteLink))
		mux.Handle("POST /v1/admin/projects/{id}/domains", admin(dh.createDomain))
		mux.Handle("DELETE /v1/admin/projects/{id}/domains/{domainId}", admin(dh.deleteDomain))
		mux.Handle("GET /v1/admin/projects/{id}/secrets", admin(dh.listSecrets))
		mux.Handle("GET /v1/admin/projects/{id}/secrets/{secretId}", admin(dh.getSecret))
		mux.Handle("POST /v1/admin/projects/{id}/secrets", admin(dh.createSecret))
		mux.Handle("DELETE /v1/admin/projects/{id}/secrets/{secretId}", admin(dh.deleteSecret))
		mux.Handle("POST /v1/admin/projects/{id}/gallery", admin(dh.createGallery))
		mux.Handle("DELETE /v1/admin/projects/{id}/gallery/{itemId}", admin(dh.deleteGallery))
	}

	if app.Media != nil {
		mh := newMediaHandler(app.Media)
		admin := func(h http.HandlerFunc) http.Handler {
			return middleware.AdminAuth(app.AdminAPIKey, h)
		}
		mux.Handle("GET /v1/admin/media", admin(mh.list))
		mux.Handle("POST /v1/admin/media", admin(mh.upload))
		mux.Handle("GET /v1/admin/media/{id}", admin(mh.get))
		mux.Handle("DELETE /v1/admin/media/{id}", admin(mh.delete))
		mux.HandleFunc("GET /v1/public/media/{id}/file", mh.serveFile)
		mux.HandleFunc("GET /v1/public/media/{id}", mh.get)
	}

	if app.Profile != nil {
		ph := newProfileHandler(app.Profile)
		admin := func(h http.HandlerFunc) http.Handler {
			return middleware.AdminAuth(app.AdminAPIKey, h)
		}
		mux.Handle("GET /v1/admin/profile", admin(ph.getAdmin))
		mux.Handle("PATCH /v1/admin/profile", admin(ph.update))
		mux.HandleFunc("GET /v1/public/profile", ph.getPublic)
	}
}

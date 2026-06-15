package httpserver

import (
	"net/http"
	"strings"

	"github.com/woragis/management/backend/server/internal/apperrors"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
)

type publicHandler struct {
	projects *devprojectsvc.Service
}

func newPublicHandler(projects *devprojectsvc.Service) *publicHandler {
	return &publicHandler{projects: projects}
}

func (h *publicHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	featuredOnly := strings.EqualFold(r.URL.Query().Get("featured"), "true")
	items, err := h.projects.ListPublic(r.Context(), featuredOnly)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apperrors.WriteJSON(w, http.StatusOK, items)
}

func (h *publicHandler) getProject(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(r.PathValue("slug"))
	if slug == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, "Project slug is invalid."))
		return
	}
	item, err := h.projects.GetPublicBySlug(r.Context(), slug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apperrors.WriteJSON(w, http.StatusOK, item)
}

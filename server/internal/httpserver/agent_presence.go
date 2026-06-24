package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
)

func (h *agentToolsHandler) listSocialPosts(w http.ResponseWriter, r *http.Request) {
	if h.presenceH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Presence service unavailable."))
		return
	}
	h.presenceH.listPosts(w, r)
}

func (h *agentToolsHandler) listPostTemplates(w http.ResponseWriter, r *http.Request) {
	if h.presenceH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Presence service unavailable."))
		return
	}
	h.presenceH.listTemplates(w, r)
}

func (h *agentToolsHandler) createSocialPost(w http.ResponseWriter, r *http.Request) {
	if h.presenceH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Presence service unavailable."))
		return
	}
	var body socialPostBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	if strings.TrimSpace(strings.ToLower(body.Status)) == "published" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Cannot create posts as published via agent; use draft or scheduled."))
		return
	}
	if body.Status == "" {
		body.Status = "draft"
	}
	row, err := h.presence.CreatePost(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *agentToolsHandler) applyPostTemplate(w http.ResponseWriter, r *http.Request) {
	if h.presence == nil || h.devProjects == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Presence or projects service unavailable."))
		return
	}
	var body struct {
		TemplateSlug string    `json:"templateSlug"`
		ProjectID    uuid.UUID `json:"projectId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	if body.ProjectID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "projectId is required."))
		return
	}
	project, err := h.devProjects.GetByID(r.Context(), body.ProjectID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	rendered, err := h.presence.RenderTemplate(r.Context(), body.TemplateSlug, presencesvc.TemplateVarsFromProject(project))
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rendered)
}

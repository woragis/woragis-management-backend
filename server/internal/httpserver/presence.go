package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
	"github.com/woragis/management/backend/server/internal/presence/repository"
)

type presenceHandler struct {
	svc *presencesvc.Service
}

func newPresenceHandler(svc *presencesvc.Service) *presenceHandler {
	return &presenceHandler{svc: svc}
}

func (h *presenceHandler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var projectID *uuid.UUID
	if pid := q.Get("projectId"); pid != "" {
		if id, err := uuid.Parse(pid); err == nil {
			projectID = &id
		}
	}
	rows, err := h.svc.ListCampaigns(r.Context(), presencesvc.CampaignFilter{
		Goal:       q.Get("goal"),
		ProjectID:  projectID,
		ActiveOnly: q.Get("active") != "false",
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *presenceHandler) getCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetCampaign(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) createCampaign(w http.ResponseWriter, r *http.Request) {
	var body campaignBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateCampaign(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *presenceHandler) updateCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body campaignUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateCampaign(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) deleteCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteCampaign(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *presenceHandler) listTemplates(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rows, err := h.svc.ListTemplates(r.Context(), presencesvc.TemplateFilter{
		Platform:   q.Get("platform"),
		Goal:       q.Get("goal"),
		ActiveOnly: q.Get("active") != "false",
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *presenceHandler) getTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetTemplate(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) createTemplate(w http.ResponseWriter, r *http.Request) {
	var body postTemplateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateTemplate(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *presenceHandler) updateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body postTemplateUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateTemplate(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteTemplate(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *presenceHandler) listPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var projectID, campaignID *uuid.UUID
	if pid := q.Get("projectId"); pid != "" {
		if id, err := uuid.Parse(pid); err == nil {
			projectID = &id
		}
	}
	if cid := q.Get("campaignId"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			campaignID = &id
		}
	}
	limit := 100
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	rows, err := h.svc.ListPosts(r.Context(), presencesvc.PostFilter{
		Platform:   q.Get("platform"),
		Goal:       q.Get("goal"),
		Status:     q.Get("status"),
		ProjectID:  projectID,
		CampaignID: campaignID,
		Limit:      limit,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *presenceHandler) getPost(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetPost(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var body socialPostBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreatePost(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *presenceHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body socialPostUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdatePost(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *presenceHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeletePost(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type campaignBody struct {
	Name        string     `json:"name"`
	Goal        string     `json:"goal"`
	Description string     `json:"description"`
	ProjectID   *uuid.UUID `json:"projectId"`
	StartDate   *string    `json:"startDate"`
	EndDate     *string    `json:"endDate"`
	Active      bool       `json:"active"`
}

func (b campaignBody) toCreate() presencesvc.CreateCampaignInput {
	return presencesvc.CreateCampaignInput{
		Name:        b.Name,
		Goal:        b.Goal,
		Description: b.Description,
		ProjectID:   b.ProjectID,
		StartDate:   parseOptionalDate(b.StartDate),
		EndDate:     parseOptionalDate(b.EndDate),
		Active:      b.Active,
	}
}

type campaignUpdateBody struct {
	Name        *string    `json:"name"`
	Goal        *string    `json:"goal"`
	Description *string    `json:"description"`
	ProjectID   *uuid.UUID `json:"projectId"`
	StartDate   *string    `json:"startDate"`
	EndDate     *string    `json:"endDate"`
	Active      *bool      `json:"active"`
}

func (b campaignUpdateBody) toUpdate() presencesvc.UpdateCampaignInput {
	in := presencesvc.UpdateCampaignInput{
		Name:        b.Name,
		Goal:        b.Goal,
		Description: b.Description,
		Active:      b.Active,
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	if b.StartDate != nil {
		in.StartDate = parseOptionalDate(b.StartDate)
		in.StartSet = true
	}
	if b.EndDate != nil {
		in.EndDate = parseOptionalDate(b.EndDate)
		in.EndSet = true
	}
	return in
}

type postTemplateBody struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Goal     string `json:"goal"`
	Body     string `json:"body"`
	Active   bool   `json:"active"`
}

func (b postTemplateBody) toCreate() presencesvc.CreateTemplateInput {
	return presencesvc.CreateTemplateInput(b)
}

type postTemplateUpdateBody struct {
	Slug     *string `json:"slug"`
	Name     *string `json:"name"`
	Platform *string `json:"platform"`
	Goal     *string `json:"goal"`
	Body     *string `json:"body"`
	Active   *bool   `json:"active"`
}

func (b postTemplateUpdateBody) toUpdate() presencesvc.UpdateTemplateInput {
	return presencesvc.UpdateTemplateInput(b)
}

type socialPostBody struct {
	ProjectID    *uuid.UUID `json:"projectId"`
	CampaignID   *uuid.UUID `json:"campaignId"`
	Platform     string     `json:"platform"`
	Goal         string     `json:"goal"`
	Status       string     `json:"status"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	Hook         string     `json:"hook"`
	CTA          string     `json:"cta"`
	TemplateSlug string     `json:"templateSlug"`
	ScheduledAt  *string    `json:"scheduledAt"`
	PublishedAt  *string    `json:"publishedAt"`
	PublishedURL string     `json:"publishedUrl"`
	Notes        string     `json:"notes"`
}

func (b socialPostBody) toCreate() presencesvc.CreatePostInput {
	return presencesvc.CreatePostInput{
		ProjectID:    b.ProjectID,
		CampaignID:   b.CampaignID,
		Platform:     b.Platform,
		Goal:         b.Goal,
		Status:       b.Status,
		Title:        b.Title,
		Body:         b.Body,
		Hook:         b.Hook,
		CTA:          b.CTA,
		TemplateSlug: b.TemplateSlug,
		ScheduledAt:  parseOptionalDateTime(b.ScheduledAt),
		PublishedAt:  parseOptionalDateTime(b.PublishedAt),
		PublishedURL: b.PublishedURL,
		Notes:        b.Notes,
	}
}

type socialPostUpdateBody struct {
	ProjectID    *uuid.UUID `json:"projectId"`
	CampaignID   *uuid.UUID `json:"campaignId"`
	Platform     *string    `json:"platform"`
	Goal         *string    `json:"goal"`
	Status       *string    `json:"status"`
	Title        *string    `json:"title"`
	Body         *string    `json:"body"`
	Hook         *string    `json:"hook"`
	CTA          *string    `json:"cta"`
	TemplateSlug *string    `json:"templateSlug"`
	ScheduledAt  *string    `json:"scheduledAt"`
	PublishedAt  *string    `json:"publishedAt"`
	PublishedURL *string    `json:"publishedUrl"`
	Notes        *string    `json:"notes"`
}

func (b socialPostUpdateBody) toUpdate() presencesvc.UpdatePostInput {
	in := presencesvc.UpdatePostInput{
		Platform:     b.Platform,
		Goal:         b.Goal,
		Status:       b.Status,
		Title:        b.Title,
		Body:         b.Body,
		Hook:         b.Hook,
		CTA:          b.CTA,
		TemplateSlug: b.TemplateSlug,
		PublishedURL: b.PublishedURL,
		Notes:        b.Notes,
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	if b.CampaignID != nil {
		in.CampaignID = b.CampaignID
		in.CampaignSet = true
	}
	if b.ScheduledAt != nil {
		in.ScheduledAt = parseOptionalDateTime(b.ScheduledAt)
		in.ScheduledSet = true
	}
	if b.PublishedAt != nil {
		in.PublishedAt = parseOptionalDateTime(b.PublishedAt)
		in.PublishedSet = true
	}
	return in
}

func parseOptionalDate(s *string) *time.Time {
	if s == nil {
		return nil
	}
	return repository.ParseDate(*s)
}

func parseOptionalDateTime(s *string) *time.Time {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	raw := strings.TrimSpace(*s)
	if t := repository.ParseDate(raw); t != nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t
	}
	return nil
}

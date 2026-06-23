package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
)

type messagingHandler struct {
	svc *messagingsvc.Service
}

func newMessagingHandler(svc *messagingsvc.Service) *messagingHandler {
	return &messagingHandler{svc: svc}
}

func (h *messagingHandler) listDestinations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := messagingsvc.DestinationFilter{
		Channel:    q.Get("channel"),
		Query:      q.Get("q"),
		ActiveOnly: q.Get("active") != "false",
	}
	rows, err := h.svc.ListDestinations(r.Context(), f)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *messagingHandler) getDestination(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetDestination(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *messagingHandler) createDestination(w http.ResponseWriter, r *http.Request) {
	var body destinationBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateDestination(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *messagingHandler) updateDestination(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body destinationUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateDestination(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *messagingHandler) deleteDestination(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteDestination(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *messagingHandler) listTemplates(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := messagingsvc.TemplateFilter{
		ProgramSlug: q.Get("programSlug"),
		ActiveOnly:  q.Get("active") != "false",
	}
	if did := q.Get("destinationId"); did != "" {
		if id, err := uuid.Parse(did); err == nil {
			f.DestinationID = &id
		}
	}
	rows, err := h.svc.ListTemplates(r.Context(), f)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *messagingHandler) getTemplate(w http.ResponseWriter, r *http.Request) {
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

func (h *messagingHandler) createTemplate(w http.ResponseWriter, r *http.Request) {
	var body messageTemplateBody
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

func (h *messagingHandler) updateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body messageTemplateUpdateBody
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

func (h *messagingHandler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
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

func (h *messagingHandler) listJobs(w http.ResponseWriter, r *http.Request) {
	enabledOnly := r.URL.Query().Get("enabled") == "true"
	rows, err := h.svc.ListJobs(r.Context(), enabledOnly)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *messagingHandler) getJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetJob(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *messagingHandler) createJob(w http.ResponseWriter, r *http.Request) {
	var body jobBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateJob(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *messagingHandler) updateJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body jobUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateJob(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *messagingHandler) deleteJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteJob(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *messagingHandler) listDeliveries(w http.ResponseWriter, r *http.Request) {
	var destID *uuid.UUID
	if did := r.URL.Query().Get("destinationId"); did != "" {
		if id, err := uuid.Parse(did); err == nil {
			destID = &id
		}
	}
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	rows, err := h.svc.ListDeliveries(r.Context(), destID, limit)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

type destinationBody struct {
	Channel          string         `json:"channel"`
	ExternalID       string         `json:"externalId"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Responsibilities string         `json:"responsibilities"`
	Tags             []string       `json:"tags"`
	Metadata         map[string]any `json:"metadata"`
	Active           bool           `json:"active"`
}

func (b destinationBody) toCreate() messagingsvc.CreateDestinationInput {
	return messagingsvc.CreateDestinationInput(b)
}

type destinationUpdateBody struct {
	Channel          *string        `json:"channel"`
	ExternalID       *string        `json:"externalId"`
	Name             *string        `json:"name"`
	Description      *string        `json:"description"`
	Responsibilities *string        `json:"responsibilities"`
	Tags             []string       `json:"tags"`
	Metadata         map[string]any `json:"metadata"`
	Active           *bool          `json:"active"`
}

func (b destinationUpdateBody) toUpdate() messagingsvc.UpdateDestinationInput {
	in := messagingsvc.UpdateDestinationInput{
		Channel:          b.Channel,
		ExternalID:       b.ExternalID,
		Name:             b.Name,
		Description:      b.Description,
		Responsibilities: b.Responsibilities,
		Active:           b.Active,
	}
	if b.Tags != nil {
		in.Tags = b.Tags
		in.TagsSet = true
	}
	if b.Metadata != nil {
		in.Metadata = b.Metadata
		in.MetadataSet = true
	}
	return in
}

type messageTemplateBody struct {
	DestinationID *uuid.UUID `json:"destinationId"`
	ProgramSlug   string     `json:"programSlug"`
	Slug          string     `json:"slug"`
	Name          string     `json:"name"`
	Body          string     `json:"body"`
	ComposeMode   string     `json:"composeMode"`
	AIPromptHint  string     `json:"aiPromptHint"`
	Active        bool       `json:"active"`
}

func (b messageTemplateBody) toCreate() messagingsvc.CreateTemplateInput {
	return messagingsvc.CreateTemplateInput(b)
}

type messageTemplateUpdateBody struct {
	DestinationID *uuid.UUID `json:"destinationId"`
	ProgramSlug   *string    `json:"programSlug"`
	Slug          *string    `json:"slug"`
	Name          *string    `json:"name"`
	Body          *string    `json:"body"`
	ComposeMode   *string    `json:"composeMode"`
	AIPromptHint  *string    `json:"aiPromptHint"`
	Active        *bool      `json:"active"`
}

func (b messageTemplateUpdateBody) toUpdate() messagingsvc.UpdateTemplateInput {
	in := messagingsvc.UpdateTemplateInput{
		ProgramSlug:  b.ProgramSlug,
		Slug:         b.Slug,
		Name:         b.Name,
		Body:         b.Body,
		ComposeMode:  b.ComposeMode,
		AIPromptHint: b.AIPromptHint,
		Active:       b.Active,
	}
	if b.DestinationID != nil {
		in.DestinationID = b.DestinationID
		in.DestinationSet = true
	}
	return in
}

type jobBody struct {
	Name          string    `json:"name"`
	DestinationID uuid.UUID `json:"destinationId"`
	TemplateSlug  string    `json:"templateSlug"`
	ProgramAction string    `json:"programAction"`
	CronExpr      string    `json:"cronExpr"`
	Timezone      string    `json:"timezone"`
	Enabled       bool      `json:"enabled"`
}

func (b jobBody) toCreate() messagingsvc.CreateJobInput {
	return messagingsvc.CreateJobInput(b)
}

type jobUpdateBody struct {
	Name          *string    `json:"name"`
	DestinationID *uuid.UUID `json:"destinationId"`
	TemplateSlug  *string    `json:"templateSlug"`
	ProgramAction *string    `json:"programAction"`
	CronExpr      *string    `json:"cronExpr"`
	Timezone      *string    `json:"timezone"`
	Enabled       *bool      `json:"enabled"`
}

func (b jobUpdateBody) toUpdate() messagingsvc.UpdateJobInput {
	return messagingsvc.UpdateJobInput(b)
}

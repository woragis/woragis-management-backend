package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/messaging/executor"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	msgtemplaterender "github.com/woragis/management/backend/server/internal/messaging/templaterender"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
	"github.com/woragis/management/backend/server/internal/telegramworkerclient"
	"gorm.io/datatypes"
)

type messagingHandler struct {
	svc       *messagingsvc.Service
	scheduler *executor.Executor
	whatsapp  *whatsappworkerclient.Client
	telegram  *telegramworkerclient.Client
}

func newMessagingHandler(svc *messagingsvc.Service, scheduler *executor.Executor, whatsapp *whatsappworkerclient.Client, telegram *telegramworkerclient.Client) *messagingHandler {
	return &messagingHandler{svc: svc, scheduler: scheduler, whatsapp: whatsapp, telegram: telegram}
}

func (h *messagingHandler) syncWhatsAppDestinations(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.SyncWhatsAppDestinations(r.Context(), h.whatsapp)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, res)
}

func (h *messagingHandler) syncTelegramDestinations(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.SyncTelegramDestinations(r.Context(), h.telegram)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, res)
}

func (h *messagingHandler) resolveDestination(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	externalID := r.URL.Query().Get("externalId")
	row, err := h.svc.ResolveDestination(r.Context(), channel, externalID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *messagingHandler) catalogFields(w http.ResponseWriter, r *http.Request) {
	program := r.URL.Query().Get("program")
	if program == "" {
		program = "custom"
	}
	apperrors.WriteJSON(w, http.StatusOK, msgtemplaterender.CatalogFields(program))
}

func (h *messagingHandler) previewTemplate(w http.ResponseWriter, r *http.Request) {
	var body templatePreviewBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	tplID, err := uuid.Parse(strings.TrimSpace(body.TemplateID))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid template id."))
		return
	}
	tpl, err := h.svc.GetTemplate(r.Context(), tplID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	job := &models.ScheduledJob{
		ProgramAction: strings.TrimSpace(body.ProgramAction),
		DataSource:    jsonRawMap(body.DataSource),
	}
	if h.scheduler == nil {
		apperrors.WriteJSON(w, http.StatusOK, map[string]any{"body": tpl.Body, "skipped": false})
		return
	}
	res, err := h.scheduler.PreviewTemplate(r.Context(), tpl, job)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, map[string]any{
		"body":        res.Body,
		"data":        res.Data,
		"skipped":     res.Skipped,
		"skipReason":  res.SkipReason,
		"externalRef": res.ExternalRef,
	})
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
	DestinationID *uuid.UUID        `json:"destinationId"`
	ProgramSlug   string            `json:"programSlug"`
	Slug          string            `json:"slug"`
	Name          string            `json:"name"`
	Body          string            `json:"body"`
	ComposeMode   string            `json:"composeMode"`
	AIPromptHint  string            `json:"aiPromptHint"`
	Bindings      map[string]string `json:"bindings"`
	Active        bool              `json:"active"`
}

func (b messageTemplateBody) toCreate() messagingsvc.CreateTemplateInput {
	return messagingsvc.CreateTemplateInput(b)
}

type messageTemplateUpdateBody struct {
	DestinationID *uuid.UUID        `json:"destinationId"`
	ProgramSlug   *string           `json:"programSlug"`
	Slug          *string           `json:"slug"`
	Name          *string           `json:"name"`
	Body          *string           `json:"body"`
	ComposeMode   *string           `json:"composeMode"`
	AIPromptHint  *string           `json:"aiPromptHint"`
	Bindings      map[string]string `json:"bindings"`
	Active        *bool             `json:"active"`
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
	if b.Bindings != nil {
		in.Bindings = b.Bindings
		in.BindingsSet = true
	}
	return in
}

type jobBody struct {
	Name          string         `json:"name"`
	DestinationID uuid.UUID      `json:"destinationId"`
	TemplateSlug  string         `json:"templateSlug"`
	ProgramAction string         `json:"programAction"`
	DataSource    map[string]any `json:"dataSource"`
	CronExpr      string         `json:"cronExpr"`
	Timezone      string         `json:"timezone"`
	Enabled       bool           `json:"enabled"`
}

func (b jobBody) toCreate() messagingsvc.CreateJobInput {
	return messagingsvc.CreateJobInput(b)
}

type jobUpdateBody struct {
	Name          *string        `json:"name"`
	DestinationID *uuid.UUID     `json:"destinationId"`
	TemplateSlug  *string        `json:"templateSlug"`
	ProgramAction *string        `json:"programAction"`
	DataSource    map[string]any `json:"dataSource"`
	CronExpr      *string        `json:"cronExpr"`
	Timezone      *string        `json:"timezone"`
	Enabled       *bool          `json:"enabled"`
}

func (b jobUpdateBody) toUpdate() messagingsvc.UpdateJobInput {
	in := messagingsvc.UpdateJobInput{
		Name:          b.Name,
		DestinationID: b.DestinationID,
		TemplateSlug:  b.TemplateSlug,
		ProgramAction: b.ProgramAction,
		CronExpr:      b.CronExpr,
		Timezone:      b.Timezone,
		Enabled:       b.Enabled,
	}
	if b.DataSource != nil {
		in.DataSource = b.DataSource
		in.DataSourceSet = true
	}
	return in
}

type templatePreviewBody struct {
	TemplateID    string         `json:"templateId"`
	ProgramAction string         `json:"programAction"`
	DataSource    map[string]any `json:"dataSource"`
}

func jsonRawMap(m map[string]any) datatypes.JSON {
	if len(m) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	b, _ := json.Marshal(m)
	return datatypes.JSON(b)
}

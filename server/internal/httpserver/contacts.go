package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	contactssvc "github.com/woragis/management/backend/server/internal/contacts/service"
	financesvc "github.com/woragis/management/backend/server/internal/finance/service"
)

type contactsHandler struct {
	svc     *contactssvc.Service
	finance *financesvc.Service
}

func newContactsHandler(svc *contactssvc.Service, finance *financesvc.Service) *contactsHandler {
	return &contactsHandler{svc: svc, finance: finance}
}

func (h *contactsHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := contactssvc.ListFilter{
		Query:        q.Get("q"),
		Relationship: q.Get("relationship"),
		Organization: q.Get("organization"),
		Stage:        q.Get("stage"),
		ActiveOnly:   q.Get("active") != "false",
	}
	if pid := q.Get("projectId"); pid != "" {
		if id, err := uuid.Parse(pid); err == nil {
			f.ProjectID = &id
		}
	}
	rows, err := h.svc.List(r.Context(), f)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contactsHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contactsHandler) create(w http.ResponseWriter, r *http.Request) {
	var body contactBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.Create(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *contactsHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body contactUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.Update(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contactsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *contactsHandler) listInteractions(w http.ResponseWriter, r *http.Request) {
	contactID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	rows, err := h.svc.ListInteractions(r.Context(), contactID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contactsHandler) createInteraction(w http.ResponseWriter, r *http.Request) {
	contactID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body interactionBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateInteraction(r.Context(), contactID, body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *contactsHandler) contactFinance(w http.ResponseWriter, r *http.Request) {
	if h.finance == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	out, err := h.finance.ContactFinance(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, out)
}

type contactBody struct {
	Name           string     `json:"name"`
	DisplayName    string     `json:"displayName"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Telegram       string     `json:"telegram"`
	Whatsapp       string     `json:"whatsapp"`
	Organization   string     `json:"organization"`
	RoleTitle      string     `json:"roleTitle"`
	Relationship   string     `json:"relationship"`
	Stage          string     `json:"stage"`
	Source         string     `json:"source"`
	Notes          string     `json:"notes"`
	Tags           []string   `json:"tags"`
	ProjectID      *uuid.UUID `json:"projectId"`
	NextFollowUpAt *time.Time `json:"nextFollowUpAt"`
	Active         bool       `json:"active"`
}

func (b contactBody) toCreate() contactssvc.CreateContactInput {
	return contactssvc.CreateContactInput(b)
}

type contactUpdateBody struct {
	Name           *string    `json:"name"`
	DisplayName    *string    `json:"displayName"`
	Email          *string    `json:"email"`
	Phone          *string    `json:"phone"`
	Telegram       *string    `json:"telegram"`
	Whatsapp       *string    `json:"whatsapp"`
	Organization   *string    `json:"organization"`
	RoleTitle      *string    `json:"roleTitle"`
	Relationship   *string    `json:"relationship"`
	Stage          *string    `json:"stage"`
	Source         *string    `json:"source"`
	Notes          *string    `json:"notes"`
	Tags           []string   `json:"tags"`
	ProjectID      *uuid.UUID `json:"projectId"`
	NextFollowUpAt *time.Time `json:"nextFollowUpAt"`
	Active         *bool      `json:"active"`
}

func (b contactUpdateBody) toUpdate() contactssvc.UpdateContactInput {
	in := contactssvc.UpdateContactInput{
		Name:         b.Name,
		DisplayName:  b.DisplayName,
		Email:        b.Email,
		Phone:        b.Phone,
		Telegram:     b.Telegram,
		Whatsapp:     b.Whatsapp,
		Organization: b.Organization,
		RoleTitle:    b.RoleTitle,
		Relationship: b.Relationship,
		Stage:        b.Stage,
		Source:       b.Source,
		Notes:        b.Notes,
		Active:       b.Active,
	}
	if b.Tags != nil {
		in.Tags = b.Tags
		in.TagsSet = true
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	if b.NextFollowUpAt != nil {
		in.NextFollowUpAt = b.NextFollowUpAt
		in.NextFollowUpSet = true
	}
	return in
}

type interactionBody struct {
	Type       string    `json:"type"`
	Channel    string    `json:"channel"`
	Summary    string    `json:"summary"`
	HappenedAt time.Time `json:"happenedAt"`
}

func (b interactionBody) toCreate() contactssvc.CreateInteractionInput {
	return contactssvc.CreateInteractionInput(b)
}

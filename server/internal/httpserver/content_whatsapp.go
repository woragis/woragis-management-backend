package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/woragis/management/backend/server/internal/apperrors"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
)

func (h *contentHandler) getSettings(w http.ResponseWriter, r *http.Request) {
	row, err := h.svc.GetSettings(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var body settingsUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateSettings(r.Context(), body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) listWhatsappTemplates(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListWhatsappTemplates(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contentHandler) updateWhatsappTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body whatsappTemplateUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateWhatsappTemplate(r.Context(), id, contentsvc.WhatsappTemplateInput{
		Name: body.Name,
		Body: body.Body,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) whatsappPreview(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	tpl := r.URL.Query().Get("type")
	if tpl == "" {
		tpl = r.URL.Query().Get("template")
	}
	msg, err := h.svc.PreviewWhatsapp(r.Context(), id, tpl)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func (h *contentHandler) whatsappSendNow(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	res, err := h.svc.SendWhatsappNow(r.Context(), id, body.Type)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, res)
}

func (h *contentHandler) whatsappWorkerStatus(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.WhatsappWorkerStatus(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, res)
}

func (h *contentHandler) whatsappWorkerQR(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.WhatsappWorkerQR(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, res)
}

type settingsUpdateBody struct {
	Timezone             *string `json:"timezone"`
	ProblemPostTime      *string `json:"problemPostTime"`
	DiscussionPostTime   *string `json:"discussionPostTime"`
	SolutionPostTime     *string `json:"solutionPostTime"`
	WeeklySummaryDay     *string `json:"weeklySummaryDay"`
	WeeklySummaryTime    *string `json:"weeklySummaryTime"`
	DiscussionEnabled    *bool   `json:"discussionEnabled"`
	InviteLink           *string `json:"inviteLink"`
	DefaultStudyPlanSlug *string `json:"defaultStudyPlanSlug"`
	NextTheme            *string `json:"nextTheme"`
}

func (b settingsUpdateBody) toUpdate() contentsvc.UpdateSettingsInput {
	return contentsvc.UpdateSettingsInput(b)
}

type whatsappTemplateUpdateBody struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func handleInternalDispatch(svc *contentsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dispatchType := r.URL.Query().Get("type")
		date := r.URL.Query().Get("date")
		res, err := svc.Dispatch(r.Context(), dispatchType, date)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, res)
	}
}

func handleInternalWhatsappStatus(svc *contentsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseUUID(r.PathValue("id"))
		if err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
			return
		}
		var body struct {
			ProblemSent    bool `json:"problemSent"`
			DiscussionSent bool `json:"discussionSent"`
			SolutionSent   bool `json:"solutionSent"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
			return
		}
		if err := svc.PatchWhatsappStatus(r.Context(), id, contentsvc.WhatsappStatusPatch(body)); err != nil {
			apperrors.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleInternalSettings(svc *contentsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		row, err := svc.GetSettings(r.Context())
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, row)
	}
}
